CLUSTER_NAME := saas-dev
KIND_CONFIG := infra/kind/kind.yaml

INGRESS_NAMESPACE := ingress-nginx

.PHONY: dev-up dev-down cluster-up cluster-down ingress-up status

dev-up: cluster-up ingress-up
	@echo "Local environment is ready"

dev-down: cluster-down

cluster-up:
	@if kind get clusters | grep -qx "$(CLUSTER_NAME)"; then \
		echo "kind cluster $(CLUSTER_NAME) already exists"; \
	else \
		echo "Creating kind cluster $(CLUSTER_NAME)..."; \
		kind create cluster \
			--name $(CLUSTER_NAME) \
			--config $(KIND_CONFIG); \
	fi

cluster-down:
	kind delete cluster --name $(CLUSTER_NAME)

ingress-up:
	@echo "Installing/upgrading ingress-nginx..."
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx \
		--force-update
	helm repo update ingress-nginx
	helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
		--namespace $(INGRESS_NAMESPACE) \
		--create-namespace \
		--set controller.hostPort.enabled=true \
		--set controller.service.type=ClusterIP
	kubectl rollout status deployment/ingress-nginx-controller \
		--namespace $(INGRESS_NAMESPACE) \
		--timeout=120s

status:
	@echo "=== Nodes ==="
	kubectl get nodes
	@echo
	@echo "=== Pods ==="
	kubectl get pods -A
