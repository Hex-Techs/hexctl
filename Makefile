all: darwin linux
.PHONY: all

test:
		gofmt -w .
		go vet ./...
		go test -v ./...

darwin: test main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=darwin go build -ldflags "-w -s" -v -o ./bin/hexctl_darwin main.go

linux: test main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -v -o ./bin/hexctl_linux main.go

clean:
		rm -rf $(shell pwd)/bin
