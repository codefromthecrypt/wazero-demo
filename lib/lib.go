package lib

import (
	"context"
	_ "embed"
	"log"

	"github.com/tetratelabs/wazero"
)

//go:embed sum.wasm
var sumWASMBytes []byte

func Start() {
	log.SetFlags(log.Lshortfile)

	runtime := wazero.NewRuntime()

	// required even if we don't use  WASI
	wasi, err := runtime.InstantiateModule(wazero.WASISnapshotPreview1())
	if err != nil {
		log.Panicln(err)
	}
	defer wasi.Close()

	module, err := runtime.InstantiateModuleFromSource(sumWASMBytes)
	if err != nil {
		log.Panicln(err)
	}
	defer module.Close()

	sum := module.ExportedFunction("sum")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := sum.Call(module.WithContext(ctx), 30, 12)
	if err != nil {
		log.Panicln(err)
	}

	log.Println(result[0])
}
