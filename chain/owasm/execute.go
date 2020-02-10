package owasm

import (
	"fmt"

	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/exec"
)

type Resolver struct{}

func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	panic(fmt.Errorf("unknown module: %s", module))
}

func (r *Resolver) ResolveGlobal(module, field string) int64 {
	panic(fmt.Errorf("unknown module: %s", module))
}

// Execute runs an Owasm script code by via the script's entryID. Note that
// both result and err can be nil concurrently if the function terminates
// successfully without `saveReturnData` getting called.
func Execute(
	env ExecutionEnvironment,
	code []byte,
	entryID int32,
	calldata []byte,
	gasLimit int64,
) (result []byte, gasUsed int64, err error) {
	vm, err := exec.NewVirtualMachine(code, exec.VMConfig{
		DefaultMemoryPages: 64,
		DefaultTableSize:   65536,
		GasLimit:           uint64(gasLimit),
	}, new(Resolver), &compiler.SimpleGasPolicy{GasPerInstruction: 1})
	if err != nil {
		return nil, 0, err
	}
	_, err = vm.Run(int(entryID))
	// TODO: Implement the way to get result
	return nil, int64(vm.Gas), err
}
