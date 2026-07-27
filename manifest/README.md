# Kubernetes Manifests — GitOps with Argo CD and Kustomize

Everything that runs inside the EKS cluster, declared as Kustomize bases and per-environment overlays. Argo CD watches this folder and reconciles the cluster against it, so this directory is the source of truth for cluster state — nobody runs `kubectl apply` by hand.

## How the reconciliation works

Terraform bootstraps exactly one Argo CD Application (see [infras/README.md](../infras/README.md)), pointing at `argocd-appsets/overlays/main`. That single entry point applies three ApplicationSets, and each one discovers its own Applications from the Git tree:

```text
                Terraform (one-time bootstrap)
                            │
                            ▼
          argocd-appsets/overlays/main   ← the root Application
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
  infrastructure-      kubernetes-         service-appset
     appset              appset                  │
        │                   │            git directory generator:
        ▼                   ▼            manifest/services/**/overlays/main
 infrastructures/     kubernetes/                │
  overlays/main       overlays/main    ┌─────────┴─────────┐
        │                   │          ▼                   ▼
   ServiceAccounts     StorageClass  chat-backend      chat-frontend
   (ns: chat-main)     (ns: kube-    (ns: chat-main)  (ns: chat-main)
                        system)
```

The service ApplicationSet uses a Git directory generator with a `manifest/services/**/overlays/main` glob, so adding a new service means creating a folder — no Argo CD config changes, nothing to register. All three sets run with `prune: true` and `selfHeal: true`, meaning drift is corrected automatically and deleted manifests are removed from the cluster.

## Layout

```text
manifest/
├── argocd-appsets/           # The ApplicationSets that discover everything else
│   └── overlays/main/
│       ├── infrastructure-appset.yml
│       ├── kubernetes-appset.yml
│       └── service-appset.yml
├── infrastructures/          # Shared per-namespace resources (ServiceAccounts)
│   └── overlays/{develop,main}/
├── kubernetes/               # Cluster-scoped resources (gp3 StorageClass)
│   └── overlays/{develop,main}/
└── services/                 # One folder per application
    ├── chat-backend-service/
    │   ├── base/             # Deployment, Service, HPA
    │   └── overlays/{develop,main}/
    └── chat-frontend-service/
        ├── base/
        └── overlays/{develop,main}/
```

Every service follows the same base/overlay split: the base holds the environment-agnostic shape of the workload, and each overlay supplies the image tag, replica count, ConfigMap, Secret, Ingress, and resource patches for that environment. `develop` targets the `chat-develop` namespace, `main` targets `chat-main`.

## What's in a service

Taking `chat-backend-service` as the example:

**Base** — a Deployment with `replicas: 0` (the overlay sets the real number, so the base is never accidentally deployable), liveness and readiness probes against `/healthz`, CPU/memory requests and limits, config and secrets pulled in via `envFrom`, pod name and node IP injected through the downward API, and a `nodeSelector` pinning it to on-demand capacity. Alongside it, a ClusterIP Service and an HPA scaling on 80% CPU with separate scale-up and scale-down behaviours.

**Overlay** — namespace, image name and tag, replica count, a ConfigMap holding the app's environment variables, an nginx Ingress on `api.v2d27demo.online` (backend) or `chat.v2d27demo.online` (frontend), and strategic-merge patches for resources and HPA bounds.

The frontend follows the same pattern, plus a `topologySpreadConstraints` block that spreads pods across capacity types.

## The image-tag handoff

This is the part that connects CI to CD. The `image: deployment-image` placeholder in each base Deployment is a named target, and the overlay's `kustomization.yml` rewrites it:

```yaml
images:
  - name: deployment-image
    newName: hercules9/real-time-chat-app-backend
    newTag: main-923d4c14c642eeb6e2fcf3a6b506852d02e5999f
```

After building and pushing an image, the GitHub Actions workflow runs `kustomize edit set image` against the right overlay and commits the result back to the repo as a bot. Argo CD sees the commit and syncs. CI never touches the cluster and needs no cluster credentials — the only thing it has is write access to this repository.

## Design decisions worth calling out

- **Replicas ignored on sync.** The service ApplicationSet declares `ignoreDifferences` on `/spec/replicas` for both Deployments, so the HPA can scale freely without Argo CD flagging the cluster as out of sync and reverting it. Without this, GitOps and autoscaling fight each other.
- **Base defaults to zero replicas.** Applying a base directly does nothing, which makes the base/overlay contract explicit rather than conventional.
- **Health probes everywhere.** Both services expose `/healthz` (the frontend's is served by nginx itself), so rollouts wait for readiness instead of shifting traffic to a container that has started but isn't serving.
- **Retention on storage.** The `gp3-resizable` StorageClass uses `reclaimPolicy: Retain`, `allowVolumeExpansion: true`, and `WaitForFirstConsumer` binding, restricted to the two AZs the cluster runs in — so volumes are provisioned where the pod is actually scheduled, can grow in place, and survive a PVC deletion.

## Working with it locally

```bash
# Render what Argo CD would apply
kustomize build services/chat-backend-service/overlays/main

# Apply by hand (only for a scratch cluster — Argo CD owns the real ones)
kubectl apply -k services/chat-backend-service/overlays/main

# Update an image tag the same way CI does
cd services/chat-backend-service/overlays/main
kustomize edit set image deployment-image=hercules9/real-time-chat-app-backend:main-abc1234
```

## Adding a service

1. Create `services/<name>/base/` with `deployment.yml`, `service.yml`, `hpa.yml`, and a `kustomization.yml`.
2. Create `services/<name>/overlays/main/` (and `develop/`) with the ConfigMap, Ingress, patches, and image reference.
3. Commit. The service ApplicationSet's directory glob picks it up on the next reconciliation — no Argo CD change needed.

## Related

- [infras/README.md](../infras/README.md) — the cluster and the Argo CD bootstrap
- [app/README.md](../app/README.md) — the application source these manifests deploy
- [.github/workflows/](../.github/workflows/) — the build-and-commit pipeline
