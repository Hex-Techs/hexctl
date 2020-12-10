.PHONY: build build-vendor clean

build: main.go
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build -v -o ./bin/hexctl main.go

build-vendor: main.go vendor/
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build -mod vendor -v -o ./bin/hexctl main.go

install:
		cp $(shell pwd)/bin/hexctl $(shell go env GOPATH)/bin/

clean:
		rm -rf $(shell pwd)/bin
