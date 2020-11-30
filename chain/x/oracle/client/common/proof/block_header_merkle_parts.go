package proof

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

// BlockHeaderMerkleParts stores a group of hashes using for computing Tendermint's block
// header hash from app hash, and height.
//
// In Tendermint, a block header hash is the Merkle hash of a binary tree with 14 leaf nodes.
// Each node encodes a data piece of the blockchain. The notable data leaves are: [A] app_hash,
// [2] height. All data pieces are combined into one 32-byte hash to be signed
// by block validators. The structure of the Merkle tree is shown below.
//
//                                   [BlockHeader]
//                                /                \
//                   [3A]                                    [3B]
//                 /      \                                /      \
//         [2A]                [2B]                [2C]                [2D]
//        /    \              /    \              /    \              /    \
//    [1A]      [1B]      [1C]      [1D]      [1E]      [1F]        [C]    [D]
//    /  \      /  \      /  \      /  \      /  \      /  \
//  [0]  [1]  [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [A]  [B]
//
//  [0] - version               [1] - chain_id            [2] - height        [3] - time
//  [4] - last_block_id         [5] - last_commit_hash    [6] - data_hash     [7] - validators_hash
//  [8] - next_validators_hash  [9] - consensus_hash      [A] - app_hash      [B] - last_results_hash
//  [C] - evidence_hash         [D] - proposer_address
//
// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
// root hash, since we only want to validate the correctness of [2], [3], and [A]. In fact, only
// [1A], [2B], [1E], [B], and [2D] are needed in order to compute [BlockHeader].
type BlockHeaderMerkleParts struct {
	VersionAndChainIdHash             tmbytes.HexBytes `json:"versionAndChainIdHash"`
	Height                            uint64           `json:"height"`
	TimeSecond                        uint64           `json:"timeSecond"`
	TimeNanoSecond                    uint32           `json:"timeNanoSecond"`
	LastBlockIDAndOther               tmbytes.HexBytes `json:"lastBlockIDAndOther"`
	NextValidatorHashAndConsensusHash tmbytes.HexBytes `json:"nextValidatorHashAndConsensusHash"`
	LastResultsHash                   tmbytes.HexBytes `json:"lastResultsHash"`
	EvidenceAndProposerHash           tmbytes.HexBytes `json:"evidenceAndProposerHash"`
}

// BlockHeaderMerklePartsEthereum is an Ethereum version of BlockHeaderMerkleParts for solidity ABI-encoding.
type BlockHeaderMerklePartsEthereum struct {
	VersionAndChainIdHash             common.Hash
	Height                            uint64
	TimeSecond                        uint64
	TimeNanoSecond                    uint32
	LastBlockIDAndOther               common.Hash
	NextValidatorHashAndConsensusHash common.Hash
	LastResultsHash                   common.Hash
	EvidenceAndProposerHash           common.Hash
}

func (bp *BlockHeaderMerkleParts) encodeToEthFormat() BlockHeaderMerklePartsEthereum {
	return BlockHeaderMerklePartsEthereum{
		VersionAndChainIdHash:             common.BytesToHash(bp.VersionAndChainIdHash),
		Height:                            bp.Height,
		TimeSecond:                        bp.TimeSecond,
		TimeNanoSecond:                    bp.TimeNanoSecond,
		LastBlockIDAndOther:               common.BytesToHash(bp.LastBlockIDAndOther),
		NextValidatorHashAndConsensusHash: common.BytesToHash(bp.NextValidatorHashAndConsensusHash),
		LastResultsHash:                   common.BytesToHash(bp.LastResultsHash),
		EvidenceAndProposerHash:           common.BytesToHash(bp.EvidenceAndProposerHash),
	}
}

// GetBlockHeaderMerkleParts converts Tendermint block header struct into BlockHeaderMerkleParts for gas-optimized proof verification.
func GetBlockHeaderMerkleParts(codec *codec.Codec, block *types.Header) BlockHeaderMerkleParts {
	return BlockHeaderMerkleParts{
		VersionAndChainIdHash: merkle.SimpleHashFromByteSlices([][]byte{
			cdcEncode(codec, block.Version),
			cdcEncode(codec, block.ChainID),
		}),
		Height:         uint64(block.Height),
		TimeSecond:     uint64(block.Time.Unix()),
		TimeNanoSecond: uint32(block.Time.Nanosecond()),
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
