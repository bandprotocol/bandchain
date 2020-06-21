package api

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func readWatFile(fileName string) []byte {
	code, err := ioutil.ReadFile(fmt.Sprintf("./../wasm/%s.wat", fileName))
	if err != nil {
		panic(err)
	}
	return code
}

func readWasmFile(fileName string) []byte {
	code, err := ioutil.ReadFile(fmt.Sprintf("./../wasm/%s.wasm", fileName))
	if err != nil {
		panic(err)
	}
	return code
}

func TestSuccessWatToOwasm(t *testing.T) {
	code := readWatFile("test")
	wasm := wat2wasm(code)
	expectedWasm := readWasmFile("test")
	require.Equal(t, expectedWasm, wasm)
}

func TestFailCompileInvalidContent(t *testing.T) {
	code := []byte("invalid content")
	spanSize := 1 * 1024 * 1024
	_, err := Compile(code, spanSize)
	require.Equal(t, ErrValidateFail, err)
}
func TestRuntimeError(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  i32.const 0
		  i32.const 0
		  i32.div_s
		  drop
		)
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

		`))
	code, _ := Compile(wasm, spanSize)

	err := Prepare(code, 100000, NewMockEnv([]byte("")))
	require.Equal(t, ErrRuntimeError, err)
}

func TestInvaildSignature(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(func (param i64 i64 i32 i64)
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))
	  `))
	code, _ := Compile(wasm, spanSize)

	err := Prepare(code, 100000, NewMockEnv([]byte("")))

	require.Equal(t, ErrInvalidSignatureFunction, err)
}

func TestGasLimit(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))
	  `))
	code, err := Compile(wasm, spanSize)
	err = Prepare(code, 100000, NewMockEnv([]byte("")))
	require.NoError(t, err)

	err = Prepare(code, 70000, NewMockEnv([]byte("")))
	require.Equal(t, ErrGasLimitExceeded, err)
}

func TestCompileErrorNoMemory(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `))
	code, err := Compile(wasm, spanSize)
	require.Equal(t, ErrNoMemoryWasm, err)
	require.Equal(t, []uint8([]byte{}), code)
}

func TestCompileErrorMinimumMemoryExceed(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 512)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `))
	_, err := Compile(wasm, spanSize)
	require.NoError(t, err)

	wasm = wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 513)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `))
	_, err = Compile(wasm, spanSize)
	require.Equal(t, ErrNoMemoryWasm, err)
}

func TestCompileErrorSetMaximumMemory(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (set_local $idx (i32.const 0))
		  (block
			  (loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17 20)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `))
	code, err := Compile(wasm, spanSize)
	require.Equal(t, ErrNoMemoryWasm, err)
	require.Equal(t, []uint8([]byte{}), code)
}

func TestCompileErrorCheckWasmImports(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(import "env" "beeb" (func (type 0)))
		(func
		(local $idx i32)
		(set_local $idx (i32.const 0))
		(block
				(loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 1000000000)))
				)
			)
		)
		(func)
		(memory 17)
		(data (i32.const 1048576) "beeb")
		(export "prepare" (func 0))
		(export "execute" (func 1)))
		`))
	code, err := Compile(wasm, spanSize)
	require.Equal(t, ErrCheckWasmImports, err)
	require.Equal(t, []uint8([]byte{}), code)
}

func TestCompileErrorCheckWasmExports(t *testing.T) {
	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm([]byte(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(import "env" "ask_external_data" (func (type 0)))
		(func
		(local $idx i32)
		(set_local $idx (i32.const 0))
		(block
				(loop
				(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (get_local $idx) (i32.const 1000000000)))
				)
			)
		)
		(memory 17)
		(data (i32.const 1048576) "beeb")
		(export "prepare" (func 0)))
		`))
	code, err := Compile(wasm, spanSize)
	require.Equal(t, ErrCheckWasmExports, err)
	require.Equal(t, []uint8([]byte{}), code)
}
