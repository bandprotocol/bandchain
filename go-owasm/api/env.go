package api

// #include "bindings.h"
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
	ext     EnvInterface
	extData map[[2]int64][]byte
}

func createEnvIntl(ext EnvInterface) *envIntl {
	return &envIntl{
		ext:     ext,
		extData: make(map[[2]int64][]byte),
	}
}

func parseErrorToC(err error) C.Error {
	switch err {
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	default:
		return C.Error_UnknownError
	}
}

//export cGetCalldata
func cGetCalldata(e *C.env_t, calldata *C.Span) C.Error {
	data := (*(*envIntl)(unsafe.Pointer(e))).ext.GetCalldata()
	return writeSpan(calldata, data)
}

//export cSetReturnData
func cSetReturnData(e *C.env_t, span C.Span) C.Error {
	err := (*(*envIntl)(unsafe.Pointer(e))).ext.SetReturnData(readSpan(span))
	if err != nil {
		return parseErrorToC(err)
	}
	return C.Error_NoError
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
func cGetAnsCount(e *C.env_t, val *C.int64_t) C.Error {
	v, err := (*(*envIntl)(unsafe.Pointer(e))).ext.GetAnsCount()
	if err != nil {
		return parseErrorToC(err)
	}
	*val = C.int64_t(v)
	return C.Error_NoError
}

//export cAskExternalData
func cAskExternalData(e *C.env_t, eid C.int64_t, did C.int64_t, span C.Span) C.Error {
	err := (*(*envIntl)(unsafe.Pointer(e))).ext.AskExternalData(int64(eid), int64(did), readSpan(span))
	if err != nil {
		return parseErrorToC(err)
	}
	return C.Error_NoError
}

//export cGetExternalDataStatus
func cGetExternalDataStatus(e *C.env_t, eid C.int64_t, vid C.int64_t, status *C.int64_t) C.Error {
	s, err := (*(*envIntl)(unsafe.Pointer(e))).ext.GetExternalDataStatus(int64(eid), int64(vid))
	if err != nil {
		return parseErrorToC(err)
	}
	*status = C.int64_t(s)
	return C.Error_NoError
}

//export cGetExternalData
func cGetExternalData(e *C.env_t, eid C.int64_t, vid C.int64_t, data *C.Span) C.Error {
	key := [2]int64{int64(eid), int64(vid)}
	env := (*(*envIntl)(unsafe.Pointer(e)))
	if _, ok := env.extData[key]; !ok {
		data, err := env.ext.GetExternalData(int64(eid), int64(vid))
		if err != nil {
			return parseErrorToC(err)
		}
		if data == nil {
			return C.Error_UnavailableExternalDataError
		}
		env.extData[key] = data
	}
	return writeSpan(data, env.extData[key])
}
