package goborsh

import (
	"encoding/binary"
)

// BorshEncoder stores the information necessary for the encoding functions
type BorshEncoder struct {
	data []byte
}

// NewBorshEncoder returns an empty BorshEncoder struct
func NewBorshEncoder() BorshEncoder {
	return BorshEncoder{
		data: []byte{},
	}
}

// EncodeU8 takes a `uint8` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeU8(value uint8) []byte {
	bytes := []byte{value}
	encoder.data = append(encoder.data, bytes...)
	return bytes
}

// EncodeU16 takes a `uint16` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeU16(value uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, value)
	encoder.data = append(encoder.data, bytes...)
	return bytes
}

// EncodeU32 takes a `uint32` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeU32(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	encoder.data = append(encoder.data, bytes...)
	return bytes
}

// EncodeU64 takes a `uint64` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeU64(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	encoder.data = append(encoder.data, bytes...)
	return bytes
}

// EncodeSigned8 takes a `int8` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeSigned8(value int8) []byte {
	uintValue := uint8(value)
	return encoder.EncodeU8(uintValue)
}

// EncodeSigned16 takes a `int8` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeSigned16(value int16) []byte {
	uintValue := uint16(value)
	return encoder.EncodeU16(uintValue)
}

// EncodeSigned32 takes a `int32` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeSigned32(value int32) []byte {
	uintValue := uint32(value)
	return encoder.EncodeU32(uintValue)
}

// EncodeSigned64 takes a `uint64` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeSigned64(value int64) []byte {
	uintValue := uint64(value)
	return encoder.EncodeU64(uintValue)
}

// EncodeBytes takes a `[]byte` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeBytes(value []byte) []byte {
	encoder.EncodeU32(uint32(len(value)))
	encoder.data = append(encoder.data, value...)
	return value
}

// EncodeString takes a `string` variable and encodes it into a byte array
func (encoder *BorshEncoder) EncodeString(value string) []byte {
	stringBytes := []byte(value)
	encoder.EncodeU32(uint32(len(stringBytes)))
	encoder.data = append(encoder.data, stringBytes...)
	return stringBytes
}
