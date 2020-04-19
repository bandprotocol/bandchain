package owasm

import (
	"fmt"

	"github.com/perlin-network/life/exec"
)

// Executor is the type of any function that supports Owasm execution.
type Executor func(
	env ExecutionEnvironment, code []byte, entry string, calldata []byte, gasLimit uint64,
) (result []byte, gasUsed uint64, err error)

// Execute runs an Owasm script code by via the script's entryID. Note that result will be an
// empty byte slice if the function terminates successfully without saveReturnData getting called.
func Execute(
	env ExecutionEnvironment, code []byte, entry string, calldata []byte, gasLimit uint64,
) (result []byte, gasUsed uint64, err error) {
	resolver := NewResolver(env, calldata)
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
	if err != nil {
		return nil, 0, err
	}
	entryID, ok := vm.GetFunctionExport(entry)
	if !ok {
		return nil, 0, fmt.Errorf("Execute: invalid owasm entry: %s", entry)
	}
	_, err = vm.Run(int(entryID))
	return resolver.result, vm.Gas, err
}
