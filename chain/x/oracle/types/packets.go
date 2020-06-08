package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
)

// GetBytes returns the bytes representation of this oracle request packet data.
func (p OracleRequestPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

// GetBytes returns the bytes representation of this oracle response packet data.
func (p OracleResponsePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

// CalculateEncodedResult returns append obi encode of request and response packet.
func CalculateEncodedResult(req OracleRequestPacketData, res OracleResponsePacketData) []byte {
	encodedReq := obi.MustEncode(req)
	encodedRes := obi.MustEncode(res)
	return append(encodedReq, encodedRes...)
}
