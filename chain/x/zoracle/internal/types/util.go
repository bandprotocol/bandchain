package types

import (
	"bytes"
	"encoding/binary"
)

func int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	// Encode int64 must not have error
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func int64ToBytes(num int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(num))
	return result
}
