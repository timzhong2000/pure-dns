.PHONY: \
	all \
	build.linux/386 \
	build.linux/amd64 \
	build.linux/arm \
	build.linux/arm64 \
	build.windows/386 \
	build.windows/amd64 \
	build.windows/arm \
	build.windows/arm64 \
	run \
	format \
	clean \
	help

BINARY="pure-dns"
VERSION="0.1.0"

all: format build.linux/amd64

build.all: \
	build.linux/386 \
	build.linux/amd64 \
	build.linux/arm \
	build.linux/arm64 \
	build.windows/386 \
	build.windows/amd64 \
	build.windows/arm \
	build.windows/arm64

build.linux/386:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_linux_386

build.linux/amd64:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_linux_amd64

build.linux/arm:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_linux_arm

build.linux/arm64:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_linux_arm64

build.windows/386:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_windows_386.exe

build.windows/amd64:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_windows_amd64.exe

build.windows/arm:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=windows GOARCH=arm go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_windows_arm.exe

build.windows/arm64:
	@[ -d build ] || mkdir build
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o ./build/${BINARY}_${VERSION}_windows_arm64.exe

run:
	@go run .

format:
	go fmt .
	go vet .

clean:
	@if [ -d build ] ; then rm -r build ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成x86_64架构适用于linux系统的二进制文件"
	@echo "make build.<platform/arch> - 编译 Go 代码, 生成对应平台架构的二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除build文件夹
	@echo "make format - 运行 Go 工具 'fmt' 和 'vet'"
