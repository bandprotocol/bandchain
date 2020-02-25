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
	result, gasUsed, err := Execute(&mockExecutionEnvironment{
		requestID:               1,
		requestedValidatorCount: 2,
	}, code, "execute", []byte{}, 10000)
	require.Nil(t, err)
	require.Equal(t, int64(3), int64(binary.LittleEndian.Uint64(result)))
	require.Equal(t, int64(1056), gasUsed)
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
	require.Equal(t, int64(1279), gasUsed)
}

// DefaultPageSize = 65536 (≈64KB)
// MaxMemoryPages = 1024
// MaxUsageMemory = DefaultPageSize * MaxMemoryPages (≈64MB)
// overflow.wasm is the allocated memory script. It allocates Vec<i64> (i64 ≈ 8 Bytes).
// Ex. Length of Vec<i64> is 5,000,000. It means the Vector will allocate around 38.146 MB. (≈ 611 Pages)
func TestOverflowSuccess(t *testing.T) {
	code, err := ioutil.ReadFile("./res/overflow.wasm")
	require.Nil(t, err)

	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(5000000))

	_, _, err = Execute(&mockExecutionEnvironment{}, code, "execute", size, 100000000000000000)
	require.Nil(t, err)
}

// Ex. Length of Vec<i64> is 8,500,000. It means the Vector will allocate around 64.84 MB. (≈ 1,038 Pages)
func TestOverflowFailWithExceedMemory(t *testing.T) {
	code, err := ioutil.ReadFile("./res/overflow.wasm")
	require.Nil(t, err)

	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(8500000))

	_, _, errExecute := Execute(&mockExecutionEnvironment{}, code, "execute", size, 100000000000000000)
	require.NotNil(t, errExecute)
}
