package proof

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func TestVerifySignature(t *testing.T) {
	signatures := []TMSignature{
		{
			R:                hexToBytes("81AC28C67F636974BDC70D1A694BA050652FBAA9AA83A1F8B7B10F84C6BC9171"),
			S:                hexToBytes("5D022B62644E496504FC6AF1DD138544CD311F26A979DF4B4DAF92E60CA0F762"),
			V:                28,
			SignedDataSuffix: hexToBytes("12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F50510D581FBA001320962616E64636861696E"),
		},
		{
			R:                hexToBytes("F4ABDF0CB47604292B9B0D9636692A0D5379B646EA3246180004BFEAD2D7CA8A"),
			S:                hexToBytes("1AF744F61921AA03D5327333F654747C928CDD1D324A27FF79181AC8E1F6841E"),
			V:                27,
			SignedDataSuffix: hexToBytes("12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F50510FA86D8A701320962616E64636861696E"),
		},
		{
			R:                hexToBytes("4258784CC9659EEC320EA86AB7DD1C41C7BF8E9F22035B9E50FA8B527A6079BE"),
			S:                hexToBytes("35C1D785DA88F2D0D563E3AA64B15B96E7C53D025E85895D37F25D99AD11CA14"),
			V:                27,
			SignedDataSuffix: hexToBytes("12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F5051089FE8FA601320962616E64636861696E"),
		},
	}

	addresses := []common.Address{
		common.HexToAddress("652D89a66Eb4eA55366c45b1f9ACfc8e2179E1c5"),
		common.HexToAddress("aAA22E077492CbaD414098EBD98AA8dc1C7AE8D9"),
		common.HexToAddress("B956589b6fC5523eeD0d9eEcfF06262Ce84ff260"),
	}

	signedDataPrefix := hexToBytes("6E080211400200000000000022480A20")
	blockHash := hexToBytes("8D1897D04B6B4746021EAF4BF80F0FF9E5A4DDDD451A2FA814DA1C21380D69F8")

	for i, sig := range signatures {
		msg := append(signedDataPrefix, blockHash...)
		msg = append(msg, sig.SignedDataSuffix...)

		sigBytes := append(sig.R, sig.S...)
		sigBytes = append(sigBytes, sig.V-27)
		pub, err := crypto.SigToPub(tmhash.Sum(msg), sigBytes)
		require.Nil(t, err)
		require.Equal(t, addresses[i], crypto.PubkeyToAddress(*pub))
	}
}
