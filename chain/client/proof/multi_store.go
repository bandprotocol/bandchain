package proof

import (
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

type MultiStoreProof struct {
	AccToMintStoresMerkleHash          tmbytes.HexBytes `json:"accToMintStoresMerkleHash"`
	OracleIAVLStateHash                tmbytes.HexBytes `json:"oracleIAVLStateHash"`
	ParamsStoresMerkleHash             tmbytes.HexBytes `json:"paramsStoresMerkleHash"`
	SlashingAndStakingStoresMerkleHash tmbytes.HexBytes `json:"slashingAndStakingStoresMerkleHash"`
	SupplyAndUpgradeStoresMerkleHash   tmbytes.HexBytes `json:"supplyAndUpgradeStoresMerkleHash"`
}

type MultiStoreProofEthereum struct {
	AccToMintStoresMerkleHash          common.Hash
	OracleStoresMerkleHash             common.Hash
	ParamsStoresMerkleHash             common.Hash
	SlashingAndStakingStoresMerkleHash common.Hash
	SupplyAndUpgradeStoresMerkleHash   common.Hash
}

func (m *MultiStoreProof) encodeToEthFormat() MultiStoreProofEthereum {
	return MultiStoreProofEthereum{
		AccToMintStoresMerkleHash:          common.BytesToHash(m.AccToMintStoresMerkleHash),
		OracleStoresMerkleHash:             common.BytesToHash(m.OracleIAVLStateHash),
		ParamsStoresMerkleHash:             common.BytesToHash(m.ParamsStoresMerkleHash),
		SlashingAndStakingStoresMerkleHash: common.BytesToHash(m.SlashingAndStakingStoresMerkleHash),
		SupplyAndUpgradeStoresMerkleHash:   common.BytesToHash(m.SupplyAndUpgradeStoresMerkleHash),
	}
}

func GetMultiStoreProof(proof rootmulti.MultiStoreProofOp) MultiStoreProof {
	m := make(map[string][]byte, len(proof.Proof.StoreInfos))
	for _, info := range proof.Proof.StoreInfos {
		m[info.Name] = info.Core.CommitID.Hash
	}
	return MultiStoreProof{
		AccToMintStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(auth.StoreKey, m[auth.StoreKey]),
			encodeStoreMerkleHash(bank.StoreKey, m[bank.StoreKey]),
			encodeStoreMerkleHash(distribution.StoreKey, m[distribution.StoreKey]),
			encodeStoreMerkleHash(evidence.StoreKey, m[evidence.StoreKey]),
			encodeStoreMerkleHash(gov.StoreKey, m[gov.StoreKey]),
			encodeStoreMerkleHash(ibc.StoreKey, m[ibc.StoreKey]),
			encodeStoreMerkleHash("main", m["main"]),
			encodeStoreMerkleHash(mint.StoreKey, m[mint.StoreKey]),
		}),
		OracleIAVLStateHash: m[oracle.StoreKey],
		ParamsStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(params.StoreKey, m[params.StoreKey]),
		}),
		SlashingAndStakingStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(slashing.StoreKey, m[slashing.StoreKey]),
			encodeStoreMerkleHash(staking.StoreKey, m[staking.StoreKey]),
		}),
		SupplyAndUpgradeStoresMerkleHash: merkle.SimpleHashFromByteSlices([][]byte{
			encodeStoreMerkleHash(supply.StoreKey, m[supply.StoreKey]),
			encodeStoreMerkleHash(upgrade.StoreKey, m[upgrade.StoreKey]),
		}),
	}
}
