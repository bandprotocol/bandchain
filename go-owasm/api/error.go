package api

import (
	"errors"
)

var (
	ErrCompliationFail      = errors.New("compile fail")
	ErrRunFail              = errors.New("run fail")
	ErrParseFail            = errors.New("parse fail")
	ErrWriteBinaryFail      = errors.New("write binary fail")
	ErrResolvesNamesFail    = errors.New("resolve names fail")
	ErrValidateFail         = errors.New("validate fail")
	ErrSpanExceededCapacity = errors.New("span exceeded capacity")
	ErrDeserializeFail      = errors.New("deserialize fail")
	ErrGasCounterInjectFail = errors.New("gas counter inject fail")
	ErrSerializetFail       = errors.New("serialize fail")
	ErrUnknownError         = errors.New("unknown error")
	ErrCompliationError     = errors.New("compile fail")
	ErrRunError             = errors.New("run fail")
	ErrParseError           = errors.New("parse fail")
	ErrWriteBinaryError     = errors.New("write binary fail")
	ErrValidateError        = errors.New("validate fail")
	ErrFunctionNotFound     = errors.New("can't find prepare or execute funtion")
	ErrGasLimitExceeded     = errors.New("gas limit exceeded")
)

// parseError - returns parsed error from errors code on bindings.h
func parseError(code int32) error {
	switch code {
	case 0:
		return nil
	case 1:
		return ErrCompliationFail
	case 2:
		return ErrRunFail
	case 3:
		return ErrParseFail
	case 4:
		return ErrWriteBinaryFail
	case 5:
		return ErrResolvesNamesFail
	case 6:
		return ErrValidateFail
	case 7:
		return ErrSpanExceededCapacity
	case 8:
		return ErrFunctionNotFound
	case 9:
		return ErrGasLimitExceeded
	case 10:
		return ErrSpanExceededCapacity
	case 11:
		return ErrDeserializeFail
	case 12:
		return ErrGasCounterInjectFail
	case 13:
		return ErrSerializetFail
	default:
		return ErrUnknownError
	}

}
