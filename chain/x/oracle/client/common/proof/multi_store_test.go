package proof

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func TestGetMultiStoreProof(t *testing.T) {
	multiStoreProofOp, err := rootmulti.MultiStoreProofOpDecoder(merkle.ProofOp{
		Type: "multistore",
		Key:  base64ToBytes("b3JhY2xl"), // oracle
		Data: base64ToBytes("vQQKugQKEwoIZXZpZGVuY2USBwoFCM21uAEKMAoDYWNjEikKJwjNtbgBEiCpnMcSKJ6M1amkAImRZTArIRHOeu/50Jo5rO3rvJfCqAozCgZvcmFjbGUSKQonCM21uAESIE+QC4tCXPhasqHtKQfUgwv2dHA8ZJVIR9qu6bgaCboxCjQKB3N0YWtpbmcSKQonCM21uAESIETcYTgf66RphwCHSdu3yuAVbSV9glVZXdFUnxmNnH5FCjAKA2dvdhIpCicIzbW4ARIgYM0TfBli7KxhY4nWgDSDPykhUJwtKFql9RU5l86WinQKMwoGcGFyYW1zEikKJwjNtbgBEiC/Owz4hHzwDasNUjeI+xGUW6ma+mq3LSz8KBe9+fZmZgoSCgd1cGdyYWRlEgcKBQjNtbgBCjEKBG1pbnQSKQonCM21uAESIP4eMZYd/UsxicOW2AUQ5QYXGhOl3uk4V6nnux9mciljCjEKBG1haW4SKQonCM21uAESILjQXXrLNbKLc8dmAUD3qSPSHHXTU8102Vf0C6PEQhfICjMKBnN1cHBseRIpCicIzbW4ARIgjgXbnPGDvS0BwswS3v1DkohXxqQ7kEZ36pSVKr507CEKNQoIc2xhc2hpbmcSKQonCM21uAESIOieO4Oh/Z+x7IC4cN3ORjl3K1s3olg5HZ/OvRuWVi6oCjkKDGRpc3RyaWJ1dGlvbhIpCicIzbW4ARIgbU8UcR/jBa7x3svYjLE/jB2hrPfYS1V1mbI0z2nMZFg="),
	})
	require.Nil(t, err)
	mp := GetMultiStoreProof(multiStoreProofOp.(rootmulti.MultiStoreProofOp))
	// Must get app hash from next block of multi proof.
	expectAppHash := hexToBytes("91C6C90AD6765C3080CEF2AEB25B1DBDD8ABE6EB409F400C3D6F8DC2767980F6")
	oraclePrefix := hexToBytes("066f7261636c6520")
	appHash := branchHash(
		branchHash(
			mp.AccToGovStoresMerkleHash,
			branchHash(
				mp.MainAndMintStoresMerkleHash,
				branchHash(
					leafHash(append(oraclePrefix, tmhash.Sum(tmhash.Sum(mp.OracleIAVLStateHash))...)),
					mp.ParamsStoresMerkleHash,
				),
			),
		),
		mp.SlashingToUpgradeStoresMerkleHash,
	)
	require.Equal(t, expectAppHash, appHash)
}
