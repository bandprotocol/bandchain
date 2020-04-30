package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBasicMsg = sdkerrors.Register(ModuleName, 1, "InvalidBasicMsg")
	ErrBadDataValue    = sdkerrors.Register(ModuleName, 2, "BadDataValue")
	ErrItemDuplication = sdkerrors.Register(ModuleName, 4, "ItemDuplication")
	ErrItemNotFound    = sdkerrors.Register(ModuleName, 5, "ItemNotFound")
	ErrInvalidState    = sdkerrors.Register(ModuleName, 6, "InvalidState")

	ErrBadWasmExecution = sdkerrors.Register(ModuleName, 7, "bad wasm execution")

	ErrBadDataLength        = sdkerrors.Register(ModuleName, 10, "bad data length")
	ErrDataSourceNotFound   = sdkerrors.Register(ModuleName, 11, "data source not found")
	ErrOracleScriptNotFound = sdkerrors.Register(ModuleName, 12, "oracle script not found")
	ErrRequestNotFound      = sdkerrors.Register(ModuleName, 13, "request not found")
	ErrRawRequestNotFound   = sdkerrors.Register(ModuleName, 14, "raw request not found")
	ErrReporterNotFound     = sdkerrors.Register(ModuleName, 15, "reporter not found")
	ErrResultNotFound       = sdkerrors.Register(ModuleName, 16, "result not found")

	ErrRawRequestAlreadyExists = sdkerrors.Register(ModuleName, 20, "raw request already exists")
	ErrReporterAlreadyExists   = sdkerrors.Register(ModuleName, 21, "reporter already exists")

	ErrValidatorNotRequested    = sdkerrors.Register(ModuleName, 30, "validator not requested")
	ErrValidatorAlreadyReported = sdkerrors.Register(ModuleName, 31, "validator already reported")

	ErrInvalidDataSourceCount = sdkerrors.Register(ModuleName, 40, "invalid data source count")

	ErrReporterNotAuthorized = sdkerrors.Register(ModuleName, 50, "reporter not authorized")
	ErrEditorNotAuthorized   = sdkerrors.Register(ModuleName, 51, "editor not authorized")

	ErrValidatorOutOfRange = sdkerrors.Register(ModuleName, 60, "validator out of range")
	ErrTooManyRawRequests  = sdkerrors.Register(ModuleName, 61, "too many raw requests")
)
