all: run

GO_FILES = go.mod go.sum $(shell find . -type f -name '*.go')
LIB_FILES = go.mod go.sum $(shell find sum -type f -name '*.go')

run: build/exe
	./build/exe

build/exe: ${GO_FILES} lib/sum.wasm
	go build -o build/exe main.go

build/ios/arm64/Lib.xcframework: ${GO_FILES} lib/sum.wasm
	gomobile bind -target ios/arm64 -o build/ios/arm64/Lib.xcframework github.com/birros/wazero-demo/lib

build/android/arm64/lib.aar: ${GO_FILES} lib/sum.wasm
	mkdir -p build/android/arm64
	gomobile bind -target android/arm64 -o build/android/arm64/lib.aar github.com/birros/wazero-demo/lib

lib/sum.wasm: ${LIB_FILES}
	GOROOT=$(shell go env GOROOT) tinygo build -o lib/sum.wasm -target wasm sum/sum.go

.PHONY: clean
clean:
	rm -rf build
	find . -type f -name "*.wasm" -exec rm -f {} \;
