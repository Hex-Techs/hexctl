.PHONY: build

build: main.go
		rm -rf ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=darwin go build -ldflags "-w -s" -v -o ./bin/hexctl_darwin main.go
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -v -o ./bin/hexctl_linux main.go

upx: ./bin/hexctl_darwin ./bin/hexctl_linux
		upx -9 -o ./bin/hexctl_linux_upx ./bin/hexctl_linux
		upx -9 -o ./bin/hexctl_darwin_upx ./bin/hexctl_darwin

clean:
		rm -rf $(shell pwd)/bin
