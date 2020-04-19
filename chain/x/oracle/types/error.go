package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBasicMsg        = sdkerrors.Register(ModuleName, 1, "InvalidBasicMsg")
	ErrBadDataValue           = sdkerrors.Register(ModuleName, 2, "BadDataValue")
	ErrUnauthorizedPermission = sdkerrors.Register(ModuleName, 3, "UnauthorizedPermission")
	ErrItemDuplication        = sdkerrors.Register(ModuleName, 4, "ItemDuplication")
	ErrItemNotFound           = sdkerrors.Register(ModuleName, 5, "ItemNotFound")
	ErrInvalidState           = sdkerrors.Register(ModuleName, 6, "InvalidState")
	ErrBadWasmExecution       = sdkerrors.Register(ModuleName, 7, "BadWasmExecution")

	ErrBadDataLength        = sdkerrors.Register(ModuleName, 10, "bad data length")
	ErrDataSourceNotFound   = sdkerrors.Register(ModuleName, 11, "data source not found")
	ErrOracleScriptNotFound = sdkerrors.Register(ModuleName, 12, "oracle script not found")
	ErrRawRequestNotFound   = sdkerrors.Register(ModuleName, 13, "raw request not found")

	ErrReporterAlreadyExists = sdkerrors.Register(ModuleName, 30, "reporter already exists")
	ErrReporterNotFound      = sdkerrors.Register(ModuleName, 31, "reporter not found")
)
