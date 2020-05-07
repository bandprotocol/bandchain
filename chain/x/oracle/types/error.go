package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBasicMsg             = sdkerrors.Register(ModuleName, 1, "invalid basic message")
	ErrBadDataValue                = sdkerrors.Register(ModuleName, 2, "bad data value")
	ErrItemNotFound                = sdkerrors.Register(ModuleName, 3, "item not found")
	ErrBadWasmExecution            = sdkerrors.Register(ModuleName, 4, "bad wasm execution")
	ErrBadDataLength               = sdkerrors.Register(ModuleName, 5, "bad data length")
	ErrDataSourceNotFound          = sdkerrors.Register(ModuleName, 6, "data source not found")
	ErrOracleScriptNotFound        = sdkerrors.Register(ModuleName, 7, "oracle script not found")
	ErrRequestNotFound             = sdkerrors.Register(ModuleName, 8, "request not found")
	ErrRawRequestNotFound          = sdkerrors.Register(ModuleName, 9, "raw request not found")
	ErrReporterNotFound            = sdkerrors.Register(ModuleName, 10, "reporter not found")
	ErrResultNotFound              = sdkerrors.Register(ModuleName, 11, "result not found")
	ErrRawRequestAlreadyExists     = sdkerrors.Register(ModuleName, 12, "raw request already exists")
	ErrReporterAlreadyExists       = sdkerrors.Register(ModuleName, 13, "reporter already exists")
	ErrValidatorNotRequested       = sdkerrors.Register(ModuleName, 14, "validator not requested")
	ErrValidatorAlreadyReported    = sdkerrors.Register(ModuleName, 15, "validator already reported")
	ErrInvalidDataSourceCount      = sdkerrors.Register(ModuleName, 16, "invalid data source count")
	ErrReporterNotAuthorized       = sdkerrors.Register(ModuleName, 17, "reporter not authorized")
	ErrEditorNotAuthorized         = sdkerrors.Register(ModuleName, 18, "editor not authorized")
	ErrValidatorOutOfRange         = sdkerrors.Register(ModuleName, 19, "validator out of range")
	ErrTooManyRawRequests          = sdkerrors.Register(ModuleName, 20, "too many raw requests")
	ErrValidatorReportInfoNotFound = sdkerrors.Register(ModuleName, 21, "validator report info not found")
)
