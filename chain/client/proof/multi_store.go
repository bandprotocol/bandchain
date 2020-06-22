package proof

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
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
)

type MultiStoreProof struct {
	AccToGovStoresMerkleHash          tmbytes.HexBytes `json:"accToGovStoresMerkleHash"`
	MainAndMintStoresMerkleHash       tmbytes.HexBytes `json:"mainAndMintStoresMerkleHash"`
	OracleIAVLStateHash               tmbytes.HexBytes `json:"oracleIAVLStateHash"`
	ParamsStoresMerkleHash            tmbytes.HexBytes `json:"paramsStoresMerkleHash"`
	SlashingToUpgradeStoresMerkleHash tmbytes.HexBytes `json:"stakingAndUpgradeStoresMerkleHash"`
}

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
		OracleIAVLStateHash:    m[oracle.StoreKey],
		ParamsStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{encodeStoreMerkleHash(params.StoreKey, m[params.StoreKey])}),
		SlashingToUpgradeStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(slashing.StoreKey, m[slashing.StoreKey]),
			encodeStoreMerkleHash(staking.StoreKey, m[staking.StoreKey]),
			encodeStoreMerkleHash(supply.StoreKey, m[supply.StoreKey]),
			encodeStoreMerkleHash(upgrade.StoreKey, m[upgrade.StoreKey]),
		}),
	}
}
