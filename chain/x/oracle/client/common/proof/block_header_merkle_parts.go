package proof

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

type BlockHeaderMerkleParts struct {
	VersionAndChainIdHash             tmbytes.HexBytes `json:"versionAndChainIdHash"`
	TimeHash                          tmbytes.HexBytes `json:"timeHash"`
	LastBlockIDAndOther               tmbytes.HexBytes `json:"lastBlockIDAndOther"`
	NextValidatorHashAndConsensusHash tmbytes.HexBytes `json:"nextValidatorHashAndConsensusHash"`
	LastResultsHash                   tmbytes.HexBytes `json:"lastResultsHash"`
	EvidenceAndProposerHash           tmbytes.HexBytes `json:"evidenceAndProposerHash"`
}

type BlockHeaderMerklePartsEthereum struct {
	VersionAndChainIdHash             common.Hash
	TimeHash                          common.Hash
	LastBlockIDAndOther               common.Hash
	NextValidatorHashAndConsensusHash common.Hash
	LastResultsHash                   common.Hash
	EvidenceAndProposerHash           common.Hash
}

func (bp *BlockHeaderMerkleParts) encodeToEthFormat() BlockHeaderMerklePartsEthereum {
	return BlockHeaderMerklePartsEthereum{
		VersionAndChainIdHash:             common.BytesToHash(bp.VersionAndChainIdHash),
		TimeHash:                          common.BytesToHash(bp.TimeHash),
		LastBlockIDAndOther:               common.BytesToHash(bp.LastBlockIDAndOther),
		NextValidatorHashAndConsensusHash: common.BytesToHash(bp.NextValidatorHashAndConsensusHash),
		LastResultsHash:                   common.BytesToHash(bp.LastResultsHash),
		EvidenceAndProposerHash:           common.BytesToHash(bp.EvidenceAndProposerHash),
	}
}

func GetBlockHeaderMerkleParts(codec *codec.Codec, block *types.Header) BlockHeaderMerkleParts {
	return BlockHeaderMerkleParts{
		VersionAndChainIdHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.Version),
			cdcEncode(codec, block.ChainID),
		}),
		TimeHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.Time),
		}),
		LastBlockIDAndOther: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.LastBlockID),
			cdcEncode(codec, block.LastCommitHash),
			cdcEncode(codec, block.DataHash),
			cdcEncode(codec, block.ValidatorsHash),
		}),
		NextValidatorHashAndConsensusHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.NextValidatorsHash),
			cdcEncode(codec, block.ConsensusHash),
		}),
		LastResultsHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.LastResultsHash),
		}),
		EvidenceAndProposerHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.EvidenceHash),
			cdcEncode(codec, block.ProposerAddress),
		}),
	}
}
