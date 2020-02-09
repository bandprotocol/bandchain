package owasm

import "github.com/perlin-network/life/exec"

func GetLocalInt64(vm *exec.VirtualMachine, index int) int64 {
	return vm.GetCurrentFrame().Locals[index]
}
