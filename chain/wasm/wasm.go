package wasm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func allocateInner(instance wasm.Instance, data []byte) (int64, error) {
	sz := len(data)
	res, err := instance.Exports["__allocate"](sz)
	if err != nil {
		return 0, err
	}
	ptr := res.ToI32()
	mem := instance.Memory.Data()[ptr:]
	if len(mem) < sz {
		return 0, errors.New("allocateInner: invalid memory size")
	}
	copy(mem, data)
	return int64(sz<<32) | int64(ptr), nil
}

func allocate(instance wasm.Instance, data [][]byte) (int64, error) {
	sz := len(data)
	tmp := make([]byte, 8*sz)
	for idx, each := range data {
		loc, err := allocateInner(instance, each)
		if err != nil {
			return 0, err
		}
		binary.LittleEndian.PutUint64(tmp[8*idx:], uint64(loc))
	}

	res, err := instance.Exports["__allocate"](8 * sz)
	if err != nil {
		return 0, err
	}
	ptr := res.ToI32()
	mem := instance.Memory.Data()[ptr:]
	if len(mem) < 8*sz {
		return 0, errors.New("allocate: invalid memory size")
	}
	copy(mem, tmp)
	return int64(sz<<32) | int64(ptr), nil
}

func parseOutput(instance wasm.Instance, ptr int64) ([]byte, error) {
	sz, pointer := int(ptr>>32), (ptr & ((1 << 32) - 1))
	mem := instance.Memory.Data()[pointer:]
	if len(mem) < sz {
		return nil, errors.New("parseOutput: invalid memory size")
	}
	res := make([]byte, sz)
	for idx := 0; idx < sz; idx++ {
		res[idx] = mem[idx]
	}
	return res, nil
}

func storeRawBytes(instance wasm.Instance, rawBytes []byte) (int64, error) {
	return allocateInner(instance, rawBytes)
}

func ParamsInfo(code []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()

	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__params_info")
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func ParseParams(code []byte, params []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeRawBytes(instance, params)
	if err != nil {
		return nil, err
	}

	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__parse_params", paramsInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func RawDataInfo(code []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()

	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__raw_data_info")
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func ParseRawData(code []byte, params []byte, data []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeRawBytes(instance, params)
	dataInput, err := allocateInner(instance, data)
	if err != nil {
		return nil, err
	}
	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__parse_raw_data", paramsInput, dataInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func SerializeParams(code []byte, params []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeRawBytes(instance, params)
	if err != nil {
		return nil, err
	}
	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__serialize_params", paramsInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func ResultInfo(code []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	fn := instance.Exports["__result_info"]
	if fn == nil {
		return nil, errors.New("__result_info not implemented")
	}
	ptr, err := fn()
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func ParseResult(code []byte, result []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	resultInput, err := storeRawBytes(instance, result)
	if err != nil {
		return nil, err
	}
	fn := instance.Exports["__parse_result"]
	if fn == nil {
		return nil, errors.New("__parse_result not implemented")
	}
	ptr, err := fn(resultInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func Prepare(code []byte, params []byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeRawBytes(instance, params)
	if err != nil {
		return nil, err
	}
	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__prepare", paramsInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func executeWithTimeout(limitTimeout time.Duration, inst *wasm.Instance, funcName string, args ...interface{}) (wasm.Value, error) {
	type wasmOutput struct {
		ptr wasm.Value
		err error
	}

	fn := inst.Exports[funcName]
	if fn == nil {
		return wasm.Value{}, errors.New(funcName + " not implemented")
	}

	chanWasmOutput := make(chan wasmOutput, 1)
	go func() {
		ptr, err := fn(args...)
		chanWasmOutput <- wasmOutput{ptr, err}
	}()

	var res wasmOutput
	select {
	case <-time.After(limitTimeout):
		return wasm.Value{}, fmt.Errorf("wasm execution timeout")
	case res = <-chanWasmOutput:
		if res.err != nil {
			return wasm.Value{}, res.err
		} else {
			return res.ptr, res.err
		}
	}
}

func Execute(code []byte, params []byte, inputs [][]byte) ([]byte, error) {
	instance, err := wasm.NewInstance(code)
	if err != nil {
		return nil, err
	}
	defer instance.Close()
	paramsInput, err := storeRawBytes(instance, params)
	if err != nil {
		return nil, err
	}
	wasmInput, err := allocate(instance, inputs)
	if err != nil {
		return nil, err
	}
	ptr, err := executeWithTimeout(100*time.Millisecond, &instance, "__execute", paramsInput, wasmInput)
	if err != nil {
		return nil, err
	}
	return parseOutput(instance, ptr.ToI64())
}

func ReadBytes(filename string) ([]byte, error) {
	return wasm.ReadBytes(filename)
}
