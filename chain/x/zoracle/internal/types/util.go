package types

import (
	"bytes"
	"encoding/binary"
)

func uint64ToBytes(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	// Encode uint64 must not have error
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	// Encode int64 must not have error
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
