package lib

import (
	"context"
	_ "embed"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

//go:embed sum.wasm
var sumWASMBytes []byte

func Start() {
	log.SetFlags(log.Lshortfile)

	// Assign a Go context used during instantiation and any function calls.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime()
	defer runtime.Close(ctx) // This closes everything this Runtime created.

	// sum.wasm was compiled with TinyGo, which requires being instantiated as a
	// WASI command (to initialize memory).
	// This is required by TinyGo even if the source (../sum/sum.go in this
	// case) doesn't directly use I/O or memory.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		log.Panicln(err)
	}

	module, err := runtime.InstantiateModuleFromBinary(ctx, sumWASMBytes)
	if err != nil {
		log.Panicln(err)
	}

	sum := module.ExportedFunction("sum")

	result, err := sum.Call(ctx, 30, 12)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(result[0])
}
