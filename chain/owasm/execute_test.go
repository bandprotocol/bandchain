package owasm

import (
	"encoding/binary"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteCanCallEnv(t *testing.T) {
	code, err := ioutil.ReadFile("./res/main.wasm")
	if err != nil {
		panic(err)
	}
	result, gasUsed, err := Execute(NewMockExecutionEnvironment(), code, "execute", []byte{}, 10000)
	if err != nil {
		panic(err)
	}
	require.Equal(t, int64(42), int64(binary.LittleEndian.Uint64(result)))
	require.Equal(t, int64(22), gasUsed)
}

func TestExecuteOutOfGas(t *testing.T) {
	code, err := ioutil.ReadFile("./res/main.wasm")
	if err != nil {
		panic(err)
	}
	_, _, err = Execute(NewMockExecutionEnvironment(), code, "execute", []byte{}, 10)
	require.EqualError(t, err, "gas limit exceeded")
}
