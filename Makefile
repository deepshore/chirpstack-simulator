IMG_NAME ?= schwarzit-int-sit-siot-docker-local.jfrog.io/chirpstack/chirpstack-simulator
IMG_TAG ?= v22

KUBECONFIG ?= $(HOME)/.kube/config

.PHONY: build

build:
	@echo "Building the project..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--push \
		-t $(IMG_NAME):$(IMG_TAG) \
		--build-arg IMG_NAME=$(IMG_NAME) \
		--build-arg IMG_TAG=$(IMG_TAG) \
		.

.PHONY: deploy
deploy:
	@echo "Deploying the project..."
	@kubectl apply -k manifests \
		--kubeconfig=$(KUBECONFIG)
	@kubectl rollout restart deployment/chirpstack-simulator-nextevolution -n chirpstack --kubeconfig=$(KUBECONFIG)