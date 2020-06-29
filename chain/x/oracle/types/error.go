package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrOwasmCompilation         = sdkerrors.Register(ModuleName, 1, "owasm compilation failed")
	ErrBadWasmExecution         = sdkerrors.Register(ModuleName, 2, "bad wasm execution")
	ErrDataSourceNotFound       = sdkerrors.Register(ModuleName, 3, "data source not found")
	ErrOracleScriptNotFound     = sdkerrors.Register(ModuleName, 4, "oracle script not found")
	ErrRequestNotFound          = sdkerrors.Register(ModuleName, 5, "request not found")
	ErrRawRequestNotFound       = sdkerrors.Register(ModuleName, 6, "raw request not found")
	ErrReporterNotFound         = sdkerrors.Register(ModuleName, 7, "reporter not found")
	ErrResultNotFound           = sdkerrors.Register(ModuleName, 8, "result not found")
	ErrReporterAlreadyExists    = sdkerrors.Register(ModuleName, 9, "reporter already exists")
	ErrValidatorNotRequested    = sdkerrors.Register(ModuleName, 10, "validator not requested")
	ErrValidatorAlreadyReported = sdkerrors.Register(ModuleName, 11, "validator already reported")
	ErrInvalidReportSize        = sdkerrors.Register(ModuleName, 12, "invalid report size")
	ErrReporterNotAuthorized    = sdkerrors.Register(ModuleName, 13, "reporter not authorized")
	ErrEditorNotAuthorized      = sdkerrors.Register(ModuleName, 14, "editor not authorized")
	ErrReportInfoNotFound       = sdkerrors.Register(ModuleName, 16, "validator report info not found")
	ErrUncompressionFailed      = sdkerrors.Register(ModuleName, 17, "uncompression failed")
	ErrTooLongName              = sdkerrors.Register(ModuleName, 18, "too long name")
	ErrTooLongDescription       = sdkerrors.Register(ModuleName, 19, "too long description")
	ErrEmptyExecutable          = sdkerrors.Register(ModuleName, 20, "empty executable")
	ErrEmptyWasmCode            = sdkerrors.Register(ModuleName, 21, "empty wasm code")
	ErrTooLargeExecutable       = sdkerrors.Register(ModuleName, 22, "too large executable")
	ErrTooLargeWasmCode         = sdkerrors.Register(ModuleName, 23, "too large wasm code")
	ErrInvalidMinCount          = sdkerrors.Register(ModuleName, 24, "invalid min count")
	ErrInvalidAskCount          = sdkerrors.Register(ModuleName, 25, "invalid ask count")
	ErrTooLargeCalldata         = sdkerrors.Register(ModuleName, 26, "too large calldata")
	ErrTooLongClientID          = sdkerrors.Register(ModuleName, 27, "too long client id")
	ErrEmptyReport              = sdkerrors.Register(ModuleName, 28, "empty report")
	ErrDuplicateExternalID      = sdkerrors.Register(ModuleName, 29, "duplicate external id")
	ErrTooLongSchema            = sdkerrors.Register(ModuleName, 30, "too long schema")
	ErrTooLongURL               = sdkerrors.Register(ModuleName, 31, "too long url")
	ErrTooLargeRawReportData    = sdkerrors.Register(ModuleName, 32, "too large raw report data")
	ErrInsufficientValidators   = sdkerrors.Register(ModuleName, 33, "insufficent available validators")
	ErrCreateWithDoNotModify    = sdkerrors.Register(ModuleName, 34, "cannot create with [do-not-modify] content")
	ErrSelfReferenceAsReporter  = sdkerrors.Register(ModuleName, 35, "cannot reference self as reporter")
	ErrOBIDecode                = sdkerrors.Register(ModuleName, 36, "obi decode failed")
)

// WrapMaxError wraps an error message with additional info of the current and max values.
func WrapMaxError(err error, got int, max int) error {
	return sdkerrors.Wrapf(err, "got: %d, max: %d", got, max)
}
