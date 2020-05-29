package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

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

// CalculateResultHash returs hash of append hash of request and response packet.
func CalculateResultHash(req OracleRequestPacketData, res OracleResponsePacketData) []byte {
	reqPacketHash := tmhash.Sum(obi.MustEncode(req))
	resPacketHash := tmhash.Sum(obi.MustEncode(res))
	return tmhash.Sum(append(reqPacketHash, resPacketHash...))
}
