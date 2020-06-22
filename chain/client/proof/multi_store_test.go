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
		Data: base64ToBytes("pQQKogQKLgoDZ292EicKJQiWAhIgYM0TfBli7KxhY4nWgDSDPykhUJwtKFql9RU5l86WinQKMwoIc2xhc2hpbmcSJwolCJYCEiBfW/Loa4v0ZtxLkuWPQUbz6SVJuKg/E7X4EiG7Hkj/rwovCgRtYWluEicKJQiWAhIgI5l329mAyqoP1+RHld/OmKTaBPWa9NSaifZ6D3aCybsKNwoMZGlzdHJpYnV0aW9uEicKJQiWAhIgAnl+IKwaWnB6M5DPJxnTu+vh51ghtz16f0p4hpsklVEKMQoGcGFyYW1zEicKJQiWAhIgtxWMDPGF+SS8y9rUgmfiSnGykZ3birnksXSJLVsKgn8KLgoDYWNjEicKJQiWAhIgUYQhvmWNF7BiiczyeWHl2ErSxi7AwryrVx1NwhzagfYKMQoGc3VwcGx5EicKJQiWAhIgySh5gpbubLQHNcPKTpNN/CuKDkESvM+1XcfaGw+fydgKMQoGb3JhY2xlEicKJQiWAhIgFTrjBbDFiOtarT5oX7dAsnD9KH7OCLzIXWaPyxDjwfUKEAoHdXBncmFkZRIFCgMIlgIKMgoHc3Rha2luZxInCiUIlgISIPSPOKwdrwBrRV3dGze5ZM0VjHqf6K/7E1E1gjPQTQvcChEKCGV2aWRlbmNlEgUKAwiWAgovCgRtaW50EicKJQiWAhIgVoD/WGhQAMfwiRBMeyy0n/eg11aKq5Aof00pfjTRqCs="),
	})
	require.Nil(t, err)
	mp := GetMultiStoreProof(multiStoreProofOp.(rootmulti.MultiStoreProofOp))
	// Must get app hash from next block of multi proof.
	expectAppHash := hexToBytes("0DEF6341481C4370D561D546C268FB6AFED8520689B70ABC615C47CFB2A0EEE8")
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
