package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrRelayBlock             = sdkerrors.Register(ModuleName, 1, "relay block failed")
	ErrVerifyProofFail        = sdkerrors.Register(ModuleName, 2, "verify proof failed")
	ErrRelayTooOldBlockHeight = sdkerrors.Register(ModuleName, 3, "relay too old block height")
	ErrResponsePacketOutdated = sdkerrors.Register(ModuleName, 4, "response packet is outdated")
	ErrAppHashNotFound        = sdkerrors.Register(ModuleName, 5, "app hash not found")
)
