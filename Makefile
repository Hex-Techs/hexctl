.PHONY: build build-vendor clean

build: main.go
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build -v -o ./bin/n main.go

build-vendor: main.go vendor/
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 go build -mod vendor -v -o ./bin/n main.go

clean:
		rm -rf ./bin
