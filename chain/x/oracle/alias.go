package oracle

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	PortID            = types.PortID
)

var (
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	ModuleCdc               = types.ModuleCdc
	RegisterCodec           = types.RegisterCodec
	EventTypeRequestExecute = types.EventTypeRequestExecute
)

type (
	Keeper                   = keeper.Keeper
	MsgRequestData           = types.MsgRequestData
	MsgReportData            = types.MsgReportData
	MsgCreateDataSource      = types.MsgCreateDataSource
	MsgEditDataSource        = types.MsgEditDataSource
	MsgCreateOracleScript    = types.MsgCreateOracleScript
	MsgEditOracleScript      = types.MsgEditOracleScript
	MsgAddReporter           = types.MsgAddReporter
	MsgRemoveReporter        = types.MsgRemoveReporter
	OracleRequestPacketData  = types.OracleRequestPacketData
	OracleResponsePacketData = types.OracleResponsePacketData
)
