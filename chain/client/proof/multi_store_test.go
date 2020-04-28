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
		Data: base64ToBytes("CvYEChAKB3VwZ3JhZGUSBQoDCJcFChMKCmNhcGFiaWxpdHkSBQoDCJcFCi8KBGJhbmsSJwolCJcFEiAH4idukNC6oNbM2oVMWC2p+Sy6yPOEGWJJFnqogL0omQovCgRtaW50EicKJQiXBRIgLvOu+wCxtoTxBkQcVejiXt9HE6a+JY7TlMn4Vccav3EKMQoGc3VwcGx5EicKJQiXBRIg7lXJWNUTysK0qF4YGY1LOf8ZMBs5eI9ViehyzGTqbdQKNwoMZGlzdHJpYnV0aW9uEicKJQiXBRIgVAY5wSAUnHco8s2vBueJj3Es5p49YzbHd16SJoKntIoKLwoEbWFpbhInCiUIlwUSICOZd9vZgMqqD9fkR5Xfzpik2gT1mvTUmon2eg92gsm7CjEKBnBhcmFtcxInCiUIlwUSIItCIwSHU2bIQ8Ggtr1rJ2ZBvDhlzprR1ObYiQGIRgDzCjEKBm9yYWNsZRInCiUIlwUSINDuKe2xqA+AttwsBYsH6FhG4qHU7En84d0M8blGzPRWCi4KA2FjYxInCiUIlwUSIMKNnINAAgLun6gVLILvB/6oUDxkenAiZLf1in6tD4qMChEKCGV2aWRlbmNlEgUKAwiXBQouCgNnb3YSJwolCJcFEiBgzRN8GWLsrGFjidaANIM/KSFQnC0oWqX1FTmXzpaKdAoMCgNpYmMSBQoDCJcFCjIKB3N0YWtpbmcSJwolCJcFEiACWdb3d8lIViqQ8ySrrEN9/3o2f3Z5p/L5396OhiGE/gozCghzbGFzaGluZxInCiUIlwUSIFa61WCgLNn11FrW2g4xYzGbYmWfuP32bQywySwABY7m"),
	})
	require.Nil(t, err)
	mp := GetMultiStoreProof(multiStoreProofOp.(rootmulti.MultiStoreProofOp))
	expectAppHash := hexToBytes("2F3BEAC1586C205052B74E1CF3D284CD022F739200B74CB51B910D0F3D0BF13D")
	oraclePrefix := hexToBytes("066f7261636c6520")
	appHash := branchHash(
		mp.AccToMainStoresMerkleHash,
		branchHash(
			branchHash(
				branchHash(
					mp.MintStoresMerkleHash,
					leafHash(append(oraclePrefix, tmhash.Sum(tmhash.Sum(mp.OracleIAVLStateHash))...)),
				),
				mp.ParamsAndSlashingStoresMerkleHash,
			),
			mp.StakingToUpgradeStoresMerkleHash,
		),
	)
	require.Equal(t, expectAppHash, appHash)
}
