package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	// TODO: Change this to 101 after old errors were cleared.
	CodeInvalidBasicMsg        sdk.CodeType = 201
	CodeBadDataValue           sdk.CodeType = 202
	CodeUnauthorizedPermission sdk.CodeType = 203

	CodeInvalidInput       sdk.CodeType = 101
	CodeInvalidValidator   sdk.CodeType = 102
	CodeRequestNotFound    sdk.CodeType = 103
	CodeResultNotFound     sdk.CodeType = 104
	CodeInvalidOwner       sdk.CodeType = 105
	CodeOutOfReportPeriod  sdk.CodeType = 106
	CodeDuplicateValidator sdk.CodeType = 107
	CodeDuplicateRequest   sdk.CodeType = 108
	CodeReportNotFound     sdk.CodeType = 109

	WasmError sdk.CodeType = 105
)

// ErrRequestNotFound is the error for invalid request id
func ErrRequestNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotFound, "request not found")
}

func ErrResultNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeResultNotFound, "result not found")
}

func ErrInvalidValidator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator invalid")
}

func ErrReportNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeReportNotFound, "report not found")
}

func ErrCodeAlreadyExisted(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "code hash already existed")
}

func ErrInvalidOwner(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, "invalid owner")
}

func ErrOutOfReportPeriod(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeOutOfReportPeriod, "report period already ended")
}

func ErrDuplicateValidator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateValidator, "duplicate validator")
}

func ErrDuplicateRequest(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateRequest, "duplicate request")
}

func ErrRawDataRequestNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRequestNotFound, "raw data request not found")
}

func ErrInvalidBasicMsg(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidBasicMsg, fmt.Sprintf(format, args...))
}

func ErrBadDataValue(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBadDataValue, fmt.Sprintf(format, args...))
}

func ErrUnauthorizedPermission(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnauthorizedPermission, fmt.Sprintf(format, args...))
}
