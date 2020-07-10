package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module.
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers the module's concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequestData{}, "oracle/Request", nil)
	cdc.RegisterConcrete(MsgReportData{}, "oracle/Report", nil)
	cdc.RegisterConcrete(MsgCreateDataSource{}, "oracle/CreateDataSource", nil)
	cdc.RegisterConcrete(MsgEditDataSource{}, "oracle/EditDataSource", nil)
	cdc.RegisterConcrete(MsgCreateOracleScript{}, "oracle/CreateOracleScript", nil)
	cdc.RegisterConcrete(MsgEditOracleScript{}, "oracle/EditOracleScript", nil)
	cdc.RegisterConcrete(MsgActivate{}, "oracle/Activate", nil)
	cdc.RegisterConcrete(MsgAddReporter{}, "oracle/AddReporter", nil)
	cdc.RegisterConcrete(MsgRemoveReporter{}, "oracle/RemoveReporter", nil)
	cdc.RegisterConcrete(OracleRequestPacketData{}, "oracle/OracleRequestPacketData", nil)
	cdc.RegisterConcrete(OracleResponsePacketData{}, "oracle/OracleResponsePacketData", nil)
}
