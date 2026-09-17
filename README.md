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

## Backend Configuration

The backend loads configuration from three sources, in ascending order of priority:

1. built-in defaults;
2. an optional YAML config file;
3. environment variables.

A field set in a higher-priority source overrides the same field from a lower one. Fields absent from the YAML file simply keep their default value; unrecognized keys in the file are ignored.

| Field  | YAML key | Env var    | Default | Description             |
|--------|----------|------------|---------|--------------------------|
| Port   | `port`   | `LTG_PORT` | `8080`  | TCP port the server listens on |

To point the backend at a YAML config file, use the `-config` flag or the `LTG_CONFIG_FILE` environment variable (the flag takes precedence):

```bash
go run . -config /path/to/config.yaml
```

See [backend/config/config.example.yaml](backend/config/config.example.yaml) for an example config file.

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

## Run tests

To run backend unit tests:

```bash
make unit-test
```

This command will run e2e tests. Cluster should be up and ready for it to sucseed.

```bash
make e2e
```

## Linting and Formatting

The backend uses [`gofmt`](https://pkg.go.dev/cmd/gofmt) for formatting and [`golangci-lint`](https://golangci-lint.run/) for linting. The frontend uses [Prettier](https://prettier.io/) for formatting and [`oxlint`](https://oxc.rs/docs/guide/usage/linter.html) for linting.

Check formatting and linting for both projects:

```bash
make fmt-check
make lint
```

Auto-fix formatting:

```bash
make fmt
```

Individual targets are also available: `make backend-fmt`, `make backend-fmt-check`, `make backend-lint`, `make frontend-fmt`, `make frontend-fmt-check`, `make frontend-lint`.

`golangci-lint` must be installed locally to run `make backend-lint` (or `make lint`) outside of CI. See the [installation guide](https://golangci-lint.run/welcome/install/).

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