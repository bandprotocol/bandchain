package testapp

import (
	"crypto/sha256"
	"encoding/hex"
)

// An extra Owasm code to test creating or editing oracle scripts.
var WasmExtra1 []byte = wat2wasm([]byte(`
(module
	(type $t0 (func))
	(type $t2 (func (param i64 i64)))
	(import "env" "set_return_data" (func $set_return_data (type $t2)))
	(func $prepare (export "prepare") (type $t0))
	(func $execute (export "execute") (type $t0))
	(memory $memory (export "memory") 17))

`))

// Another extra Owasm code to test creating or editing oracle scripts.
var WasmExtra2 []byte = wat2wasm([]byte(`
(module
	(type $t0 (func))
	(type $t1 (func (param i64 i64 i64 i64)))
	(type $t2 (func (param i64 i64)))
	(import "env" "ask_external_data" (func $ask_external_data (type $t1)))
	(import "env" "set_return_data" (func $set_return_data (type $t2)))
	(func $prepare (export "prepare") (type $t0))
	(func $execute (export "execute") (type $t0))
	(memory $memory (export "memory") 17))
`))

var WasmExtra1FileName string
var WasmExtra2FileName string

func init() {
	wasm1CompiledHash := sha256.Sum256(compile(WasmExtra1))
	wasm2CompiledHash := sha256.Sum256(compile(WasmExtra2))
	WasmExtra1FileName = hex.EncodeToString(wasm1CompiledHash[:])
	WasmExtra2FileName = hex.EncodeToString(wasm2CompiledHash[:])
}
