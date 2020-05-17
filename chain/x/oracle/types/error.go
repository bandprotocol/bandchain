package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBasicMsg          = sdkerrors.Register(ModuleName, 1, "invalid basic message")
	ErrBadDataValue             = sdkerrors.Register(ModuleName, 2, "bad data value")
	ErrItemNotFound             = sdkerrors.Register(ModuleName, 3, "item not found")
	ErrBadWasmExecution         = sdkerrors.Register(ModuleName, 4, "bad wasm execution")
	ErrBadDataLength            = sdkerrors.Register(ModuleName, 5, "bad data length")
	ErrDataSourceNotFound       = sdkerrors.Register(ModuleName, 6, "data source not found")
	ErrOracleScriptNotFound     = sdkerrors.Register(ModuleName, 7, "oracle script not found")
	ErrRequestNotFound          = sdkerrors.Register(ModuleName, 8, "request not found")
	ErrRawRequestNotFound       = sdkerrors.Register(ModuleName, 9, "raw request not found")
	ErrReporterNotFound         = sdkerrors.Register(ModuleName, 10, "reporter not found")
	ErrResultNotFound           = sdkerrors.Register(ModuleName, 11, "result not found")
	ErrRawRequestAlreadyExists  = sdkerrors.Register(ModuleName, 12, "raw request already exists")
	ErrReporterAlreadyExists    = sdkerrors.Register(ModuleName, 13, "reporter already exists")
	ErrValidatorNotRequested    = sdkerrors.Register(ModuleName, 14, "validator not requested")
	ErrValidatorAlreadyReported = sdkerrors.Register(ModuleName, 15, "validator already reported")
	ErrInvalidDataSourceCount   = sdkerrors.Register(ModuleName, 16, "invalid data source count")
	ErrReporterNotAuthorized    = sdkerrors.Register(ModuleName, 17, "reporter not authorized")
	ErrEditorNotAuthorized      = sdkerrors.Register(ModuleName, 18, "editor not authorized")
	ErrValidatorOutOfRange      = sdkerrors.Register(ModuleName, 19, "validator out of range")
	ErrTooManyRawRequests       = sdkerrors.Register(ModuleName, 20, "too many raw requests")

	ErrTooLongName              = sdkerrors.Register(ModuleName, 30, "too long name")
	ErrTooLongDescription       = sdkerrors.Register(ModuleName, 31, "too long description")
	ErrEmptyExecutable          = sdkerrors.Register(ModuleName, 32, "empty executable")
	ErrEmptyWasmCode            = sdkerrors.Register(ModuleName, 33, "empty wasm code")
	ErrTooLargeExecutable       = sdkerrors.Register(ModuleName, 34, "too large executable")
	ErrTooLargeWasmCode         = sdkerrors.Register(ModuleName, 35, "too large wasm code")
	ErrInvalidMinCount          = sdkerrors.Register(ModuleName, 37, "invalid min count")
	ErrAskCountLessThanMinCount = sdkerrors.Register(ModuleName, 38, "ask count < min count")
	ErrTooLargeCalldata         = sdkerrors.Register(ModuleName, 39, "too large calldata")
	ErrTooLongClientID          = sdkerrors.Register(ModuleName, 40, "too long client id")
	ErrEmptyReport              = sdkerrors.Register(ModuleName, 41, "empty report")
	ErrDuplicateExternalID      = sdkerrors.Register(ModuleName, 42, "duplicate external id")
	ErrTooLongSchema            = sdkerrors.Register(ModuleName, 43, "too long schema")
	ErrTooLongURL               = sdkerrors.Register(ModuleName, 44, "too long url")
	ErrTooLargeRawReportData    = sdkerrors.Register(ModuleName, 45, "too large raw report data")
)

// WrapMaxError wraps an error message with additional info of the current and max values.
func WrapMaxError(err error, got int, max int) error {
	return sdkerrors.Wrapf(err, "got: %d, max: %d", got, max)
}
