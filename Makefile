GO=env go

BIN=chaos-os
BUILD_TARGET=target
BUILD_TARGET_PKG=$(BUILD_TARGET)/chaos-os
BUILD_TARGET_BIN=$(BUILD_TARGET_PKG)/bin
GO_FLAGS=-v

build: prebuild build_bin build_main

prebuild:
	rm -rf $(BUILD_TARGET)
	mkdir -p $(BUILD_TARGET)

build_bin: build_nettc

build_nettc: exec/bin/nettc.go exec/bin/nettc_create.go exec/bin/util.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/nettc $^

build_main: main.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_PKG)/$(BIN) $<

build_on_centos:
	docker build -f "build/Dockerfile" -t go_centos .
	docker run -v "$(shell pwd)":/usr/src/chaos-os -w /usr/src/chaos-os -e GOOS=linux -e GOARCH=amd64 --name chaos-os go_centos make

.PHONY:clean

clean: 
	rm -rf $(BUILD_TARGET)
