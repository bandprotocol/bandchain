package proof

import (
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

type MultiStoreProof struct {
	AccToMemCapStoresMerkleHash       tmbytes.HexBytes `json:"accToMemCapStoresMerkleHash"`
	MintStoresMerkleHash              tmbytes.HexBytes `json:"mintStoresMerkleHash"`
	OracleIAVLStateHash               tmbytes.HexBytes `json:"oracleIAVLStateHash"`
	ParamsAndSlashingStoresMerkleHash tmbytes.HexBytes `json:"paramsAndSlashingStoresMerkleHash"`
	StakingAndUpgradeStoresMerkleHash tmbytes.HexBytes `json:"stakingAndUpgradeStoresMerkleHash"`
}

type MultiStoreProofEthereum struct {
	AccToMemCapStoresMerkleHash       common.Hash
	MintStoresMerkleHash              common.Hash
	OracleIAVLStateHash               common.Hash
	ParamsAndSlashingStoresMerkleHash common.Hash
	StakingAndUpgradeStoresMerkleHash common.Hash
}

func (m *MultiStoreProof) encodeToEthFormat() MultiStoreProofEthereum {
	return MultiStoreProofEthereum{
		AccToMemCapStoresMerkleHash:       common.BytesToHash(m.AccToMemCapStoresMerkleHash),
		MintStoresMerkleHash:              common.BytesToHash(m.MintStoresMerkleHash),
		OracleIAVLStateHash:               common.BytesToHash(m.OracleIAVLStateHash),
		ParamsAndSlashingStoresMerkleHash: common.BytesToHash(m.ParamsAndSlashingStoresMerkleHash),
		StakingAndUpgradeStoresMerkleHash: common.BytesToHash(m.StakingAndUpgradeStoresMerkleHash),
	}
}

func GetMultiStoreProof(proof rootmulti.MultiStoreProofOp) MultiStoreProof {
	m := make(map[string][]byte, len(proof.Proof.StoreInfos))
	for _, info := range proof.Proof.StoreInfos {
		m[info.Name] = info.Core.CommitID.Hash
	}
	return MultiStoreProof{
		// AccToMemCapStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
		// 	encodeStoreMerkleHash(auth.StoreKey, m[auth.StoreKey]),
		// 	encodeStoreMerkleHash(bank.StoreKey, m[bank.StoreKey]),
		// 	encodeStoreMerkleHash(capability.StoreKey, m[capability.StoreKey]),
		// 	encodeStoreMerkleHash(distribution.StoreKey, m[distribution.StoreKey]),
		// 	encodeStoreMerkleHash(evidence.StoreKey, m[evidence.StoreKey]),
		// 	encodeStoreMerkleHash(gov.StoreKey, m[gov.StoreKey]),
		// 	encodeStoreMerkleHash(ibc.StoreKey, m[ibc.StoreKey]),
		// 	encodeStoreMerkleHash(capability.MemStoreKey, m[capability.MemStoreKey]),
		// }),
		// MintStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
		// 	encodeStoreMerkleHash(mint.StoreKey, m[mint.StoreKey]),
		// }),
		// OracleIAVLStateHash: m[oracle.StoreKey],
		// ParamsAndSlashingStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
		// 	encodeStoreMerkleHash(params.StoreKey, m[params.StoreKey]),
		// 	encodeStoreMerkleHash(slashing.StoreKey, m[slashing.StoreKey]),
		// }),
		// StakingAndUpgradeStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
		// 	encodeStoreMerkleHash(staking.StoreKey, m[staking.StoreKey]),
		// 	encodeStoreMerkleHash(upgrade.StoreKey, m[upgrade.StoreKey]),
		// }),
	}
}
