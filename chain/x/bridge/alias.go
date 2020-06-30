package bridge

import (
	"github.com/bandprotocol/bandchain/chain/x/bridge/keeper"
	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper            = keeper.Keeper
	MsgRelayAndVerify = types.MsgRelayAndVerify
)
