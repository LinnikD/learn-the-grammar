# learn-the-grammar
SaaS project for studying language grammar


## Local Development

The project runs locally on Kubernetes using kind with Docker as the container runtime.

Prerequisites

Make sure the following tools are installed and available in your PATH:

Docker — container runtime
kubectl — Kubernetes CLI
kind — local Kubernetes cluster
Helm — Kubernetes package manager
Make — local development commands

Docker must be running before starting the local environment.

### Local Kubernetes Cluster

The kind cluster configuration is located at:

`infra/kind/kind.yaml`

The local environment uses a cluster named saas-dev.

The cluster exposes ports 80 and 443 to the host, allowing services to be accessed from the browser via localhost.

### Make Commands
Start the local environment
`make dev-up`

This command:

Creates the saas-dev kind cluster if it does not already exist.
Installs or upgrades the NGINX Ingress Controller.
Waits until the Ingress Controller is ready.

Check the environment
`make status`

Displays Kubernetes nodes and all running pods.

Stop and remove the local environment
`make dev-down`

This deletes the entire saas-dev kind cluster, including all workloads deployed inside it.

The environment can be recreated at any time with:
`make dev-up` 
