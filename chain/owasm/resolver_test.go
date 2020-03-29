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
	require.Equal(t, err, nil)

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
	require.Equal(t, err, nil)

	//resolveReadExternalData need len(vm.GetCurrentFrame().Locals) >= 5
	localSize := 5
	functionID := -1
	for i := 0; i < len(vm.FunctionCode); i++ {
		if vm.FunctionCode[i].NumParams >= localSize {
			functionID = i
			break
		}
	}
	require.NotEqual(t, functionID, -1)

	// Ignite initializes the first call frame.
	params := make([]int64, 5)
	vm.Ignite(functionID, params...)

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
	require.Equal(t, err, nil)

	//resolveReadExternalData need len(vm.GetCurrentFrame().Locals) >= 5
	localSize := 5
	functionID := -1
	for i := 0; i < len(vm.FunctionCode); i++ {
		if vm.FunctionCode[i].NumParams >= localSize {
			functionID = i
			break
		}
	}
	require.NotEqual(t, functionID, -1)

	// Ignite initializes the first call frame.
	params := make([]int64, 5)
	vm.Ignite(functionID, params...)

	extID := GetLocalInt64(vm, 0)
	valIndex := GetLocalInt64(vm, 1)

	externalData, statusCode, err := resolver.getExternalDataFromCache(extID, valIndex)
	require.Nil(t, err)
	require.Equal(t, []byte("RETURN_DATA"), externalData)
	require.Equal(t, uint8(0), statusCode)
}

func TestGetExternalDataSize(t *testing.T) {

	externalData := "RETURN_DATA"
	env := &mockExecutionEnvironment{
		requestID:                         1,
		requestedValidatorCount:           2,
		externalDataResults:               [][][]byte{{[]byte(externalData), nil}},
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
	require.Equal(t, err, nil)

	//resolveReadExternalData need len(vm.GetCurrentFrame().Locals) >= 5
	localSize := 5
	functionID := -1
	for i := 0; i < len(vm.FunctionCode); i++ {
		if vm.FunctionCode[i].NumParams >= localSize {
			functionID = i
			break
		}
	}
	require.NotEqual(t, functionID, -1)

	// Ignite initializes the first call frame.
	params := make([]int64, 5)
	vm.Ignite(functionID, params...)

	externalDataSize := resolver.resolveGetExternalDataSize(vm)
	require.Equal(t, len(externalData), int(externalDataSize))

}

func TestReadExternalData(t *testing.T) {

	externalData := "RETURN_DATA"
	env := &mockExecutionEnvironment{
		requestID:                         1,
		requestedValidatorCount:           2,
		externalDataResults:               [][][]byte{{[]byte(externalData), nil}},
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
	require.Equal(t, err, nil)

	//resolveReadExternalData need len(vm.GetCurrentFrame().Locals) >= 5
	localSize := 5
	functionID := -1
	for i := 0; i < len(vm.FunctionCode); i++ {
		if vm.FunctionCode[i].NumParams >= localSize {
			functionID = i
			break
		}
	}
	require.NotEqual(t, functionID, -1)

	params := make([]int64, 5)
	vm.Ignite(functionID, params...)

	vm.GetCurrentFrame().Locals[2] = 0
	resultOffset := int(GetLocalInt64(vm, 2))

	vm.GetCurrentFrame().Locals[4] = int64(len(externalData))
	resultSize := int(GetLocalInt64(vm, 4))

	_ = resolver.resolveReadExternalData(vm)

	actualExternalData := vm.Memory[resultOffset : resultOffset+resultSize]

	require.Equal(t, []byte(externalData), actualExternalData)

}
