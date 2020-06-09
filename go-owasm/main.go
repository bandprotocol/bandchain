package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bandprotocol/bandchain/go-owasm/api"
)

type Env struct{}

func (e *Env) GetCalldata() []byte {
	return []byte("switza")
}

func (e *Env) SetReturnData(data []byte) {
	fmt.Println("set", string(data))
}

func (e *Env) GetAskCount() int64 {
	return 10000
}

func (e *Env) GetMinCount() int64 {
	return 20000
}

func (e *Env) GetAnsCount() int64 {
	return 30000
}

func (e *Env) AskExternalData(eid int64, did int64, data []byte) {
	fmt.Println("asked", eid, did, string(data))
}

func (e *Env) GetExternalDataStatus(eid int64, vid int64) int64 {
	return 42
}

func (e *Env) GetExternalData(eid int64, vid int64) []byte {
	return []byte("switez")
}

func WatToWasm(fileName string) error {
	code, _ := ioutil.ReadFile(fmt.Sprintf("./wasm/%s.wat", fileName))
	wasm, err := api.WatToWasm(code)
	if err != nil {
		panic(err)
	}

	fmt.Println("wasm:", wasm)
	f, err := os.Create(fmt.Sprintf("./wasm/%s.wasm", fileName))
	defer f.Close()
	n, err := f.Write(wasm)
	fmt.Printf("wrote %d bytes\n", n)
	if err != nil {
		return err
	}
	f.Sync()

	return nil
}

func main() {
	// fmt.Println("Hello, World!")
	// code, _ := ioutil.ReadFile("./wasm/test.wat")
	// wasm, e := api.WatToWasm(code)
	// fmt.Println("wasm", wasm)
	// fmt.Println(e)

	WatToWasm("test")
}
