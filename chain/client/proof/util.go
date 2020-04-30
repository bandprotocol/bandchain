package proof

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Copied from https://github.com/tendermint/tendermint/blob/master/types/utils.go
func cdcEncode(cdc *codec.Codec, item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		return cdc.MustMarshalBinaryBare(item)
	}
	return nil
}

// Go lacks a simple and safe way to see if something is a typed nil.
// See:
//  - https://dave.cheney.net/2017/08/09/typed-nils-in-go-2
//  - https://groups.google.com/forum/#!topic/golang-nuts/wnH302gBa4I/discussion
//  - https://github.com/golang/go/issues/21538
func isTypedNil(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// Returns true if it has zero length.
func isEmpty(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	default:
		return false
	}
}
func base64ToBytes(s string) []byte {
	decodedString, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return decodedString
}

func encodeStoreMerkleHash(key string, value []byte) []byte {
	bytesKey := []byte(key)
	keyBytes := append([]byte{uint8(len(bytesKey))}, bytesKey...)
	valueBytes := append([]byte{32}, tmhash.Sum(tmhash.Sum(value))...)
	return append(keyBytes, valueBytes...)
}

func encodeVarint(value int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, value)
	return buf[:n]
}

func mustParseInt64(b []byte) int64 {
	i64, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

func mustDecodeString(hexstr string) []byte {
	b, err := hex.DecodeString(hexstr)
	if err != nil {
		panic(err)
	}
	return b
}
