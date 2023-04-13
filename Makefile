# Makefile for building the golang package
REPONAME ?= signoz
IMAGE_NAME ?= golang-distributed-tracing
BACKEND_DOCKER_TAG ?= backend
FRONTEND_DOCKER_TAG ?= frontend

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

LD_FLAGS ?=

.PHONY: dependencies
dependencies:
	go mod tidy
	go mod download

.PHONY: build-binaries
build-binaries: dependencies
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOARCH)/order ./order
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOARCH)/users ./users
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(GOARCH)/payment ./payment

.PHONY: docker-backend
docker-backend:
	docker build -t $(IMAGE_NAME):$(BACKEND_DOCKER_TAG) .

.PHONY: docker-frontend
docker-frontend:
	docker build -t $(IMAGE_NAME):$(FRONTEND_DOCKER_TAG) -f ./frontend/Dockerfile .

.PHONY: docker-backend-push
docker-backend-push:
	docker buildx build --platform linux/arm64,linux/amd64 -f ./Dockerfile --push  -t $(REPONAME)/$(IMAGE_NAME):$(BACKEND_DOCKER_TAG) .

.PHONY: docker-frontend-push
docker-frontend-push:
	docker buildx build --platform linux/arm64,linux/amd64 -f ./frontend/Dockerfile --push  -t $(REPONAME)/$(IMAGE_NAME):$(FRONTEND_DOCKER_TAG) .
