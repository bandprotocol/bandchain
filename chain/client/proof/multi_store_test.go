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
		Data: base64ToBytes("CoMFCjEKBnBhcmFtcxInCiUIj04SIMbe1bnhlZ/QXKhwt2bBDkuECSyXjQJR9i163emjYerLChEKCGV2aWRlbmNlEgUKAwiPTgovCgRtaW50EicKJQiPThIgHRIn/1HjJPruCU+0YCbGTMQdN3H8zPYSnG3BcM50HVQKLgoDaWJjEicKJQiPThIg2FdO5TPqLPUPb0eDKyz75v8F8tKYibGUk1v9291o/0wKMwoIc2xhc2hpbmcSJwolCI9OEiBlMr4w0+p/W2a2IsZN2j2r+eM3dO6T1ljzYyJHFF+ywwovCgRiYW5rEicKJQiPThIgIWl08l6pVCVCwTUOupBADhftNG9IR22d+7zKSMhWZEoKLgoDYWNjEicKJQiPThIgUtwTVqGGzhyrdXk50KSnftCBxHYMsAckSCFS3yQ618AKMQoGb3JhY2xlEicKJQiPThIgf+lq7ZMJJ8giRXz4iDVQXsNTOgtvtZD8NNngF1oYSywKNwoMZGlzdHJpYnV0aW9uEicKJQiPThIgeOvxhU5e2CiVjzs+y1MOQtZDJ3iK7JAwwKar4XXzRywKEAoHdXBncmFkZRIFCgMIj04KLwoEbWFpbhInCiUIj04SICOZd9vZgMqqD9fkR5Xfzpik2gT1mvTUmon2eg92gsm7Ci4KA2dvdhInCiUIj04SIGDNE3wZYuysYWOJ1oA0gz8pIVCcLShapfUVOZfOlop0CjIKB3N0YWtpbmcSJwolCI9OEiBFTY5j3HlY+MCENXlCrvRVhnB0kmTuHldZStDNU8yCsQoxCgZzdXBwbHkSJwolCI9OEiDql6ROudqdJXyzgcLKjpQ57Ov6JleivMlj3C2rH5wPqA=="),
	})
	require.Nil(t, err)
	mp := GetMultiStoreProof(multiStoreProofOp.(rootmulti.MultiStoreProofOp))
	expectAppHash := hexToBytes("9592EB9B13206F557F123FB98E9B4BC9B4963F9F8A2FA46A67BB421944FB2B08")
	oraclePrefix := hexToBytes("066f7261636c6520")
	appHash := branchHash(
		mp.AccToMintStoresMerkleHash,
		branchHash(
			branchHash(
				branchHash(
					leafHash(append(oraclePrefix, tmhash.Sum(tmhash.Sum(mp.OracleIAVLStateHash))...)),
					mp.ParamsStoresMerkleHash,
				),
				mp.SlashingAndStakingStoresMerkleHash,
			),
			mp.SupplyAndUpgradeStoresMerkleHash,
		),
	)
	require.Equal(t, expectAppHash, appHash)
}
