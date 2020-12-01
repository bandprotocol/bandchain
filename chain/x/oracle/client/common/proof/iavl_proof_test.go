package proof

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

/*
{
	code: 0,
	log: "",
	info: "",
	index: "0",
	key: "/wAAAAAAAAAB",
	value: "AAAAAAAAAAAAAAABAAAAFgAAAANMUkMAAAADVVNEAAAAAAAPQkAAAAAAAAAABwAAAAAAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAXyquiQAAAABfKq6PAAAAAQAAAAgAAAAAAAH1Ug==",
	proof: {
		ops: [
			{
				type: "iavl:v",
				key: "/wAAAAAAAAAB",
				data: "rAgKqQgKLggwEKnlwwMYzbW4ASIguYvJR88nWJD/ZP+u7z3k9i117OAXaX/ytarsSCym/m8KLgguEKnd9gIYzLW4ASIgV7NitOHbf1qu1Do0E4sCHUSmrlsrSjtuBDnIsZMdYw4KLggsEObxwAEYzLW4ASIgN+xxrBpxT4VVX9lOVR632Wt+ELRxSM268b9Qs7JFqKIKLQgqEOXhShjMtbgBKiCnvFfCkSTeeOd5jLasn1W1fDyQVXsVC6MdNuqw2TY2xwotCCYQ8dsfGOu6oQEqIJ0WW3mEy+bv632mpowMcGfkP9hR2+TWx+u0oqyzVADnCi0IJBCQ3w8Y67qhASogZINfPk96L4lWiS480vnnlkjjJrveRpXi/tEspZZeQd4KLQgiENDfBxjruqEBKiBhvYvz9EdoZcR7cTs8uTlonvs31rYI2EKYVys5HGOhNwotCCAQ1t8DGOu6oQEqIB7II+9prcgtoyHEqPnliyra6vpwxqB31hWKDJzP1D+sCi0IHhCN4AEY67qhASogGVk5fiPoW5X1VOSii0+fVmjVa+Z1A8r4bDYoNDQGf1gKLAgcEI1gGOu6oQEqIK5DsY2tgAruXOy8FrccgYWqF6dFTsk+TB5nafcANi+LCiwIGhCNIBjruqEBKiBqagBwMFpo1AE2B/isc6rdMDmBxQkSErisY+0KarjJ6gosCBgQjRAY67qhASogTUYQMW+7MgKqMl/shQkHh3Fi/yEk2VUDEBK9xXOynI8KLAgWEI0IGOu6oQEqIDsuPucVHLcHCBh+VXy9LgzSG316T+xVEhWBs1DHeFj/CiwIFBCNBBjruqEBKiBewqe6zH2giTCqvU3/7HTGrC5/6ElcjbcWrgH7kSjjowosCBIQjQIY67qhASogwqCW1X4kSFj+wboo0R5YhXuzGHY5f0YsMsrAEZonTWMKLAgQEI0BGOu6oQEiIDHCylNPfwGe/Ilh4bRdjcOApuIDSISJefOp59+1GOU2CisIDBA7GIzBlgEqIMDb6Q2Gvbj34w1nwprY+LGWD09EiBauIJEIIb598z6mCisIChAbGIzBlgEiICg34LshAO47dcHHHOiJf5xQkUu0VsSxtqsPcjoczWhCCioICBAQGPz9JSogQ70tKx+ygOMZZSgKvz+vHoVboWUiONMYsa9NBk4Kz2QKKggGEAgY/P0lKiBhMRyrUapSbv0Dr3BhBae12D+5R6pr9VKKMhXIrc4PDAoqCAQQBBj8/SUiIIMV34E/ltIyClMWEAq/pdk6e+No7XAyoQgzrVLpNhg4CioIAhACGPz9JSIgS+4ACDIPgMy9rHHagRckXHyLd11EwOOtpIi3YHrroUUaMQoJ/wAAAAAAAAABEiCzTL8H6AlqvGdqWp4kDFGKsas0KeKjzs+0QxuWahDhtRiG/wE="
			},
			{
				type: "multistore",
				key: "b3JhY2xl",
				data: "vQQKugQKEwoIZXZpZGVuY2USBwoFCM21uAEKMAoDYWNjEikKJwjNtbgBEiCpnMcSKJ6M1amkAImRZTArIRHOeu/50Jo5rO3rvJfCqAozCgZvcmFjbGUSKQonCM21uAESIE+QC4tCXPhasqHtKQfUgwv2dHA8ZJVIR9qu6bgaCboxCjQKB3N0YWtpbmcSKQonCM21uAESIETcYTgf66RphwCHSdu3yuAVbSV9glVZXdFUnxmNnH5FCjAKA2dvdhIpCicIzbW4ARIgYM0TfBli7KxhY4nWgDSDPykhUJwtKFql9RU5l86WinQKMwoGcGFyYW1zEikKJwjNtbgBEiC/Owz4hHzwDasNUjeI+xGUW6ma+mq3LSz8KBe9+fZmZgoSCgd1cGdyYWRlEgcKBQjNtbgBCjEKBG1pbnQSKQonCM21uAESIP4eMZYd/UsxicOW2AUQ5QYXGhOl3uk4V6nnux9mciljCjEKBG1haW4SKQonCM21uAESILjQXXrLNbKLc8dmAUD3qSPSHHXTU8102Vf0C6PEQhfICjMKBnN1cHBseRIpCicIzbW4ARIgjgXbnPGDvS0BwswS3v1DkohXxqQ7kEZ36pSVKr507CEKNQoIc2xhc2hpbmcSKQonCM21uAESIOieO4Oh/Z+x7IC4cN3ORjl3K1s3olg5HZ/OvRuWVi6oCjkKDGRpc3RyaWJ1dGlvbhIpCicIzbW4ARIgbU8UcR/jBa7x3svYjLE/jB2hrPfYS1V1mbI0z2nMZFg="
			}
		]
	},
	height: "3021517",
	codespace: ""
}
*/

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
	iavlOpData := base64ToBytes("rAgKqQgKLggwEKnlwwMYzbW4ASIguYvJR88nWJD/ZP+u7z3k9i117OAXaX/ytarsSCym/m8KLgguEKnd9gIYzLW4ASIgV7NitOHbf1qu1Do0E4sCHUSmrlsrSjtuBDnIsZMdYw4KLggsEObxwAEYzLW4ASIgN+xxrBpxT4VVX9lOVR632Wt+ELRxSM268b9Qs7JFqKIKLQgqEOXhShjMtbgBKiCnvFfCkSTeeOd5jLasn1W1fDyQVXsVC6MdNuqw2TY2xwotCCYQ8dsfGOu6oQEqIJ0WW3mEy+bv632mpowMcGfkP9hR2+TWx+u0oqyzVADnCi0IJBCQ3w8Y67qhASogZINfPk96L4lWiS480vnnlkjjJrveRpXi/tEspZZeQd4KLQgiENDfBxjruqEBKiBhvYvz9EdoZcR7cTs8uTlonvs31rYI2EKYVys5HGOhNwotCCAQ1t8DGOu6oQEqIB7II+9prcgtoyHEqPnliyra6vpwxqB31hWKDJzP1D+sCi0IHhCN4AEY67qhASogGVk5fiPoW5X1VOSii0+fVmjVa+Z1A8r4bDYoNDQGf1gKLAgcEI1gGOu6oQEqIK5DsY2tgAruXOy8FrccgYWqF6dFTsk+TB5nafcANi+LCiwIGhCNIBjruqEBKiBqagBwMFpo1AE2B/isc6rdMDmBxQkSErisY+0KarjJ6gosCBgQjRAY67qhASogTUYQMW+7MgKqMl/shQkHh3Fi/yEk2VUDEBK9xXOynI8KLAgWEI0IGOu6oQEqIDsuPucVHLcHCBh+VXy9LgzSG316T+xVEhWBs1DHeFj/CiwIFBCNBBjruqEBKiBewqe6zH2giTCqvU3/7HTGrC5/6ElcjbcWrgH7kSjjowosCBIQjQIY67qhASogwqCW1X4kSFj+wboo0R5YhXuzGHY5f0YsMsrAEZonTWMKLAgQEI0BGOu6oQEiIDHCylNPfwGe/Ilh4bRdjcOApuIDSISJefOp59+1GOU2CisIDBA7GIzBlgEqIMDb6Q2Gvbj34w1nwprY+LGWD09EiBauIJEIIb598z6mCisIChAbGIzBlgEiICg34LshAO47dcHHHOiJf5xQkUu0VsSxtqsPcjoczWhCCioICBAQGPz9JSogQ70tKx+ygOMZZSgKvz+vHoVboWUiONMYsa9NBk4Kz2QKKggGEAgY/P0lKiBhMRyrUapSbv0Dr3BhBae12D+5R6pr9VKKMhXIrc4PDAoqCAQQBBj8/SUiIIMV34E/ltIyClMWEAq/pdk6e+No7XAyoQgzrVLpNhg4CioIAhACGPz9JSIgS+4ACDIPgMy9rHHagRckXHyLd11EwOOtpIi3YHrroUUaMQoJ/wAAAAAAAAABEiCzTL8H6AlqvGdqWp4kDFGKsas0KeKjzs+0QxuWahDhtRiG/wE=")
	amino.NewCodec().MustUnmarshalBinaryLengthPrefixed(iavlOpData, &iavlProof)
	paths := GetIAVLMerklePaths(&iavlProof)
	expectOracleMerkleHash := hexToBytes("4F900B8B425CF85AB2A1ED2907D4830BF674703C64954847DAAEE9B81A09BA31")
	key := base64ToBytes("/wAAAAAAAAAB")
	value := base64ToBytes("AAAAAAAAAAAAAAABAAAAFgAAAANMUkMAAAADVVNEAAAAAAAPQkAAAAAAAAAABwAAAAAAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAXyquiQAAAABfKq6PAAAAAQAAAAgAAAAAAAH1Ug==")
	leafNode := []byte{0}                                                           // Height
	leafNode = append(leafNode, 2)                                                  // subtree size = 1 * 2
	leafNode = append(leafNode, encodeVarint(iavlProof.Proof.Leaves[0].Version)...) // version
	leafNode = append(leafNode, uint8(len(key)))                                    // key length
	leafNode = append(leafNode, key...)                                             // key to result of request#1
	leafNode = append(leafNode, 32)                                                 // size of result hash must be 32
	leafNode = append(leafNode, tmhash.Sum(value)...)                               // value on this key is a result hash
	currentHash := tmhash.Sum(leafNode)
	for _, path := range paths {
		currentHash = getParentHash(path, currentHash)
	}
	require.Equal(t, expectOracleMerkleHash, currentHash)
}
