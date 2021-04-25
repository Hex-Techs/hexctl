export GIT_VERSION = $(shell git describe --tags --always)
export GIT_COMMIT = $(shell git rev-parse HEAD)
export K8S_VERSION = v1.18.6
export CRD_VERSION = v0.1.0

REPO = $(shell go list -m)

GO_BUILD_ARGS = \
	-ldflags " \
	-X '$(REPO)/internal/version.CrdVersion=$(CRD_VERSION)' \
	-X '$(REPO)/internal/version.KubernetesVersion=$(K8S_VERSION)' \
	-X '$(REPO)/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X '$(REPO)/internal/version.GitVersion=$(GIT_VERSION)' \
	" \

.PHONY: build build-vendor clean

build: main.go
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build $(GO_BUILD_ARGS) -v -o ./bin/hexctl main.go

build-vendor: main.go vendor/
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build -mod vendor $(GO_BUILD_ARGS) -v -o ./bin/hexctl main.go

install:
		cp $(shell pwd)/bin/hexctl $(shell go env GOPATH)/bin/

clean:
		rm -rf $(shell pwd)/bin
