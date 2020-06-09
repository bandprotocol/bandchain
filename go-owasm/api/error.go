package api

import (
	"errors"
)

var (
	ErrCompliationError  = errors.New("compile fail")
	ErrRunError          = errors.New("run fail")
	ErrParseError        = errors.New("parse fail")
	ErrWriteBinaryError  = errors.New("write binary fail")
	ErrResolvesNamesFail = errors.New("resolve names fail")
	ErrValidateError     = errors.New("validate fail")
	ErrUnknownError      = errors.New("unknown error")
)

// parseError - returns parsed error from errors code on bindings.h
func parseError(code int32) error {
	switch code {
	case 1:
		return ErrCompliationError
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
	default:
		return ErrUnknownError
	}

}
