package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidBasicMsg        sdk.CodeType = 101
	CodeBadDataValue           sdk.CodeType = 102
	CodeUnauthorizedPermission sdk.CodeType = 103
	CodeItemDuplication        sdk.CodeType = 104
	CodeItemNotFound           sdk.CodeType = 105
	CodeInvalidState           sdk.CodeType = 106
	CodeBadWasmExecution       sdk.CodeType = 107
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
