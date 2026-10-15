.PHONY: all build-all clean

BUILD_FLAGS := -trimpath -ldflags "-s -w"
export CGO_ENABLED := 0

all: build-all

build-all: clean
	mkdir -p bin
	# AMD64 (x86_64) - v1 for maximum compatibility
	GOOS=linux GOARCH=amd64 GOAMD64=v1 go build $(BUILD_FLAGS) -o bin/yap_x86_64 .
	# ARM64 (aarch64)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/yap_aarch64 .
	# ARMv7 (armv7h)
	GOOS=linux GOARCH=arm GOARM=7 go build $(BUILD_FLAGS) -o bin/yap_armv7h .

clean:
	rm -rf bin/
