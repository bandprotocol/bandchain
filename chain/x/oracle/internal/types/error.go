package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidInput     sdk.CodeType = 101
	CodeInvalidProvider  sdk.CodeType = 102
	CodeRequestNotFound  sdk.CodeType = 103
	CodeInvalidValidator sdk.CodeType = 104

	WasmError sdk.CodeType = 105
)

// ErrKeyProvidersInvalid is the error for invalid providers
func ErrKeyProvidersInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProvider, "providers invalid format")
}

// ErrRequestNotFound is the error for invalid request id
func ErrRequestNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotFound, "request not found")
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
