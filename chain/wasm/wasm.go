package wasm

import (
	"errors"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func allocateInner(instance wasm.Instance, data []byte) (int32, error) {
	sz := len(data)
	res, err := instance.Exports["__allocate"](4 + sz)
	if err != nil {
		return 0, err
	}
	ptr := res.ToI32()
	mem := instance.Memory.Data()[ptr:]
	if len(mem) < 4+sz {
		return 0, errors.New("allocateInner: invalid memory size")
	}
	mem[0] = byte(sz % 256)
	mem[1] = byte(sz / 256 % 256)
	mem[2] = byte(sz / 256 / 256 % 256)
	mem[3] = byte(sz / 256 / 256 / 256 % 256)
	for idx, each := range data {
		mem[4+idx] = each
	}
	return ptr, nil
}

func allocate(instance wasm.Instance, data [][]byte) (int32, error) {
	sz := len(data)
	res, err := instance.Exports["__allocate"](4 + 4*sz)
	if err != nil {
		return 0, err
	}
	ptr := res.ToI32()
	mem := instance.Memory.Data()[ptr:]
	if len(mem) < 4+4*sz {
		return 0, errors.New("allocate: invalid memory size")
	}
	mem[0] = byte(sz % 256)
	mem[1] = byte(sz / 256 % 256)
	mem[2] = byte(sz / 256 / 256 % 256)
	mem[3] = byte(sz / 256 / 256 / 256 % 256)
	for idx, each := range data {
		loc, err := allocateInner(instance, each)
		if err != nil {
			return 0, err
		}
		mem[4+4*idx] = byte(loc % 256)
		mem[4+4*idx+1] = byte(loc / 256 % 256)
		mem[4+4*idx+2] = byte(loc / 256 / 256 % 256)
		mem[4+4*idx+3] = byte(loc / 256 / 256 / 256 % 256)
	}
	return ptr, nil
}

func parseOutput(instance wasm.Instance, ptr int32) ([]byte, error) {
	mem := instance.Memory.Data()[ptr:]
	if len(mem) < 4 {
		return nil, errors.New("parseOutput: cannot decode output size")
	}
	sz := 0
	sz += int(mem[0])
	sz += int(mem[1]) * 256
	sz += int(mem[2]) * 256 * 256
	sz += int(mem[3]) * 256 * 256
	if len(mem) < 4+sz {
		return nil, errors.New("parseOutput: invalid memory size")
	}
	res := make([]byte, sz)
	for idx := 0; idx < sz; idx++ {
		res[idx] = mem[4+idx]
	}
	return res, nil
}

func storeParams(instance wasm.Instance, params []byte) (int32, error) {
	return allocateInner(instance, params)
}

func Prepare(code []byte, params []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeParams(instance, params)
	if err != nil {
		return nil, err
	}
	fn := instance.Exports["__prepare"]
	if fn == nil {
		return nil, errors.New("__prepare not implemented")
	}
	ptr, err := fn(paramsInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI32())
}

func Execute(code []byte, params []byte, inputs [][]byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeParams(instance, params)
	if err != nil {
		return nil, err
	}
	wasmInput, err := allocate(instance, inputs)
	if err != nil {
		return nil, err
	}
	fn := instance.Exports["__execute"]
	if fn == nil {
		return nil, errors.New("__execute not implemented")
	}
	ptr, err := fn(paramsInput, wasmInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI32())
}

func ReadBytes(filename string) ([]byte, error) {
	return wasm.ReadBytes(filename)
}
