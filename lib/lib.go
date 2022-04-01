package lib

import (
	"context"
	_ "embed"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi"
)

//go:embed sum.wasm
var sumWASMBytes []byte

func Start() {
	log.SetFlags(log.Lshortfile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Assign the Go context to the runtime, so it is used during instantiation
	// and any function calls.
	runtime := wazero.NewRuntimeWithConfig(
		wazero.NewRuntimeConfig().WithContext(ctx),
	)

	// sum.wasm was compiled with TinyGo, which requires being instantiated as a
	// WASI command (to initialize memory).
	// This is required by TinyGo even if the source (../sum/sum.go in this
	// case) doesn't directly use I/O or memory.
	wm, err := wasi.InstantiateSnapshotPreview1(runtime)
	if err != nil {
		log.Panicln(err)
	}
	defer wm.Close()

	module, err := runtime.InstantiateModuleFromCode(sumWASMBytes)
	if err != nil {
		log.Panicln(err)
	}
	defer module.Close()

	sum := module.ExportedFunction("sum")

	result, err := sum.Call(nil, 30, 12)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(result[0])
}
