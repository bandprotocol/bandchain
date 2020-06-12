package api

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lgo_owasm
// #include "bindings.h"
//
// typedef Span (*get_calldata_fn)(env_t*);
// Span cGetCalldata_cgo(env_t *e);
// typedef void (*set_return_data_fn)(env_t*, Span);
// void cSetReturnData_cgo(env_t *e, Span data);
// typedef int64_t (*get_ask_count_fn)(env_t*);
// int64_t cGetAskCount_cgo(env_t *e);
// typedef int64_t (*get_min_count_fn)(env_t*);
// int64_t cGetMinCount_cgo(env_t *e);
// typedef int64_t (*get_ans_count_fn)(env_t*);
// int64_t cGetAnsCount_cgo(env_t *e);
// typedef void (*ask_external_data_fn)(env_t*, int64_t eid, int64_t did);
// void cAskExternalData_cgo(env_t *e, int64_t eid, int64_t did);
// typedef int64_t (*get_external_data_status_fn)(env_t*, int64_t eid, int64_t vid);
// int64_t cGetExternalDataStatus_cgo(env_t *e, int64_t eid, int64_t vid);
// typedef Span (*get_external_data_fn)(env_t*, int64_t eid, int64_t vid);
// Span cGetExternalData_cgo(env_t *e, int64_t eid, int64_t vid);
import "C"
import (
	"unsafe"
)

func Compile(code []byte, spanSize int) ([]byte, error) {
	inputSpan := copySpan(code)
	defer freeSpan(inputSpan)
	outputSpan := newSpan(spanSize)
	defer freeSpan(outputSpan)
	err := parseError(int32(C.do_compile(inputSpan, &outputSpan)))
	return readSpan(outputSpan), err
}

func Prepare(code []byte, env EnvInterface) error {
	return run(code, true, env)
}

func Execute(code []byte, env EnvInterface) error {
	return run(code, false, env)
}

func run(code []byte, isPrepare bool, env EnvInterface) error {
	codeSpan := copySpan(code)
	defer freeSpan(codeSpan)
	envIntl := createEnvIntl(env)
	defer destroyEnvIntl(envIntl)
	return parseError(int32(C.do_run(codeSpan, C.bool(isPrepare), C.Env{
		env: (*C.env_t)(unsafe.Pointer(envIntl)),
		dis: C.EnvDispatcher{
			get_calldata:             C.get_calldata_fn(C.cGetCalldata_cgo),
			set_return_data:          C.set_return_data_fn(C.cSetReturnData_cgo),
			get_ask_count:            C.get_ask_count_fn(C.cGetAskCount_cgo),
			get_min_count:            C.get_min_count_fn(C.cGetMinCount_cgo),
			get_ans_count:            C.get_ans_count_fn(C.cGetAnsCount_cgo),
			ask_external_data:        C.ask_external_data_fn(C.cAskExternalData_cgo),
			get_external_data_status: C.get_external_data_status_fn(C.cGetExternalDataStatus_cgo),
			get_external_data:        C.get_external_data_fn(C.cGetExternalData_cgo),
		},
	})))
}

func Wat2Wasm(code []byte) ([]byte, error) {
	inputSpan := copySpan(code)
	defer freeSpan(inputSpan)
	outputSpan := newSpan(SpanSize)
	defer freeSpan(outputSpan)
	err := parseError(int32(C.do_wat2wasm(inputSpan, &outputSpan)))
	return readSpan(outputSpan), err
}
