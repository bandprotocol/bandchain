package zoracle

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey

	EventTypeRequest      = types.EventTypeRequest
	EventTypeReport       = types.EventTypeReport
	AttributeKeyRequestID = types.AttributeKeyRequestID
	AttributeKeyValidator = types.AttributeKeyValidator
)

var (
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier
	ModuleCdc                = types.ModuleCdc
	RegisterCodec            = types.RegisterCodec
	NewMsgRequestData        = types.NewMsgRequestData
	NewMsgReportData         = types.NewMsgReportData
	NewMsgCreateOracleScript = types.NewMsgCreateOracleScript
	NewMsgEditOracleScript   = types.NewMsgEditOracleScript
	NewMsgCreateDataSource   = types.NewMsgCreateDataSource
	NewMsgEditDataSource     = types.NewMsgEditDataSource

	RequestStoreKey      = types.RequestStoreKey
	ResultStoreKey       = types.ResultStoreKey
	DataSourceStoreKey   = types.DataSourceStoreKey
	OracleScriptStoreKey = types.OracleScriptStoreKey

	NewParams        = types.NewParams
	NewDataSource    = types.NewDataSource
	NewOracleScript  = types.NewOracleScript
	DefaultParams    = types.DefaultParams
	NewRawDataReport = types.NewRawDataReport

	KeyMaxDataSourceExecutableSize  = types.KeyMaxDataSourceExecutableSize
	KeyMaxOracleScriptCodeSize      = types.KeyMaxOracleScriptCodeSize
	KeyMaxCalldataSize              = types.KeyMaxCalldataSize
	KeyMaxDataSourceCountPerRequest = types.KeyMaxDataSourceCountPerRequest
	KeyMaxRawDataReportSize         = types.KeyMaxRawDataReportSize
	KeyMaxResultSize                = types.KeyMaxResultSize

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
	Keeper                = keeper.Keeper
	MsgRequestData        = types.MsgRequestData
	MsgReportData         = types.MsgReportData
	MsgCreateDataSource   = types.MsgCreateDataSource
	MsgEditDataSource     = types.MsgEditDataSource
	MsgCreateOracleScript = types.MsgCreateOracleScript
	MsgEditOracleScript   = types.MsgEditOracleScript

	RawDataReport         = types.RawDataReport
	RequestQuerierInfo    = types.RequestQuerierInfo
	DataSourceQuerierInfo = types.DataSourceQuerierInfo

	RequestID      = types.RequestID
	OracleScriptID = types.OracleScriptID
	ExternalID     = types.ExternalID
	DataSourceID   = types.DataSourceID

	DataSource   = types.DataSource
	OracleScript = types.OracleScript
)
