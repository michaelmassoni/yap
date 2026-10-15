.PHONY: all build build-all test clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD_FLAGS := -trimpath -ldflags "-s -w -X main.version=$(VERSION)"
export CGO_ENABLED := 0

all: test build-all

build:
	go build $(BUILD_FLAGS) -o yap .

build-all: clean
	mkdir -p bin
	# AMD64 (x86_64) - v1 for maximum compatibility
	GOOS=linux GOARCH=amd64 GOAMD64=v1 go build $(BUILD_FLAGS) -o bin/yap_x86_64 .
	# ARM64 (aarch64)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/yap_aarch64 .
	# ARMv7 (armv7h)
	GOOS=linux GOARCH=arm GOARM=7 go build $(BUILD_FLAGS) -o bin/yap_armv7h .

test:
	go vet ./...
	go test ./...

clean:
	rm -rf bin/ yap
