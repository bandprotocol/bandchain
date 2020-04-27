package consuming

import (
	"github.com/bandprotocol/band-consumer/x/consuming/keeper"
	"github.com/bandprotocol/band-consumer/x/consuming/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	NewQuerier    = keeper.NewQuerier
)

type (
	Keeper         = keeper.Keeper
	MsgRequestData = types.MsgRequestData
)
