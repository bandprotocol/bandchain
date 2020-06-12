package api

import (
	"errors"
)

var (
	ErrCompilationError     = errors.New("compile fail")
	ErrRunError             = errors.New("run fail")
	ErrParseError           = errors.New("parse fail")
	ErrWriteBinaryError     = errors.New("write binary fail")
	ErrResolvesNamesFail    = errors.New("resolve names fail")
	ErrValidateError        = errors.New("validate fail")
	ErrDeserializeFail      = errors.New("deserialize fail")
	ErrGasCounterInjectFail = errors.New("gas counter inject fail")
	ErrSerializetFail      = errors.New("serialize fail")
	ErrUnknownError         = errors.New("unknown error")
)

// parseError - returns parsed error from errors code on bindings.h
func parseError(code int32) error {
	switch code {
	case 0:
		return nil
	case 1:
		return ErrCompilationError
	case 2:
		return ErrRunError
	case 3:
		return ErrParseError
	case 4:
		return ErrWriteBinaryError
	case 5:
		return ErrResolvesNamesFail
	case 6:
		return ErrValidateError
	case 7:
		return ErrUnknownError
	case 9:
		return ErrDeserializeFail
	case 10:
		return ErrGasCounterInjectFail
	case 11:
		return ErrSerializetFail
	default:
		return ErrUnknownError
	}

}
