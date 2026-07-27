# Docs

Project-level documentation. Component-specific docs live with their code — see [app/README.md](../app/README.md), [infras/README.md](../infras/README.md), and [manifest/README.md](../manifest/README.md).

## Contents

- **[architecture.drawio](architecture.drawio)** — the full system diagram on a single canvas, split into three labelled sections:
  1. **Infrastructure flow** — how Terragrunt provisions the VPC, EKS cluster, and platform add-ons, and how Argo CD gets bootstrapped.
  2. **Deployment flow** — how a push to `app/` becomes a running pod, via GitHub Actions → Docker Hub → a bot commit → Argo CD.
  3. **User flow** — the request path from browser through DNS, NLB, and ingress-nginx to the frontend and backend pods.

  Nodes are colour-coded by type (human actor, AWS resource, Kubernetes object, CI/CD, design note) — the legend is on the canvas.

## Viewing and editing

The `.drawio` file is plain XML and is meant to be edited, not just read:

- **Browser** — open [app.diagrams.net](https://app.diagrams.net) and load the file.
- **VS Code** — install the *Draw.io Integration* extension (`hediet.vscode-drawio`); the file then opens as a diagram in the editor.
- **Desktop** — the [draw.io desktop app](https://github.com/jgraph/drawio-desktop/releases).

GitHub does not render `.drawio` inline. The Mermaid versions of the same three flows are embedded in the [root README](../README.md) for a quick look without opening anything.

## Exporting an image

To produce a raster or vector copy — useful for slides or a portfolio page:

```bash
# From the draw.io desktop app or the drawio CLI
drawio --export --format png --scale 2 --output architecture.png architecture.drawio
drawio --export --format svg --output architecture.svg architecture.drawio
```

In VS Code, the Draw.io Integration extension can save the file as `architecture.drawio.png` or `architecture.drawio.svg`, which keeps the editable XML embedded in the image and renders directly on GitHub.

Keep `architecture.drawio` as the source of truth and regenerate exports from it rather than editing them separately.
