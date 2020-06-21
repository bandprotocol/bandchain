package proof

// import (
// 	"testing"

// 	"github.com/cosmos/cosmos-sdk/store/rootmulti"
// 	"github.com/stretchr/testify/require"
// 	"github.com/tendermint/tendermint/crypto/merkle"
// 	"github.com/tendermint/tendermint/crypto/tmhash"
// )

// func TestGetMultiStoreProof(t *testing.T) {
// 	multiStoreProofOp, err := rootmulti.MultiStoreProofOpDecoder(merkle.ProofOp{
// 		Type: "multistore",
// 		Key:  base64ToBytes("b3JhY2xl"), // oracle
// 		Data: base64ToBytes("CpcECi4KBGJhbmsSJgokCE8SIFFp5VJl0pN7IB1zm+c7wc160Ev58w29dDHjWZahsS4NChAKDm1lbV9jYXBhYmlsaXR5CjYKDGRpc3RyaWJ1dGlvbhImCiQITxIgLgd/rDfToCd79w5sMBEd5Cz0Wx9T8tz5P0T7imhrxacKMQoHc3Rha2luZxImCiQITxIgWyzwEPdl+6QjNsYAPuWt5s3ygu2rT6Wf755pFPA3ZLMKCwoDaWJjEgQKAghPChIKCmNhcGFiaWxpdHkSBAoCCE8KLgoEbWludBImCiQITxIgn0Ay9in3n018bViYNvakJPGs5NBNApBDJlanzTkoQV8KMAoGcGFyYW1zEiYKJAhPEiD2mMDOrO1Ai0IdR5wU1o/oe9Z6lAATKJaOgfGPNre8UQoQCghldmlkZW5jZRIECgIITwoyCghzbGFzaGluZxImCiQITxIgL+ePlJmSe0nErlzaVL1/XwR5JA0T2kLjsqb/yDCngXUKLQoDZ292EiYKJAhPEiBgzRN8GWLsrGFjidaANIM/KSFQnC0oWqX1FTmXzpaKdAotCgNhY2MSJgokCE8SIDGAt5a/EEkpGtM+eLElNyck5q0BIu0WYE21u2j6bVReCg8KB3VwZ3JhZGUSBAoCCE8KMAoGb3JhY2xlEiYKJAhPEiCywp9PjsehPrpPthUq4Ef1OiNCgJcaQIGwONL8PDs+2A=="),
// 	})
// 	require.Nil(t, err)
// 	mp := GetMultiStoreProof(multiStoreProofOp.(rootmulti.MultiStoreProofOp))
// 	expectAppHash := hexToBytes("396EC5A6D5D66C0EA13BB6D5BDE8A01D80F78F3BB73D62C2F71BC14DD8B18095")
// 	oraclePrefix := hexToBytes("066f7261636c6520")
// 	appHash := branchHash(
// 		mp.AccToMemCapStoresMerkleHash,
// 		branchHash(
// 			branchHash(
// 				branchHash(
// 					mp.MintStoresMerkleHash,
// 					leafHash(append(oraclePrefix, tmhash.Sum(tmhash.Sum(mp.OracleIAVLStateHash))...)),
// 				),
// 				mp.ParamsAndSlashingStoresMerkleHash,
// 			),
// 			mp.StakingAndUpgradeStoresMerkleHash,
// 		),
// 	)
// 	require.Equal(t, expectAppHash, appHash)
// }
