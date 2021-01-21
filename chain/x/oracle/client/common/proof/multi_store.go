package proof

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// MultiStoreProof stores a compact of other Cosmos-SDK modules' storage hash in multistore to
// compute (in combination with oracle store hash) Tendermint's application state hash at a given block.
//                         ________________[AppHash]_______________
//                        /                                        \
//             _______[I9]______                          ________[I10]________
//            /                  \                       /                     \
//       __[I5]__             __[I6]__              __[I7]__               __[I8]__
//      /         \          /         \           /         \            /         \
//    [I1]       [I2]     [I3]        [I4]       [8]        [9]          [A]        [B]
//   /   \      /   \    /    \      /    \
// [0]   [1]  [2]   [3] [4]   [5]  [6]    [7]
// [0] - acc      [1] - distr   [2] - evidence  [3] - gov
// [4] - main     [5] - mint    [6] - oracle    [7] - params
// [8] - slashing [9] - staking [A] - supply    [D] - upgrade
// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
// root hash, since we only want to validate the correctness of [6] In fact, only
// [7], [I3], [I5], and [I10] are needed in order to compute [AppHash].
type MultiStoreProof struct {
	AccToGovStoresMerkleHash          tmbytes.HexBytes `json:"accToGovStoresMerkleHash"`
	MainAndMintStoresMerkleHash       tmbytes.HexBytes `json:"mainAndMintStoresMerkleHash"`
	OracleIAVLStateHash               tmbytes.HexBytes `json:"oracleIAVLStateHash"`
	ParamsStoresMerkleHash            tmbytes.HexBytes `json:"paramsStoresMerkleHash"`
	SlashingToUpgradeStoresMerkleHash tmbytes.HexBytes `json:"slashingToUpgradeStoresMerkleHash"`
}

// MultiStoreProofEthereum is an Ethereum version of MultiStoreProof for solidity ABI-encoding.
type MultiStoreProofEthereum struct {
	AccToGovStoresMerkleHash          common.Hash
	MainAndMintStoresMerkleHash       common.Hash
	OracleIAVLStateHash               common.Hash
	ParamsStoresMerkleHash            common.Hash
	SlashingToUpgradeStoresMerkleHash common.Hash
}

func (m *MultiStoreProof) encodeToEthFormat() MultiStoreProofEthereum {
	return MultiStoreProofEthereum{
		AccToGovStoresMerkleHash:          common.BytesToHash(m.AccToGovStoresMerkleHash),
		MainAndMintStoresMerkleHash:       common.BytesToHash(m.MainAndMintStoresMerkleHash),
		OracleIAVLStateHash:               common.BytesToHash(m.OracleIAVLStateHash),
		ParamsStoresMerkleHash:            common.BytesToHash(m.ParamsStoresMerkleHash),
		SlashingToUpgradeStoresMerkleHash: common.BytesToHash(m.SlashingToUpgradeStoresMerkleHash),
	}
}

// GetMultiStoreProof compacts Multi store proof from Tendermint to MultiStoreProof version.
func GetMultiStoreProof(proof rootmulti.MultiStoreProofOp) MultiStoreProof {
	m := make(map[string][]byte, len(proof.Proof.StoreInfos))
	for _, info := range proof.Proof.StoreInfos {
		m[info.Name] = info.Core.CommitID.Hash
	}
	return MultiStoreProof{
		AccToGovStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(auth.StoreKey, m[auth.StoreKey]),
			encodeStoreMerkleHash(distribution.StoreKey, m[distribution.StoreKey]),
			encodeStoreMerkleHash(evidence.StoreKey, m[evidence.StoreKey]),
			encodeStoreMerkleHash(gov.StoreKey, m[gov.StoreKey]),
		}),
		MainAndMintStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(bam.MainStoreKey, m[bam.MainStoreKey]),
			encodeStoreMerkleHash(mint.StoreKey, m[mint.StoreKey]),
		}),
		OracleIAVLStateHash:    m[types.StoreKey],
		ParamsStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{encodeStoreMerkleHash(params.StoreKey, m[params.StoreKey])}),
		SlashingToUpgradeStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(slashing.StoreKey, m[slashing.StoreKey]),
			encodeStoreMerkleHash(staking.StoreKey, m[staking.StoreKey]),
			encodeStoreMerkleHash(supply.StoreKey, m[supply.StoreKey]),
			encodeStoreMerkleHash(upgrade.StoreKey, m[upgrade.StoreKey]),
		}),
	}
}
