package owasm

import (
	"fmt"

	"github.com/perlin-network/life/exec"
)

type resolver struct {
	env      ExecEnv
	calldata []byte
	result   []byte
}

func (r *resolver) ResolveFunc(module, field string) exec.FunctionImport {
	if module != "env" {
		panic(fmt.Errorf("ResolveFunc: unknown module: %s", module))
	}
	switch field {
	case "getRequestedValidatorCount":
		return r.resolveGetRequestedValidatorCount
	case "getSufficientValidatorCount":
		return r.resolveGetSufficientValidatorCount
	case "getReceivedValidatorCount":
		return r.resolveGetReceivedValidatorCount
	case "getPrepareBlockTime":
		return r.resolveGetPrepareBlockTime
	case "getAggregateBlockTime":
		return r.resolveGetAggregateBlockTime
	case "readValidatorAddress":
		return r.resolveReadValidatorAddress
	case "getCallDataSize":
		return r.resolveGetCallDataSize
	case "readCallData":
		return r.resolveReadCallData
	case "saveReturnData":
		return r.resolveSaveReturnData
	case "requestExternalData":
		return r.resolveRequestExternalData
	case "getExternalDataStatusCode":
		return r.resolveGetExternalDataStatusCode
	case "getExternalDataSize":
		return r.resolveGetExternalDataSize
	case "readExternalData":
		return r.resolveReadExternalData
	default:
		panic(fmt.Errorf("ResolveFunc: unknown field: %s", field))
	}
}

func (r *resolver) ResolveGlobal(module, field string) int64 {
	panic(fmt.Errorf("ResolveGlobal is not supported by owasm!"))
}

func (r *resolver) resolveGetRequestedValidatorCount(vm *exec.VirtualMachine) int64 {
	return r.env.GetRequestedValidatorCount()
}

func (r *resolver) resolveGetSufficientValidatorCount(vm *exec.VirtualMachine) int64 {
	return r.env.GetSufficientValidatorCount()
}

func (r *resolver) resolveGetReceivedValidatorCount(vm *exec.VirtualMachine) int64 {
	return r.env.GetRequestedValidatorCount()
}

func (r *resolver) resolveGetPrepareBlockTime(vm *exec.VirtualMachine) int64 {
	return r.env.GetPrepareBlockTime()
}

func (r *resolver) resolveGetAggregateBlockTime(vm *exec.VirtualMachine) int64 {
	return r.env.GetAggregateBlockTime()
}

func (r *resolver) resolveReadValidatorAddress(vm *exec.VirtualMachine) int64 {
	validatorIndex := GetLocalInt64(vm, 0)
	resultOffset := int(GetLocalInt64(vm, 0))
	address, err := r.env.GetValidatorAddress(validatorIndex)
	if err != nil {
		return -1
	}
	copy(vm.Memory[resultOffset:resultOffset+len(address)], address)
	return 0
}

func (r *resolver) resolveGetCallDataSize(vm *exec.VirtualMachine) int64 {
	return int64(len(r.calldata))
}

func (r *resolver) resolveReadCallData(vm *exec.VirtualMachine) int64 {
	resultOffset := int(GetLocalInt64(vm, 0))
	seekOffset := int(GetLocalInt64(vm, 1))
	resultSize := int(GetLocalInt64(vm, 2))
	copy(vm.Memory[resultOffset:resultOffset+resultSize], r.calldata[seekOffset:seekOffset+resultSize])
	return 0
}

func (r *resolver) resolveSaveReturnData(vm *exec.VirtualMachine) int64 {
	dataOffset := int(GetLocalInt64(vm, 0))
	dataLength := int(GetLocalInt64(vm, 1))
	if dataLength > int(r.env.GetMaximumResultSize()) {
		return -1
	}
	r.result = make([]byte, dataLength)
	copy(r.result, vm.Memory[dataOffset:dataOffset+dataLength])
	return 0
}

func (r *resolver) resolveRequestExternalData(vm *exec.VirtualMachine) int64 {
	dataSourceID := GetLocalInt64(vm, 0)
	externalID := GetLocalInt64(vm, 1)
	dataOffset := int(GetLocalInt64(vm, 2))
	dataLength := int(GetLocalInt64(vm, 3))
	if dataLength > int(r.env.GetMaximumCalldataOfDataSourceSize()) {
		return -1
	}
	data := make([]byte, dataLength)
	copy(data, vm.Memory[dataOffset:dataOffset+dataLength])
	err := r.env.RequestExternalData(dataSourceID, externalID, data)
	if err != nil {
		return -1
	}
	return 0
}

func (r *resolver) resolveGetExternalDataStatusCode(vm *exec.VirtualMachine) int64 {
	externalID := GetLocalInt64(vm, 0)
	validatorIndex := GetLocalInt64(vm, 1)
	_, statusCode, err := r.env.GetExternalData(externalID, validatorIndex)
	if err != nil {
		return -1
	}
	return int64(statusCode)
}

func (r *resolver) resolveGetExternalDataSize(vm *exec.VirtualMachine) int64 {
	externalID := GetLocalInt64(vm, 0)
	validatorIndex := GetLocalInt64(vm, 1)
	externalData, _, err := r.env.GetExternalData(externalID, validatorIndex)
	if err != nil {
		return -1
	}
	return int64(len(externalData))
}

func (r *resolver) resolveReadExternalData(vm *exec.VirtualMachine) int64 {
	externalID := GetLocalInt64(vm, 0)
	validatorIndex := GetLocalInt64(vm, 1)
	resultOffset := int(GetLocalInt64(vm, 2))
	seekOffset := int(GetLocalInt64(vm, 3))
	resultSize := int(GetLocalInt64(vm, 4))
	externalData, _, err := r.env.GetExternalData(externalID, validatorIndex)
	if err != nil {
		return -1
	}
	copy(vm.Memory[resultOffset:resultOffset+resultSize], externalData[seekOffset:seekOffset+resultSize])
	return 0
}

func NewResolver(env ExecEnv, calldata []byte) *resolver {
	return &resolver{
		env:      env,
		calldata: calldata,
		result:   []byte{},
	}
}
