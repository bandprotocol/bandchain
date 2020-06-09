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
	wasm, err := Wat2Wasm(code)
	require.NoError(t, err)

	expectedWasm := readWasmFile("test")
	require.Equal(t, expectedWasm, wasm)
}

func TestFailEmptyWatContent(t *testing.T) {
	code := []byte("")
	wasm, err := Wat2Wasm(code)
	require.Equal(t, ErrParseError, err)
	require.Equal(t, []byte(""), wasm)
}

func TestFailInvalidWatContent(t *testing.T) {
	code := []byte("invalid wat content")
	wasm, err := Wat2Wasm(code)
	require.Equal(t, ErrParseError, err)
	require.Equal(t, []byte(""), wasm)
}
