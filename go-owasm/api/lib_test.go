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
	spanSize := 1 * 1024 * 1024
	wasm, err := Wat2Wasm(code, spanSize)
	require.NoError(t, err)

	expectedWasm := readWasmFile("test")
	require.Equal(t, expectedWasm, wasm)
}

func TestFailEmptyWatContent(t *testing.T) {
	code := []byte("")
	spanSize := 1 * 1024 * 1024
	_, err := Wat2Wasm(code, spanSize)
	require.Equal(t, ErrParseFail, err)
}

func TestFailInvalidWatContent(t *testing.T) {
	code := []byte("invalid wat content")
	spanSize := 1 * 1024 * 1024
	_, err := Wat2Wasm(code, spanSize)
	require.Equal(t, ErrParseFail, err)
}

func TestFailSpanExceededCapacity(t *testing.T) {
	code := readWatFile("test")
	smallSpanSize := 10
	_, err := Wat2Wasm(code, smallSpanSize)
	require.EqualError(t, err, "span exceeded capacity")
}

func TestFailCompileInvalidContent(t *testing.T) {
	code := []byte("invalid content")
	spanSize := 1 * 1024 * 1024
	_, err := Compile(code, spanSize)
	require.Equal(t, ErrValidateFail, err)
}
