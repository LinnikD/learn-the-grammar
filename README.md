# learn-the-grammar
SaaS project for studying language grammar


## Local Development

The project supports two ways of running the application locally:

1. **Development mode** — runs the backend and frontend directly in WSL. This provides the fastest feedback loop and Vite HMR for frontend development.
2. **Kubernetes mode** — builds Docker images and deploys the complete application to a local `kind` cluster. This is closer to the production architecture.

### Prerequisites

The following tools must be installed:

- Docker
- kubectl
- kind
- Helm
- Make
- Go
- Node.js / npm

You can verify the installation with:

```bash
docker --version
kubectl version --client
kind version
helm version
make --version
go version
node --version
npm --version
```

## Development Mode

Development mode runs the Go backend and Vite frontend directly, without Docker or Kubernetes.

Start the backend:

```bash
make backend-dev
```

The backend will be available at:

```text
http://localhost:8080
```

In another terminal, start the frontend:

```bash
make frontend-dev
```

The frontend will be available at:

```text
http://localhost:5173
```

Vite proxies `/api/*` requests to the Go backend:

```text
Browser
   |
   v
Vite :5173
   |
   +-- /*       -> React
   |
   +-- /api/*   -> Go backend :8080
```

Open the application at:

```text
http://localhost:5173
```

This mode is recommended for day-to-day development because frontend changes are immediately available through Vite HMR.

## Kubernetes Mode

Kubernetes mode runs the complete containerized application inside a local `kind` cluster.

The local cluster configuration is located at:

```text
infra/kind/kind.yaml
```

The cluster is named:

```text
saas-dev
```

Ports `80` and `443` are exposed from the kind node to the host.

Start the complete environment:

```bash
make app-run
```

This command:

- creates the `kind` cluster if it does not already exist;
- installs or updates the nginx Ingress Controller;
- builds the backend and frontend Docker images;
- loads the images into the kind cluster;
- deploys the Kubernetes manifests;
- waits for the backend and frontend Deployments to become ready.

The resulting request flow is:

```text
Browser
   |
   v
http://localhost
   |
   v
Ingress
   |
   +-- /api/* -> backend Service -> Go Pod
   |
   +-- /*     -> frontend Service -> nginx Pod -> React
```

Open the application at:

```text
http://localhost
```

To inspect the environment:

```bash
make status
```

To stop the environment and delete the local cluster:

```bash
make app-stop
```

## Additional Commands

The individual steps can also be executed separately:

```bash
make cluster-up
make ingress-up
make images-build
make images-load
make app-deploy
```

Delete only the application namespace:

```bash
make app-delete
```

Delete the entire kind cluster:

```bash
make cluster-down
```