package main

import (
	"github.com/birros/wazero-demo/sum/model"
	"github.com/mailru/easyjson"
	"reflect"
	"unsafe"
)

func main() {}

// _parseName is a WebAssembly export that accepts a string pointer (linear memory
// offset) and returns a pointer/size pair packed into a uint64.
//
// Note: This uses a uint64 instead of two result values for compatibility with
// WebAssembly 1.0.
//export parse_name
func _parseName(ptr, size uint32) (ptrSize uint64) {
	json := ptrToBytes(ptr, size)
	var a model.Animal
	if err := easyjson.Unmarshal(json, &a); err != nil {
		panic(err) // This propagates the error to the Wasm caller.
	}
	ptr, size = stringToPtr(a.Name)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToBytes returns a []byte from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToBytes(ptr, size uint32) (ret []byte) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	hdr.Data = uintptr(ptr)
	hdr.Cap  = uintptr(size) // Tinygo requires these as uintptrs even if they are int fields.
	hdr.Len = uintptr(size) //  ^^
	return
}

// stringToPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
func stringToPtr(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}
