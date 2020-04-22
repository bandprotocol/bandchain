package proof

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func getParentHash(path IAVLMerklePath, subtreeHash []byte) []byte {
	left, right := subtreeHash, path.SiblingHash
	if path.IsDataOnRight {
		left, right = right, left
	}
	raw := []byte{2 * path.SubtreeHeight}
	raw = append(raw, encodeVarint(int64(path.SubtreeSize))...)
	raw = append(raw, encodeVarint(int64(path.SubtreeVersion))...)
	raw = append(raw, 32)
	raw = append(raw, left...)
	raw = append(raw, 32)
	raw = append(raw, right...)
	return tmhash.Sum(raw)
}
func TestGetIAVLMerklePaths(t *testing.T) {
	var iavlProof iavl.ValueOp
	iavlOpData := base64ToBytes("twIKtAIKKQgOEDwYj04qIAqVHc4yEL3Lr0lZ0cQPJ+nMt/TaT+H5Cd/pZOcdWmWZCikIDBAdGI9OKiAOo7imPYU0LXqSghAX34+evA8tHUnu6SI3mWqUAS5rcAopCAgQDBiPTiIgVtRzIyliF58RDwWGYJIy88CWFYN1L9BiuTV3IJPnCnUKKQgGEAgY8ikqIIQMVU13Mxhj2pnxGywAfwoRaRX7qIeWcSAqI/ot7OwsCikIBBAEGIoiKiCe4B7tdcsF0tG2W+f+nMvQdzjOsQkIvqNxbNHXfpL69gopCAIQAhiNDSog/tzkcYxC546anyOplbnELEARIvc8NNk990HEd3X4dZ0aMAoJAQAAAAAAAAABEiBF01aMSserjw6ZditfkKaR1q0SHsDI0mptE2Nzxud0RBioAQ==")
	amino.NewCodec().MustUnmarshalBinaryLengthPrefixed(iavlOpData, &iavlProof)
	paths := GetIAVLMerklePaths(&iavlProof)
	expectOracleMerkleHash := hexToBytes("7FE96AED930927C822457CF88835505EC3533A0B6FB590FC34D9E0175A184B2C")
	key := base64ToBytes("AQAAAAAAAAAB")
	value := base64ToBytes("CAMSCGQAAAAAAAAAGhQNXkvEVt3wx+cADHeOiIeMO08NEhoUkVlvUlJZHqz+7gAlm/YBTD/c7OEaFMIVrZXki61HhK2Zw5304t96pQfmGhTfX4M2rg4BaNheSM5tPKMG0AQwYSAEKhTCFa2V5IutR4StmcOd9OLfeqUH5ioUDV5LxFbd8MfnAAx3joiHjDtPDRIqFJFZb1JSWR6s/u4AJZv2AUw/3OzhKhTfX4M2rg4BaNheSM5tPKMG0AQwYTCkATi6o9T0BUC4AUgC")
	leafNode := []byte{0}                                                           // Height
	leafNode = append(leafNode, 2)                                                  // subtree size = 1 * 2
	leafNode = append(leafNode, encodeVarint(iavlProof.Proof.Leaves[0].Version)...) // version
	leafNode = append(leafNode, uint8(len(key)))                                    // key length
	leafNode = append(leafNode, key...)                                             // key
	leafNode = append(leafNode, 32)                                                 // value hash size
	leafNode = append(leafNode, tmhash.Sum(value)...)                               //value
	currentHash := tmhash.Sum(leafNode)
	for _, path := range paths {
		currentHash = getParentHash(path, currentHash)
	}
	require.Equal(t, expectOracleMerkleHash, currentHash)
}
