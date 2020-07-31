package proof

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/iavl"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// IAVLMerklePath represents a Merkle step to a leaf data node in an iAVL tree.
type IAVLMerklePath struct {
	IsDataOnRight  bool             `json:"isDataOnRight"`
	SubtreeHeight  uint8            `json:"subtreeHeight"`
	SubtreeSize    uint64           `json:"subtreeSize"`
	SubtreeVersion uint64           `json:"subtreeVersion"`
	SiblingHash    tmbytes.HexBytes `json:"siblingHash"`
}

// IAVLMerklePathEthereum is an Ethereum version of IAVLMerklePath for solidity ABI-encoding.
type IAVLMerklePathEthereum struct {
	IsDataOnRight  bool
	SubtreeHeight  uint8
	SubtreeSize    *big.Int
	SubtreeVersion *big.Int
	SiblingHash    common.Hash
}

func (merklePath *IAVLMerklePath) encodeToEthFormat() IAVLMerklePathEthereum {
	return IAVLMerklePathEthereum{
		merklePath.IsDataOnRight,
		merklePath.SubtreeHeight,
		big.NewInt(int64(merklePath.SubtreeSize)),
		big.NewInt(int64(merklePath.SubtreeVersion)),
		common.BytesToHash(merklePath.SiblingHash),
	}
}

// GetIAVLMerklePaths returns the list of IAVLMerklePath elements from the given iAVL proof.
func GetIAVLMerklePaths(proof *iavl.ValueOp) []IAVLMerklePath {
	paths := make([]IAVLMerklePath, 0)
	for i := len(proof.Proof.LeftPath) - 1; i >= 0; i-- {
		path := proof.Proof.LeftPath[i]
		imp := IAVLMerklePath{}
		imp.SubtreeHeight = uint8(path.Height)
		imp.SubtreeSize = uint64(path.Size)
		imp.SubtreeVersion = uint64(path.Version)
		if len(path.Right) == 0 {
			imp.SiblingHash = path.Left
			imp.IsDataOnRight = true
		} else {
			imp.SiblingHash = path.Right
			imp.IsDataOnRight = false
		}
		paths = append(paths, imp)
	}
	return paths
}
