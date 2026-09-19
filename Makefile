CLUSTER_NAME := saas-dev
KIND_CONFIG := infra/kind/kind.yaml

NAMESPACE := learn-the-grammar
INGRESS_NAMESPACE := ingress-nginx

BACKEND_IMAGE := learn-the-grammar-backend:dev
FRONTEND_IMAGE := learn-the-grammar-frontend:dev


.PHONY: \
	app-run app-stop \
	cluster-up cluster-down \
	ingress-up \
	images-build images-load \
	app-deploy app-delete session-secret \
	status \
	frontend-dev backend-dev \
	unit-test \
	fmt fmt-check lint \
	backend-fmt backend-fmt-check backend-lint \
	frontend-fmt frontend-fmt-check frontend-lint \
	generate generate-check \
	backend-generate frontend-generate \
	e2e


app-run: cluster-up ingress-up images-build images-load app-deploy
	@echo "Local environment is ready at http://localhost"


app-stop: cluster-down


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


images-build:
	@echo "Building backend image..."
	docker build \
		-t $(BACKEND_IMAGE) \
		./backend

	@echo "Building frontend image..."
	docker build \
		-t $(FRONTEND_IMAGE) \
		./frontend


images-load:
	@echo "Loading images into kind..."
	kind load docker-image \
		$(BACKEND_IMAGE) \
		--name $(CLUSTER_NAME)
	kind load docker-image \
		$(FRONTEND_IMAGE) \
		--name $(CLUSTER_NAME)


app-deploy:
	@echo "Deploying application..."
	kubectl apply -f k8s/namespace.yaml
	@$(MAKE) session-secret
	kubectl apply -f k8s/backend-deployment.yaml
	kubectl apply -f k8s/backend-service.yaml
	kubectl apply -f k8s/frontend-deployment.yaml
	kubectl apply -f k8s/frontend-service.yaml
	kubectl apply -f k8s/ingress.yaml
	kubectl rollout restart deployment/backend \
		--namespace $(NAMESPACE)
	kubectl rollout restart deployment/frontend \
		--namespace $(NAMESPACE)

	kubectl rollout status deployment/backend \
		--namespace $(NAMESPACE) \
		--timeout=120s
	kubectl rollout status deployment/frontend \
		--namespace $(NAMESPACE) \
		--timeout=120s


session-secret:
	@if kubectl get secret backend-session --namespace $(NAMESPACE) >/dev/null 2>&1; then \
		echo "Backend session secret already exists"; \
	else \
		echo "Creating backend session secret..."; \
		kubectl create secret generic backend-session \
			--namespace $(NAMESPACE) \
			--from-literal=signing-key="$$(head -c 32 /dev/urandom | base64)"; \
	fi


app-delete:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found


status:
	@echo "=== Nodes ==="
	kubectl get nodes

	@echo
	@echo "=== Application ==="
	kubectl get pods,services,ingress \
		--namespace $(NAMESPACE)

	@echo
	@echo "=== Ingress Controller ==="
	kubectl get pods \
		--namespace $(INGRESS_NAMESPACE)

frontend-dev:
	cd frontend && npm run dev -- --host

backend-dev:
	cd backend && go run ./cmd/server

unit-test:
	cd backend && go test ./...

fmt: backend-fmt frontend-fmt

fmt-check: backend-fmt-check frontend-fmt-check

lint: backend-lint frontend-lint

backend-fmt:
	cd backend && gofmt -s -w .

backend-fmt-check:
	@cd backend && unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not gofmt'ed:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

backend-lint:
	cd backend && golangci-lint run ./...

frontend-fmt:
	cd frontend && npm run format

frontend-fmt-check:
	cd frontend && npm run format:check

frontend-lint:
	cd frontend && npm run lint

generate: backend-generate frontend-generate

backend-generate:
	cd backend && go generate ./...

frontend-generate:
	cd frontend && npm run generate

generate-check: generate
	@changed=$$(git status --porcelain -- backend/internal/api/api.gen.go frontend/src/api/schema.gen.ts); \
	if [ -n "$$changed" ]; then \
		echo "Generated API code is out of date with api/openapi.yaml. Run 'make generate' and commit the result:"; \
		echo "$$changed"; \
		exit 1; \
	fi

e2e:
	cd e2e && npm test
