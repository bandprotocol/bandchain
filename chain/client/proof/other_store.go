package proof

import "github.com/cosmos/cosmos-sdk/store/rootmulti"

type MultiStoreProof struct{}

type MultiStoreProofEthereum struct{}

func (multiStoreProof *MultiStoreProof) encodeToEthFormat() MultiStoreProofEthereum {
	return MultiStoreProofEthereum{}
}

func GetMultiStoreProof(proof rootmulti.MultiStoreProofOp) MultiStoreProof {
	return MultiStoreProof{}
}

// func MakeOtherStoresMerkleHash(mspo rootmulti.MultiStoreProofOp) (tmbytes.HexBytes, tmbytes.HexBytes, tmbytes.HexBytes) {
// 	var oracleHash []byte
// 	m := map[string][]byte{}
// 	for _, si := range mspo.Proof.StoreInfos {
// 		m[si.Name] = tmhash.Sum(tmhash.Sum(si.Core.CommitID.Hash))
// 		if si.Name == oracle.ModuleName {
// 			oracleHash = si.Core.CommitID.Hash
// 		}
// 	}

// 	keys := []string{}
// 	for k := range m {
// 		if k != "oracle" {
// 			keys = append(keys, k)
// 		}
// 	}
// 	sort.Strings(keys)

// 	bs := [][]byte{}
// 	for _, k := range keys {
// 		bs = append(bs, m[k])
// 	}

// 	h1 := tmhash.Sum(
// 		append(
// 			[]byte{1},
// 			append(
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("acc")...), append([]byte{32}, m["acc"]...)...)...)),
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{12}, []byte("distribution")...), append([]byte{32}, m["distribution"]...)...)...))...,
// 			)...,
// 		),
// 	)

// 	h2 := tmhash.Sum(
// 		append(
// 			[]byte{1},
// 			append(
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("gov")...), append([]byte{32}, m["gov"]...)...)...)),
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{4}, []byte("main")...), append([]byte{32}, m["main"]...)...)...))...,
// 			)...,
// 		),
// 	)

// 	h3 := tmhash.Sum(
// 		append(
// 			[]byte{1},
// 			append(
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{4}, []byte("mint")...), append([]byte{32}, m["mint"]...)...)...)),
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("params")...), append([]byte{32}, m["params"]...)...)...))...,
// 			)...,
// 		),
// 	)

// 	h4 := tmhash.Sum(
// 		append(
// 			[]byte{1},
// 			append(
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{8}, []byte("slashing")...), append([]byte{32}, m["slashing"]...)...)...)),
// 				tmhash.Sum(append([]byte{0}, append(append([]byte{7}, []byte("staking")...), append([]byte{32}, m["staking"]...)...)...))...,
// 			)...,
// 		),
// 	)

// 	h5 := tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("supply")...), append([]byte{32}, m["supply"]...)...)...))

// 	h6 := tmhash.Sum(append([]byte{1}, append(h1, h2...)...))

// 	h7 := tmhash.Sum(append([]byte{1}, append(h3, h4...)...))

// 	h8 := tmhash.Sum(append([]byte{1}, append(h6, h7...)...))

// 	return h5, h8, oracleHash
// }
