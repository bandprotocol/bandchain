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
	cdc.RegisterConcrete(MsgRequestData{}, "zoracle/Request", nil)
	cdc.RegisterConcrete(MsgReportData{}, "zoracle/Report", nil)
	cdc.RegisterConcrete(MsgCreateDataSource{}, "zoracle/CreateDataSource", nil)
	cdc.RegisterConcrete(MsgEditDataSource{}, "zoracle/EditDataSource", nil)
	cdc.RegisterConcrete(MsgCreateOracleScript{}, "zoracle/CreateOracleScript", nil)
	cdc.RegisterConcrete(MsgEditOracleScript{}, "zoracle/EditOracleScript", nil)
	cdc.RegisterConcrete(MsgAddOracleAddress{}, "zoracle/AddOracleAddress", nil)
	cdc.RegisterConcrete(MsgRemoveOracleAddress{}, "zoracle/RemoveOracleAddress", nil)
	cdc.RegisterConcrete(OracleRequestPacketData{}, "zoracle/OracleRequestPacketData", nil)
}
