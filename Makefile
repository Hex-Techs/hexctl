all: darwin linux
.PHONY: all

test:
		gofmt -w .
		golint ./...
		gocyclo -avg .
		go vet ./...
		go test -v ./...

darwin: test main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=darwin go build -ldflags "-w -s" -v -o ./bin/hexctl_darwin main.go

linux: test main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -v -o ./bin/hexctl_linux main.go

upx_darwin: ./bin/hexctl_darwin
		upx -9 -o ./bin/hexctl_darwin_upx ./bin/hexctl_darwin

upx_linux: ./bin/hexctl_linux
		upx -9 -o ./bin/hexctl_linux_upx ./bin/hexctl_linux

clean:
		rm -rf $(shell pwd)/bin
