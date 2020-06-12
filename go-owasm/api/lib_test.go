package api

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	SpanSize = 1 * 1024 * 1024
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
	_, err := Wat2Wasm(code)
	require.Equal(t, ErrParseError, err)
}

func TestFailInvalidWatContent(t *testing.T) {
	code := []byte("invalid wat content")
	_, err := Wat2Wasm(code)
	require.Equal(t, ErrParseError, err)
}

func TestFailCompileInvalidContent(t *testing.T) {
	code := []byte("invalid content")
	_, err := Compile(code, SpanSize)
	require.Equal(t, ErrValidateError, err)
}
