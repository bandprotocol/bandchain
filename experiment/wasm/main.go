package main

import (
	"fmt"
	"io/ioutil"

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

func main() {
	fmt.Println("Hello, World!")

	input, err := ioutil.ReadFile("./main.wat")
	if err != nil {
		fmt.Println("ERROR HERE")
		panic(err)
	}

	vm, err := exec.NewVirtualMachine(input, exec.VMConfig{
		DefaultMemoryPages: 64,
		DefaultTableSize:   65536,
		GasLimit:           100000000,
	}, new(Resolver), &compiler.SimpleGasPolicy{GasPerInstruction: 1})
	if err != nil {
		fmt.Println("ERROR HERE2")
		panic(err)
	}

	entryID, _ := vm.GetFunctionExport("infinite")

	fmt.Println("Entry ID:", entryID)

	ret, err := vm.Run(entryID, 40)
	// if err != nil {
	// 	vm.PrintStackTrace()
	// 	panic(err)
	// }

	fmt.Println(vm.Gas)
	fmt.Println(vm.GasLimitExceeded)
	fmt.Println(err)

	fmt.Printf("Return value = %d\n", ret)
}
