package owasm

import (
	"fmt"

	"github.com/perlin-network/life/exec"
)

type resolver struct {
	env      ExecutionEnvironment
	calldata []byte
	result   []byte
}

func (r *resolver) ResolveFunc(module, field string) exec.FunctionImport {
	if module != "env" {
		panic(fmt.Errorf("ResolveFunc: unknown module: %s", module))
	}
	switch field {
	case "getCurrentRequestID":
		return r.resolveGetCurrentRequestID
	case "getRequestedValidatorCount":
		return r.resolveGetRequestedValidatorCount
	case "getReceivedValidatorCount":
		return r.resolveGetReceivedValidatorCount
	case "getPrepareBlockTime":
		return r.resolveGetPrepareBlockTime
	case "getAggregateBlockTime":
		return r.resolveGetAggregateBlockTime
	case "readValidatorPubKey":
		return r.resolveReadValidatorPubKey
	case "getCallDataSize":
		return r.resolveGetCallDataSize
	case "readCallData":
		return r.resolveReadCallData
	case "saveReturnData":
		return r.resolveSaveReturnData
	case "requestExternalData":
		return r.resolveRequestExternalData
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

func (r *resolver) resolveGetCurrentRequestID(vm *exec.VirtualMachine) int64 {
	return r.env.GetCurrentRequestID()
}

func (r *resolver) resolveGetRequestedValidatorCount(vm *exec.VirtualMachine) int64 {
	return r.env.GetRequestedValidatorCount()
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

func (r *resolver) resolveReadValidatorPubKey(vm *exec.VirtualMachine) int64 {
	validatorIndex := GetLocalInt64(vm, 0)
	resultOffset := int(GetLocalInt64(vm, 0))
	pubkey, err := r.env.GetValidatorPubKey(validatorIndex)
	if err != nil {
		panic(err)
	}
	copy(vm.Memory[resultOffset:resultOffset+len(pubkey)], pubkey)
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
	r.result = make([]byte, dataLength)
	copy(r.result, vm.Memory[dataOffset:dataOffset+dataLength])
	return 0
}

func (r *resolver) resolveRequestExternalData(vm *exec.VirtualMachine) int64 {
	panic("resolveRequestExternalData: not implemented!")
}

func (r *resolver) resolveGetExternalDataSize(vm *exec.VirtualMachine) int64 {
	panic("resolveGetExternalDataSize: not implemented!")
}

func (r *resolver) resolveReadExternalData(vm *exec.VirtualMachine) int64 {
	panic("resolveReadExternalData: not implemented!")
}

func NewResolver(env ExecutionEnvironment, calldata []byte) *resolver {
	return &resolver{
		env:      env,
		calldata: calldata,
	}
}
