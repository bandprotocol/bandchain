package api

// #include "bindings.h"
import "C"
import (
	"unsafe"
)

type EnvInterface interface {
	GetCalldata() []byte
	SetReturnData([]byte)
	GetAskCount() int64
	GetMinCount() int64
	GetAnsCount() int64
	AskExternalData(eid int64, did int64, data []byte)
	GetExternalDataStatus(eid int64, vid int64) int64
	GetExternalData(eid int64, vid int64) []byte
}

type envIntl struct {
	ext      EnvInterface
	calldata C.Span
	null     C.Span
	extData  map[[2]int64]C.Span
}

func createEnvIntl(ext EnvInterface) *envIntl {
	return &envIntl{
		ext:      ext,
		calldata: copySpan(ext.GetCalldata()),
		null:     copySpan([]byte{}),
		extData:  make(map[[2]int64]C.Span),
	}
}

func destroyEnvIntl(e *envIntl) {
	freeSpan(e.calldata)
	freeSpan(e.null)
	for _, span := range e.extData {
		freeSpan(span)
	}
}

//export cGetCalldata
func cGetCalldata(e *C.env_t) C.Span {
	return (*(*envIntl)(unsafe.Pointer(e))).calldata
}

//export cSetReturnData
func cSetReturnData(e *C.env_t, span C.Span) {
	(*(*envIntl)(unsafe.Pointer(e))).ext.SetReturnData(readSpan(span))
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
func cGetAnsCount(e *C.env_t) C.int64_t {
	return C.int64_t((*(*envIntl)(unsafe.Pointer(e))).ext.GetAnsCount())
}

//export cAskExternalData
func cAskExternalData(e *C.env_t, eid C.int64_t, did C.int64_t, span C.Span) {
	(*(*envIntl)(unsafe.Pointer(e))).ext.AskExternalData(int64(eid), int64(did), readSpan(span))
}

//export cGetExternalDataStatus
func cGetExternalDataStatus(e *C.env_t, eid C.int64_t, vid C.int64_t) C.int64_t {
	return C.int64_t((*(*envIntl)(unsafe.Pointer(e))).ext.GetExternalDataStatus(int64(eid), int64(vid)))
}

//export cGetExternalData
func cGetExternalData(e *C.env_t, eid C.int64_t, vid C.int64_t) C.Span {
	key := [2]int64{int64(eid), int64(vid)}
	env := (*(*envIntl)(unsafe.Pointer(e)))
	if _, ok := env.extData[key]; !ok {
		data := env.ext.GetExternalData(int64(eid), int64(vid))
		if data == nil {
			return env.null
		}
		env.extData[key] = copySpan(env.ext.GetExternalData(int64(eid), int64(vid)))
	}
	return env.extData[key]
}
