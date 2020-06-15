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
	ErrUnknownError         = errors.New("unknown error")
	ErrSpanExceededCapacity = errors.New("span exceeded capacity")
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
		return ErrUnknownError
	case 8:
		return ErrSpanExceededCapacity
	default:
		return ErrUnknownError
	}

}
