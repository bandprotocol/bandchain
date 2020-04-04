package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidBasicMsg        = sdkerrors.Register(ModuleName, 1, "")
	ErrBadDataValue           = sdkerrors.Register(ModuleName, 2, "")
	ErrUnauthorizedPermission = sdkerrors.Register(ModuleName, 3, "")
	ErrItemDuplication        = sdkerrors.Register(ModuleName, 4, "")
	ErrItemNotFound           = sdkerrors.Register(ModuleName, 5, "")
	ErrInvalidState           = sdkerrors.Register(ModuleName, 6, "")
	ErrBadWasmExecution       = sdkerrors.Register(ModuleName, 7, "")
)
