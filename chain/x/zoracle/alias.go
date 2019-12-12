package zoracle

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	EventTypeRequest         = types.EventTypeRequest
	EventTypeReport          = types.EventTypeReport
	EventTypeRequestExecuted = types.EventTypeRequestExecuted
	AttributeKeyRequestID    = types.AttributeKeyRequestID
	AttributeKeyCodeHash     = types.AttributeKeyCodeHash
	AttributeKeyPrepare      = types.AttributeKeyPrepare
	AttributeKeyResult       = types.AttributeKeyResult
	AttributeKeyValidator    = types.AttributeKeyValidator
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewMsgRequest = types.NewMsgRequest
	NewMsgReport  = types.NewMsgReport
)

type (
	Keeper     = keeper.Keeper
	MsgRequest = types.MsgRequest
	MsgReport  = types.MsgReport
)
