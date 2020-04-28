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
	iavlOpData := base64ToBytes("4QEK3gEKKQgIEA0YlwUiINVGvQopIqUP0YSDHayhCOdCoUw9M5ySp0Z9MbzXsNrSCikIBhAIGNcEIiBwnhxzURsk792bjTzXF6UhC6IOJBGoUp6LZCxU+wAtxAopCAQQBBjXBCIgW3C/rdFuyVQJBy+zaGvyv+xIET+WysqIOXiDEQxnLxMKKQgCEAIY1wQiIBRZu8fbf/zOPezuo98GKWj049KyBrk9WfqTbjNLnsQ0GjAKCf8AAAAAAAAAARIgyy49YoXIcCE9OvooX05O2waYLsKbzjnA9Y9KhRBC0ZwY1wQ=")
	amino.NewCodec().MustUnmarshalBinaryLengthPrefixed(iavlOpData, &iavlProof)
	paths := GetIAVLMerklePaths(&iavlProof)
	expectOracleMerkleHash := hexToBytes("D0EE29EDB1A80F80B6DC2C058B07E85846E2A1D4EC49FCE1DD0CF1B946CCF456")
	key := base64ToBytes("/wAAAAAAAAAB")                                            //key to result of request#1
	value := base64ToBytes("pStdspbXz/kj9zcT3FLdqRewtPW9azeI/0hs0H7x5lE=")          // result hash
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
