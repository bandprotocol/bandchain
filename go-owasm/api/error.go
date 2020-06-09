package api

import (
	"errors"
)

var (
	ErrCompileFail    = errors.New("compile fail")
	ErrRunFail        = errors.New("run fail")
	ErrParseFail      = errors.New("parse fail")
	ErrEmptyFile      = errors.New("empty file")
	ErrNonUtfResult   = errors.New("non utf-8 result")
	ErrUndefinedError = errors.New("undefined error")
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
		return ErrEmptyFile
	case 5:
		return ErrNonUtfResult
	default:
		return ErrUndefinedError
	}

}
