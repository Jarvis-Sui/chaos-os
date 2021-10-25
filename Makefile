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

build_bin: build_nettc build_proc build_cpu build_mem cp_3rd_bin

build_nettc: exec/bin/nettc/nettc.go exec/bin/nettc/nettc_create.go exec/bin/nettc/util.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/nettc $^

build_proc: exec/bin/process/process.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/process $^

build_cpu: exec/bin/cpu/cpu.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/cpu $^

build_mem: exec/bin/memory/memory.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/memory $^

build_main: main.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_PKG)/$(BIN) $<

cp_3rd_bin:
	cp 3rdparty/bin/* $(BUILD_TARGET_BIN)/

build_on_centos:
	docker build -f "build/Dockerfile" -t go_centos .
	docker run -v "$(shell pwd)":/usr/src/chaos-os -w /usr/src/chaos-os -e GOOS=linux -e GOARCH=amd64 --name chaos-os go_centos make

.PHONY:clean

clean: 
	rm -rf $(BUILD_TARGET)
