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
	CodeItemDuplication        sdk.CodeType = 204
	CodeItemNotFound           sdk.CodeType = 205
	CodeInvalidState           sdk.CodeType = 206
	CodeBadWasmExecution       sdk.CodeType = 207
)

func ErrInvalidBasicMsg(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidBasicMsg, fmt.Sprintf(format, args...))
}

func ErrBadDataValue(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBadDataValue, fmt.Sprintf(format, args...))
}

func ErrUnauthorizedPermission(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnauthorizedPermission, fmt.Sprintf(format, args...))
}

func ErrItemDuplication(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeItemDuplication, fmt.Sprintf(format, args...))
}

func ErrItemNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeItemNotFound, fmt.Sprintf(format, args...))
}

func ErrInvalidState(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidState, fmt.Sprintf(format, args...))
}

func ErrBadWasmExecution(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeBadWasmExecution, fmt.Sprintf(format, args...))
}
