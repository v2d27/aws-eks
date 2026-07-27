# Infrastructure — AWS EKS with Terragrunt

Everything AWS-side for this project: the VPC, the EKS cluster, and the in-cluster platform components (ingress, storage driver, Argo CD). Written as Terraform modules and wired together with Terragrunt so each piece is a separate state file with explicit dependencies between them.

> If you only want to run the chat app locally with Docker, you can skip this folder entirely — see [app/README.md](../app/README.md).

## Why Terragrunt

The alternative was one large Terraform root module for the whole environment. Splitting it into Terragrunt units gives three things that matter in practice:

- **Small blast radius** — a change to the ingress controller plans against ingress state only, not the whole cluster.
- **Explicit dependency graph** — `dependency` blocks pass real outputs (VPC ID, cluster endpoint, CA data) between units, and `terragrunt run-all apply` walks them in the right order.
- **DRY configuration** — system name, region, account ID, domain, and tags are defined once and inherited down the tree, so adding a second environment means copying a folder and changing `env.hcl`.

## Layout

```text
infras/
├── root.hcl               # S3 remote state + DynamoDB locking, AWS provider generation
├── global_vars.hcl        # System name, root domain, allowed admin IPs
├── account/
│   ├── account.hcl        # AWS account ID
│   └── ap-southeast-1/
│       ├── region.hcl     # Region
│       └── develop/
│           ├── env.hcl    # Merges global/region/account vars, defines the common tag set
│           ├── vpc/
│           ├── iam/current-user/
│           └── eks/
│               ├── eks-cluster/
│               ├── eks-ebs-csi/
│               ├── eks-ingress-nginx/
│               │   └── security-groups/
│               ├── eks-argocd/
│               │   └── argocd-admin-pass/
│               └── eks-argocd-manage/
└── module/                # Terraform modules the units point at
    ├── vpc/               # terraform-aws-modules/vpc (vendored)
    ├── eks/               # terraform-aws-modules/eks (vendored)
    ├── ws-eks-helm/       # Thin helm_release wrapper with templated values files
    ├── ws-sg/             # Security group wrapper
    ├── ws-iam/iam/        # Resolves the caller identity for cluster access entries
    └── random_password/   # Stable generated password + bcrypt hash for Argo CD admin
```

The `account / region / environment` nesting is the convention that makes this scale — a second region or a `production` environment slots in without touching any existing unit.

## What gets built

- **`vpc`** — `10.0.0.0/16` VPC with four `/20` public subnets across two AZs, DNS hostnames and support enabled.
- **`iam/current-user`** — reads the caller's ARN so whoever applies gets cluster-admin access automatically.
- **`eks/eks-cluster`** — EKS 1.33 with a public endpoint and two managed node groups (on-demand and spot, `t3.medium`, scaling 0–2).
- **`eks/eks-ebs-csi`** — the `aws-ebs-csi-driver` Helm chart, so PVCs can bind to EBS volumes.
- **`eks/eks-ingress-nginx`** — `ingress-nginx` behind an internet-facing NLB, autoscaled 2–4 replicas.
- **`eks/eks-ingress-nginx/security-groups`** — the NLB security group: 80/443 open, port 5000 restricted to the IPs listed in `global_vars.hcl`.
- **`eks/eks-argocd`** — the Argo CD Helm chart with ingress, RBAC policy, and a dedicated `terraform-admin` API account.
- **`eks/eks-argocd/argocd-admin-pass`** — generates the admin password once and exposes both the plaintext and the bcrypt hash.
- **`eks/eks-argocd-manage`** — registers the AppProject, repository credentials, and the root Application that bootstraps everything under `manifest/`.

Cluster add-ons enabled on the control plane: `vpc-cni` (with prefix delegation for higher pod density per node), `coredns`, `kube-proxy`, and `eks-pod-identity-agent`.

### A note on the Argo CD bootstrap

`eks-argocd` installs Argo CD; `eks-argocd-manage` then uses the [Argo CD Terraform provider](https://registry.terraform.io/modules/v2d27/argocd/argocd/1.0.1) to create the single root Application pointing at `manifest/argocd-appsets/overlays/main`. From that point on Argo CD manages itself and the rest of the cluster from Git — Terraform is no longer in the deployment path for application changes. The admin password comes from the `random_password` module with a fixed `keepers` value, so the bcrypt hash the Helm chart receives and the plaintext the provider authenticates with stay in sync across applies.

## Prerequisites

Version-pinned via [tgswitch](https://github.com/warrensbox/tgswitch) and [tfswitch](https://github.com/warrensbox/terraform-switcher):

```sh
curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | sudo bash
tgswitch 0.82.4

curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | sudo bash
tfswitch 1.12.2
```

You also need the AWS CLI configured with credentials for the target account, plus `kubectl` and `helm` for anything post-apply.

## Before your first apply

1. Set your AWS account ID in `account/account.hcl`.
2. Set your domain and the IPs allowed through the NLB on port 5000 in `global_vars.hcl`.
3. Create the S3 state bucket and DynamoDB lock table named in `root.hcl`, or point it at ones you already have.
4. Optionally add extra cluster admins to `admin_users` in `eks/eks-cluster/terragrunt.hcl`.

## Usage

```sh
# One unit
cd account/ap-southeast-1/develop/vpc
terragrunt plan
terragrunt apply

# Whole environment, in dependency order
cd account/ap-southeast-1/develop
terragrunt run-all plan
terragrunt run-all apply

# Get kubeconfig once the cluster is up
aws eks update-kubeconfig --region ap-southeast-1 --name wsaws-develop-eks
```

Tear down in reverse: `terragrunt run-all destroy` from the environment folder.

## Cost notes

This is a demo environment and is sized like one — `t3.medium` nodes, node groups that scale to zero, one spot group alongside the on-demand group, no NAT gateway (workloads sit in public subnets with security groups doing the gating), and CloudWatch log groups for the control plane disabled. Real production would want private subnets with NAT, larger nodes, and control-plane logging on.

## Related

- [manifest/README.md](../manifest/README.md) — what Argo CD deploys into this cluster
- [app/README.md](../app/README.md) — the application itself
