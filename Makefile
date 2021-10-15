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

build_nettc: exec/bin/nettc.go exec/bin/nettc_create.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_BIN)/nettc $^

build_main: main.go
	$(GO) build $(GO_FLAGS) -o $(BUILD_TARGET_PKG)/$(BIN) $<

.PHONY:clean

clean: 
	rm -rf $(BUILD_TARGET)
