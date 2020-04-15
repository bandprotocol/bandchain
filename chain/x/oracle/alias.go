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

	EventTypeCreateDataSource   = types.EventTypeCreateDataSource
	EventTypeEditDataSource     = types.EventTypeEditDataSource
	EventTypeCreateOracleScript = types.EventTypeCreateOracleScript
	EventTypeEditOracleScript   = types.EventTypeEditOracleScript
	EventTypeRequest            = types.EventTypeRequest
	EventTypeReport             = types.EventTypeReport
	EventTypeRequestExecute     = types.EventTypeRequestExecute

	AttributeKeyID            = types.AttributeKeyID
	AttributeKeyRequestID     = types.AttributeKeyRequestID
	AttributeKeyValidator     = types.AttributeKeyValidator
	AttributeKeyResolveStatus = types.AttributeKeyResolveStatus
	AttributeKeyResult        = types.AttributeKeyResult
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	ModuleCdc                   = types.ModuleCdc
	RegisterCodec               = types.RegisterCodec
	NewMsgRequestData           = types.NewMsgRequestData
	NewMsgReportData            = types.NewMsgReportData
	NewMsgCreateOracleScript    = types.NewMsgCreateOracleScript
	NewMsgEditOracleScript      = types.NewMsgEditOracleScript
	NewMsgCreateDataSource      = types.NewMsgCreateDataSource
	NewMsgEditDataSource        = types.NewMsgEditDataSource
	NewMsgAddOracleAddress      = types.NewMsgAddOracleAddress
	NewMsgRemoveOracleAddress   = types.NewMsgRemoveOracleAddress
	NewOracleRequestPacketData  = types.NewOracleRequestPacketData
	NewOracleResponsePacketData = types.NewOracleResponsePacketData

	RequestStoreKey      = types.RequestStoreKey
	ResultStoreKey       = types.ResultStoreKey
	DataSourceStoreKey   = types.DataSourceStoreKey
	OracleScriptStoreKey = types.OracleScriptStoreKey

	NewParams              = types.NewParams
	NewDataSource          = types.NewDataSource
	NewOracleScript        = types.NewOracleScript
	DefaultParams          = types.DefaultParams
	NewRawDataReport       = types.NewRawDataReport
	NewRawDataReportWithID = types.NewRawDataReportWithID

	KeyMaxDataSourceExecutableSize      = types.KeyMaxDataSourceExecutableSize
	KeyMaxOracleScriptCodeSize          = types.KeyMaxOracleScriptCodeSize
	KeyMaxCalldataSize                  = types.KeyMaxCalldataSize
	KeyMaxDataSourceCountPerRequest     = types.KeyMaxDataSourceCountPerRequest
	KeyMaxRawDataReportSize             = types.KeyMaxRawDataReportSize
	KeyMaxResultSize                    = types.KeyMaxResultSize
	KeyMaxNameLength                    = types.KeyMaxNameLength
	KeyMaxDescriptionLength             = types.KeyMaxDescriptionLength
	KeyEndBlockExecuteGasLimit          = types.KeyEndBlockExecuteGasLimit
	KeyGasPerRawDataRequestPerValidator = types.KeyGasPerRawDataRequestPerValidator
	KeyExpirationBlockCount             = types.KeyExpirationBlockCount
	KeyExecuteGas                       = types.KeyExecuteGas
	KeyPrepareGas                       = types.KeyPrepareGas

	QueryRequestByID    = types.QueryRequestByID
	QueryRequests       = types.QueryRequests
	QueryPending        = types.QueryPending
	QueryRequestNumber  = types.QueryRequestNumber
	QueryDataSourceByID = types.QueryDataSourceByID
	QueryDataSources    = types.QueryDataSources
	QueryOracleScripts  = types.QueryOracleScripts

	ParamKeyTable = keeper.ParamKeyTable
)

type (
	Keeper                   = keeper.Keeper
	MsgRequestData           = types.MsgRequestData
	MsgReportData            = types.MsgReportData
	MsgCreateDataSource      = types.MsgCreateDataSource
	MsgEditDataSource        = types.MsgEditDataSource
	MsgCreateOracleScript    = types.MsgCreateOracleScript
	MsgEditOracleScript      = types.MsgEditOracleScript
	MsgAddOracleAddress      = types.MsgAddOracleAddress
	MsgRemoveOracleAddress   = types.MsgRemoveOracleAddress
	OracleRequestPacketData  = types.OracleRequestPacketData
	OracleResponsePacketData = types.OracleResponsePacketData

	RawDataReport         = types.RawDataReport
	RawDataReportWithID   = types.RawDataReportWithID
	RequestQuerierInfo    = types.RequestQuerierInfo
	DataSourceQuerierInfo = types.DataSourceQuerierInfo

	RequestID      = types.RequestID
	OracleScriptID = types.OracleScriptID
	ExternalID     = types.ExternalID
	DataSourceID   = types.DataSourceID

	DataSource    = types.DataSource
	OracleScript  = types.OracleScript
	ResolveStatus = types.ResolveStatus
)
