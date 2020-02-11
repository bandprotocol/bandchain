package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRequest{}, "zoracle/Request", nil)
	cdc.RegisterConcrete(MsgReport{}, "zoracle/Report", nil)
	cdc.RegisterConcrete(MsgStoreCode{}, "zoracle/Store", nil)
	cdc.RegisterConcrete(MsgDeleteCode{}, "zoracle/Delete", nil)
	cdc.RegisterConcrete(MsgRequestData{}, "zoracle/RequestData", nil)
	cdc.RegisterConcrete(MsgReportData{}, "zoracle/ReportData", nil)
}
