.PHONY: all build-all clean

all: build-all

build-all: clean
	mkdir -p bin
	# AMD64 (x86_64) - v1 for maximum compatibility
	GOOS=linux GOARCH=amd64 GOAMD64=v1 go build -ldflags "-s -w" -o bin/yap_x86_64 main.go
	# ARM64 (aarch64)
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/yap_aarch64 main.go
	# ARMv7 (armv7h)
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o bin/yap_armv7h main.go

clean:
	rm -rf bin/
