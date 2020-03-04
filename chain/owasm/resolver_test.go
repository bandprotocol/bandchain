package owasm

import (
	"io/ioutil"
	"testing"

	"github.com/perlin-network/life/exec"
	"github.com/stretchr/testify/require"
)

func TestResolveGetCallDataSizet(t *testing.T) {

	env := &mockExecutionEnvironment{
		requestID:               1,
		requestedValidatorCount: 2,
	}
	gasLimit := uint64(10000)
	code, err := ioutil.ReadFile("./res/allocate.wasm")
	callData := []byte("calldata")
	require.Nil(t, err)

	resolver := NewResolver(env, callData)

	vm, err := exec.NewVirtualMachine(code, exec.VMConfig{
		EnableJIT:                false,
		MaxMemoryPages:           1024,
		MaxTableSize:             1024,
		MaxValueSlots:            65536,
		MaxCallStackDepth:        128,
		DefaultMemoryPages:       64,
		DefaultTableSize:         65536,
		GasLimit:                 uint64(gasLimit),
		DisableFloatingPoint:     false,
		ReturnOnGasLimitExceeded: false,
	}, resolver, &BandChainGasPolicy{})

	dataSize := resolver.resolveGetCallDataSize(vm)
	require.Equal(t, dataSize, int64(len(callData)))
}

func TestResolveReadExternalDataSuccess(t *testing.T) {

	env := &mockExecutionEnvironment{
		requestID:                         1,
		requestedValidatorCount:           2,
		externalDataResults:               [][][]byte{{[]byte("RETURN_DATA"), nil}},
		requestExternalDataResultsCounter: [][]int64{{0, 0}, {0, 0}},
	}
	gasLimit := uint64(10000)

	code, err := ioutil.ReadFile("./res/allocate.wasm")
	callData := []byte("calldata")
	require.Nil(t, err)

	resolver := NewResolver(env, callData)

	vm, err := exec.NewVirtualMachine(code, exec.VMConfig{
		EnableJIT:                false,
		MaxMemoryPages:           1024,
		MaxTableSize:             1024,
		MaxValueSlots:            65536,
		MaxCallStackDepth:        128,
		DefaultMemoryPages:       64,
		DefaultTableSize:         65536,
		GasLimit:                 uint64(gasLimit),
		DisableFloatingPoint:     false,
		ReturnOnGasLimitExceeded: false,
	}, resolver, &BandChainGasPolicy{})

	params := [5]int64{}

	// functioncode[69].Numparams is 5
	vm.Ignite(69, params[:5]...)
	dataSize := resolver.resolveGetCallDataSize(vm)
	require.Equal(t, dataSize, int64(len(callData)))

	readExternalStatus := resolver.resolveReadExternalData(vm)
	require.Equal(t, int64(0), readExternalStatus)
}

func TestGetExternalDataFromCacheSuccess(t *testing.T) {

	env := &mockExecutionEnvironment{
		requestID:                         1,
		requestedValidatorCount:           2,
		externalDataResults:               [][][]byte{{[]byte("RETURN_DATA"), nil}},
		requestExternalDataResultsCounter: [][]int64{{0, 0}, {0, 0}},
	}
	gasLimit := uint64(10000)

	code, err := ioutil.ReadFile("./res/allocate.wasm")
	callData := []byte("calldata")
	require.Nil(t, err)

	resolver := NewResolver(env, callData)

	vm, err := exec.NewVirtualMachine(code, exec.VMConfig{
		EnableJIT:                false,
		MaxMemoryPages:           1024,
		MaxTableSize:             1024,
		MaxValueSlots:            65536,
		MaxCallStackDepth:        128,
		DefaultMemoryPages:       64,
		DefaultTableSize:         65536,
		GasLimit:                 uint64(gasLimit),
		DisableFloatingPoint:     false,
		ReturnOnGasLimitExceeded: false,
	}, resolver, &BandChainGasPolicy{})

	params := [5]int64{}

	// functioncode[69].Numparams is 5
	vm.Ignite(69, params[:5]...)

	dataSize := resolver.resolveGetCallDataSize(vm)
	require.Equal(t, dataSize, int64(len(callData)))

	extID := GetLocalInt64(vm, 0)
	valIndex := GetLocalInt64(vm, 1)

	externalData, err := resolver.getExternalDataFromCache(extID, valIndex)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), externalData)

}

func TestSpamGetExternalDataFromCache(t *testing.T) {

	env := &mockExecutionEnvironment{
		requestID:                         1,
		requestedValidatorCount:           2,
		externalDataResults:               [][][]byte{{[]byte("RETURN_DATA"), nil}},
		requestExternalDataResultsCounter: [][]int64{{0, 0}, {0, 0}},
	}
	gasLimit := uint64(10000)

	code, err := ioutil.ReadFile("./res/allocate.wasm")
	callData := []byte("calldata")
	require.Nil(t, err)

	resolver := NewResolver(env, callData)

	vm, err := exec.NewVirtualMachine(code, exec.VMConfig{
		EnableJIT:                false,
		MaxMemoryPages:           1024,
		MaxTableSize:             1024,
		MaxValueSlots:            65536,
		MaxCallStackDepth:        128,
		DefaultMemoryPages:       64,
		DefaultTableSize:         65536,
		GasLimit:                 uint64(gasLimit),
		DisableFloatingPoint:     false,
		ReturnOnGasLimitExceeded: false,
	}, resolver, &BandChainGasPolicy{})

	params := [5]int64{}

	// functioncode[69].Numparams is 5
	vm.Ignite(69, params[:5]...)

	dataSize := resolver.resolveGetCallDataSize(vm)
	require.Equal(t, dataSize, int64(len(callData)))

	extID := GetLocalInt64(vm, 0)
	valIndex := GetLocalInt64(vm, 1)

	externalData, err := resolver.getExternalDataFromCache(extID, valIndex)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), externalData)

	externalData, err = resolver.getExternalDataFromCache(extID, valIndex)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), externalData)

	externalData, err = resolver.getExternalDataFromCache(extID, valIndex)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), externalData)

	require.Equal(t, int64(1), env.requestExternalDataResultsCounter[extID][valIndex])
}
