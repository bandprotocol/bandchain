package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidInput      sdk.CodeType = 101
	CodeInvalidValidator  sdk.CodeType = 102
	CodeRequestNotFound   sdk.CodeType = 103
	CodeResultNotFound    sdk.CodeType = 104
	CodeInvalidOwner      sdk.CodeType = 105
	CodeOutOfReportPeriod sdk.CodeType = 106

	WasmError sdk.CodeType = 105
)

// ErrRequestNotFound is the error for invalid request id
func ErrRequestNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotFound, "request not found")
}

func ErrResultNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeResultNotFound, "result not found")
}

// ErrCodeHashNotFound is the error for invalid code hash
func ErrCodeHashNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "code hash not found")
}

func ErrCodeValidatorNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator invalid")
}

func ErrReportNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotFound, "report not found")
}

func ErrCodeAlreadyExisted(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "code hash already existed")
}

func ErrInvalidOwner(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, "invalid owner")
}

func ErrOutOfReportPeriod(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeOutOfReportPeriod, "report period already ended.")
}
