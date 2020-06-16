package api

// #include "bindings.h"
// #include <string.h>
import "C"
import (
	"unsafe"
)

type EnvInterface interface {
	GetCalldata() []byte
	SetReturnData([]byte) error
	GetAskCount() int64
	GetMinCount() int64
	GetAnsCount() (int64, error)
	AskExternalData(eid int64, did int64, data []byte) error
	GetExternalDataStatus(eid int64, vid int64) (int64, error)
	GetExternalData(eid int64, vid int64) ([]byte, error)
}

type envIntl struct {
	ext      EnvInterface
	calldata C.Span
	null     C.Span
	extData  map[[2]int64][]byte
}

func createEnvIntl(ext EnvInterface) *envIntl {
	return &envIntl{
		ext:      ext,
		calldata: copySpan(ext.GetCalldata()),
		null:     copySpan([]byte{}),
		extData:  make(map[[2]int64][]byte),
	}
}

func destroyEnvIntl(e *envIntl) {
	freeSpan(e.calldata)
	freeSpan(e.null)
	// for _, span := range e.extData {
	// 	freeSpan(span)
	// }
}

//export cGetCalldata
func cGetCalldata(e *C.env_t) C.Span {
	return (*(*envIntl)(unsafe.Pointer(e))).calldata
}

//export cSetReturnData
func cSetReturnData(e *C.env_t, span C.Span) C.GoResult {
	err := (*(*envIntl)(unsafe.Pointer(e))).ext.SetReturnData(readSpan(span))
	if err != nil {
		if err == ErrSetReturnDataWrongPeriod {
			return C.SetReturnDataWrongPeriod
		}
		return C.Other
	}
	return C.Ok
}

//export cGetAskCount
func cGetAskCount(e *C.env_t) C.int64_t {
	return C.int64_t((*(*envIntl)(unsafe.Pointer(e))).ext.GetAskCount())
}

//export cGetMinCount
func cGetMinCount(e *C.env_t) C.int64_t {
	return C.int64_t((*(*envIntl)(unsafe.Pointer(e))).ext.GetMinCount())
}

//export cGetAnsCount
func cGetAnsCount(e *C.env_t, val *C.int64_t) C.GoResult {
	v, err := (*(*envIntl)(unsafe.Pointer(e))).ext.GetAnsCount()
	if err != nil {
		if err == ErrAnsCountWrongPeriod {
			return C.AnsCountWrongPeriod
		}
		return C.Other
	}
	*val = C.int64_t(v)
	return C.Ok
}

//export cAskExternalData
func cAskExternalData(e *C.env_t, eid C.int64_t, did C.int64_t, span C.Span) C.GoResult {
	err := (*(*envIntl)(unsafe.Pointer(e))).ext.AskExternalData(int64(eid), int64(did), readSpan(span))
	if err != nil {
		if err == ErrAskExternalDataWrongPeriod {
			return C.AskExternalDataWrongPeriod
		}
		return C.Other
	}
	return C.Ok
}

//export cGetExternalDataStatus
func cGetExternalDataStatus(e *C.env_t, eid C.int64_t, vid C.int64_t, status *C.int64_t) C.GoResult {
	s, err := (*(*envIntl)(unsafe.Pointer(e))).ext.GetExternalDataStatus(int64(eid), int64(vid))
	if err != nil {
		if err == ErrGetExternalDataStatusWrongPeriod {
			return C.GetExternalDataStatusWrongPeriod
		}
		return C.Other
	}
	*status = C.int64_t(s)
	return C.Ok
}

//export cGetExternalData
func cGetExternalData(e *C.env_t, eid C.int64_t, vid C.int64_t, data *C.Span) C.GoResult {
	key := [2]int64{int64(eid), int64(vid)}
	env := (*(*envIntl)(unsafe.Pointer(e)))
	if _, ok := env.extData[key]; !ok {
		data, err := env.ext.GetExternalData(int64(eid), int64(vid))
		if err != nil {
			if err == ErrGetExternalDataWrongPeriod {
				return C.GetExternalDataWrongPeriod
			}
			return C.Other
		}
		if data == nil {
			return C.GetExternalDataFromUnreportedValidator
		}
		env.extData[key] = data
	}
	d := env.extData[key]
	dataLength := len(env.extData[key])
	if int(data.cap) < dataLength {
		return C.SpanExceededCapacity
	}
	C.memcpy(unsafe.Pointer(data.ptr), unsafe.Pointer(&d[0]), C.size_t(dataLength))
	data.len = C.uintptr_t(dataLength)
	return C.Ok
}
