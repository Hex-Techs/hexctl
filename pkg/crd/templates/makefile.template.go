package template

const MakefileTemp = `image := $(IMAGE)

.PHONY: all #            Build image and push image to the harbor, must need IMAGE, example: make all IMAGE=hex-techs/crd:v0.0.0
all: docker-build docker-push

.PHONY: docker-build #   Build the docker image, must need IMAGE, example: make docker-build IMAGE=hex-techs/crd:v0.0.0
docker-build:
    ifndef tag
		@echo "IMAGE must need, please try 'make help' for more information"
		@exit 1
    endif
	docker build -t $(IMAGE) .

.PHONY: docker-push #    Push the docker image, must need IMAGE, example: make docker-push IMAGE=hex-techs/crd:v0.0.0
docker-push:
    ifndef tag
		@echo "IMAGE must need, please try 'make help' for more information"
		@exit 1
    endif
	docker push $(IMAGE)

.PHONY: run #            Run against the configured Kubernetes cluster in ~/.kube/config 
run:
	go run cmd/main.go -v=5 --kubeconfig=$(HOME)/.kube/config

.PHONY: code-generator # Code-generator will update apis
code-generator:
	go mod vendor
	chmod +x $(shell pwd)/vendor/k8s.io/code-generator/generate-groups.sh
	cd ./hack && GOPATH=$(HOME)/go ./update-codegen.sh

.PHONY: help #           help will list all targets
help:
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20`
