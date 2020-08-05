.PHONY: build

build: main.go
		test -d ./bin && rm -rf ./bin
		mkdir ./bin
		CGO_ENABLED=0 go build -v -o ./bin/n main.go

build-vendor: main.go
		test -d ./bin && rm -rf ./bin
		mkdir ./bin
		CGO_ENABLED=0 go build -mod vendor -v -o ./bin/n main.go

clean:
		rm -rf ./bin
