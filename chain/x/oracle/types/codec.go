package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
)

// ModuleCdc is the codec for the module.
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	channel.RegisterCodec(ModuleCdc)
	commitmenttypes.RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequestData{}, "oracle/Request", nil)
	cdc.RegisterConcrete(MsgReportData{}, "oracle/Report", nil)
	cdc.RegisterConcrete(MsgCreateDataSource{}, "oracle/CreateDataSource", nil)
	cdc.RegisterConcrete(MsgEditDataSource{}, "oracle/EditDataSource", nil)
	cdc.RegisterConcrete(MsgCreateOracleScript{}, "oracle/CreateOracleScript", nil)
	cdc.RegisterConcrete(MsgEditOracleScript{}, "oracle/EditOracleScript", nil)
	cdc.RegisterConcrete(MsgAddOracleAddress{}, "oracle/AddOracleAddress", nil)
	cdc.RegisterConcrete(MsgRemoveOracleAddress{}, "oracle/RemoveOracleAddress", nil)
	cdc.RegisterConcrete(OracleRequestPacketData{}, "oracle/OracleRequestPacketData", nil)
	cdc.RegisterConcrete(OracleResponsePacketData{}, "oracle/OracleResponsePacketData", nil)
}
