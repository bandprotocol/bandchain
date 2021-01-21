package oracle

import (
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	PortID            = types.PortID
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper                   = keeper.Keeper
	MsgRequestData           = types.MsgRequestData
	MsgReportData            = types.MsgReportData
	MsgCreateDataSource      = types.MsgCreateDataSource
	MsgEditDataSource        = types.MsgEditDataSource
	MsgCreateOracleScript    = types.MsgCreateOracleScript
	MsgEditOracleScript      = types.MsgEditOracleScript
	MsgActivate              = types.MsgActivate
	MsgAddReporter           = types.MsgAddReporter
	MsgRemoveReporter        = types.MsgRemoveReporter
	OracleRequestPacketData  = types.OracleRequestPacketData
	OracleResponsePacketData = types.OracleResponsePacketData
)
