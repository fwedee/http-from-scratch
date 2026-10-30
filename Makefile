# GO Vars
BINARY := httpserver
MAIN_PACKAGE := ./cmd/httpserver
.DEFAULT_GOAL := build

# Deployment Vars
APP := my-http-server
TAG := $(shell date -u +%Y%m%d%H%M%S)
IMAGE := $(APP):$(TAG)

.PHONY: deploy fmt vet build run test clean fix

create: setup deploy

setup:
	kubectl --context docker-desktop apply -f k8s/namespaces.yaml
	helm install traefik traefik/traefik --namespace traefik --create-namespace
	helm install vmstack vm/victoria-metrics-k8s-stack --namespace monitoring --create-namespace
	helm upgrade --install vts vm/victoria-traces-single --namespace monitoring --kube-context docker-desktop
	kubectl --context docker-desktop apply -f k8s/ingress.yaml -f k8s/ingress2.yaml

deploy:
	docker build -t $(IMAGE) .
	sed 's|__IMAGE__|$(IMAGE)|g' k8s/deployment.yaml | kubectl --context docker-desktop apply -f k8s/service.yaml -f -
	kubectl --context docker-desktop -n demo rollout status deployment/$(APP)

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	CGO_ENABLED=0 go build -o bin/$(BINARY) $(MAIN_PACKAGE)

run:
	CGO_ENABLED=0 go run $(MAIN_PACKAGE)

test:
	go test ./...

clean:
	rm -f bin/$(BINARY)

fix:
	go fix ./...
