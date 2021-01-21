package testapp

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	"github.com/GeoDB-Limited/odincore/go-owasm/api"
)

func compile(code []byte) []byte {
	compiled, err := api.Compile(code, types.MaxCompiledWasmCodeSize)
	if err != nil {
		panic(err)
	}
	return compiled
}

func wat2wasm(wat []byte) []byte {
	inputFile, err := ioutil.TempFile("", "input")
	if err != nil {
		panic(err)
	}
	defer os.Remove(inputFile.Name())
	outputFile, err := ioutil.TempFile("", "output")
	if err != nil {
		panic(err)
	}
	defer os.Remove(outputFile.Name())
	if _, err := inputFile.Write(wat); err != nil {
		panic(err)
	}
	if err := exec.Command("wat2wasm", inputFile.Name(), "-o", outputFile.Name()).Run(); err != nil {
		panic(err)
	}
	output, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		panic(err)
	}
	return output
}
