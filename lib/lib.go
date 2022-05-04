package lib

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi"
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

	// sum.wasm was compiled with TinyGo, which requires being instantiated as a
	// WASI command (to initialize memory).
	// This is required by TinyGo even if the source (../sum/sum.go in this
	// case) doesn't directly use I/O or memory.
	wm, err := wasi.InstantiateSnapshotPreview1(ctx, runtime)
	if err != nil {
		log.Panicln(err)
	}
	defer wm.Close(ctx)

	module, err := runtime.InstantiateModuleFromCode(ctx, sumWASMBytes)
	if err != nil {
		log.Panicln(err)
	}
	defer module.Close(ctx)

	// These are undocumented, but exported. See tinygo-org/tinygo#2788
	malloc := module.ExportedFunction("malloc")
	free := module.ExportedFunction("free")
	var json = `{"Name": "Platypus", "Order": "Monotremata"}`
	jsonSize := uint64(len(json))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, jsonSize)
	if err != nil {
		log.Fatal(err)
	}
	jsonPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, jsonPtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !module.Memory().Write(ctx, uint32(jsonPtr), []byte(json)) {
		log.Fatalf("Memory.Write(%d, %d) out of range of memory size %d",
			jsonPtr, jsonSize, module.Memory().Size(ctx))
	}

	// Now, we can call "parse_name", which reads the string we wrote to memory!
	parse := module.ExportedFunction("parse_name")
	ptrSize, err := parse.Call(ctx, jsonPtr, jsonSize)
	if err != nil {
		log.Fatal(err)
	}

	// Note: This pointer is still owned by TinyGo, so don't try to free it!
	greetingPtr := uint32(ptrSize[0] >> 32)
	greetingSize := uint32(ptrSize[0])
	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := module.Memory().Read(ctx, greetingPtr, greetingSize); !ok {
		log.Fatalf("Memory.Read(%d, %d) out of range of memory size %d",
			greetingPtr, greetingSize, module.Memory().Size(ctx))
	} else {
		fmt.Println("name:", string(bytes))
	}
}
