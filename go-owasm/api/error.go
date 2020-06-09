package api

import (
	"errors"
)

var (
	ErrCompileFail       = errors.New("compile fail")
	ErrRunFail           = errors.New("run fail")
	ErrParseFail         = errors.New("parse fail")
	ErrWriteBinaryFail   = errors.New("write binary fail")
	ErrResolvesNamesFail = errors.New("resolve names fail")
	ErrValidateFail      = errors.New("validate fail")
	ErrUnknownError      = errors.New("unknown error")
)

// parseError - returns parsed error from errors code on bindings.h
func parseError(code int32) error {
	switch code {
	case 1:
		return ErrCompileFail
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
	default:
		return ErrUnknownError
	}

}
