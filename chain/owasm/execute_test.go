package owasm

import (
	"encoding/binary"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteCanCallEnv(t *testing.T) {
	code, err := ioutil.ReadFile("./res/main.wasm")
	require.Nil(t, err)
	result, gasUsed, err := Execute(&mockExecutionEnvironment{}, code, "execute", []byte{}, 10000)
	if err != nil {
		panic(err)
	}
	require.Equal(t, int64(42), int64(binary.LittleEndian.Uint64(result)))
	require.Equal(t, int64(22), gasUsed)
}

func TestExecuteOutOfGas(t *testing.T) {
	code, err := ioutil.ReadFile("./res/main.wasm")
	require.Nil(t, err)
	_, _, err = Execute(&mockExecutionEnvironment{}, code, "execute", []byte{}, 10)
	require.EqualError(t, err, "gas limit exceeded")
}

func TestExecuteEndToEnd(t *testing.T) {
	code, err := ioutil.ReadFile("./res/silly.wasm")
	require.Nil(t, err)
	env := &mockExecutionEnvironment{
		externalDataResults: [][][]byte{nil, {[]byte("RETURN_DATA")}},
	}

	// It should log "RequestExternalData: DataSourceID = 1, ExternalDataID = 1"
	_, _, err = Execute(env, code, "prepare", []byte{}, 10000)
	require.Nil(t, err)

	// It should return "RETURN_DATA" as the code return data from externalID = 1, validatorIndex = 0
	result, gasUsed, err := Execute(env, code, "execute", []byte{}, 10000)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), result)
	require.Equal(t, int64(1061), gasUsed)
}
