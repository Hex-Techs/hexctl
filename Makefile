all: darwin linux
.PHONY: all

darwin: main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=darwin go build -ldflags "-w -s" -v -o ./bin/hexctl_darwin main.go

linux: main.go
		test -d ./bin || mkdir ./bin
		CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -v -o ./bin/hexctl_linux main.go

upx_darwin: ./bin/hexctl_darwin
		upx -9 -o ./bin/hexctl_darwin_upx ./bin/hexctl_darwin

upx_linux: ./bin/hexctl_linux
		upx -9 -o ./bin/hexctl_linux_upx ./bin/hexctl_linux

clean:
		rm -rf $(shell pwd)/bin
