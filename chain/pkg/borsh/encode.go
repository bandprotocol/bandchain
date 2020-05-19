package borsh

import (
	"encoding/binary"
)

// Encoder stores the information necessary for the encoding functions
type Encoder struct {
	data []byte
}

// NewEncoder returns an empty Encoder struct
func NewEncoder() Encoder {
	return Encoder{
		data: []byte{},
	}
}

// GetEncodedData returns the `data` byte slice containing all the previously encoded data
func (encoder *Encoder) GetEncodedData() []byte {
	return encoder.data
}

// EncodeU8 takes an `uint8` variable and encodes it into a byte array
func (encoder *Encoder) EncodeU8(value uint8) {
	encoder.data = append(encoder.data, value)
}

// EncodeU16 takes an `uint16` variable and encodes it into a byte array
func (encoder *Encoder) EncodeU16(value uint16) {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, value)
	encoder.data = append(encoder.data, bytes...)
}

// EncodeU32 takes an `uint32` variable and encodes it into a byte array
func (encoder *Encoder) EncodeU32(value uint32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	encoder.data = append(encoder.data, bytes...)
}

// EncodeU64 takes an `uint64` variable and encodes it into a byte array
func (encoder *Encoder) EncodeU64(value uint64) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	encoder.data = append(encoder.data, bytes...)
}

// EncodeSigned8 takes an `int8` variable and encodes it into a byte array
func (encoder *Encoder) EncodeSigned8(value int8) {
	uintValue := uint8(value)
	encoder.EncodeU8(uintValue)
}

// EncodeSigned16 takes an `int16` variable and encodes it into a byte array
func (encoder *Encoder) EncodeSigned16(value int16) {
	uintValue := uint16(value)
	encoder.EncodeU16(uintValue)
}

// EncodeSigned32 takes an `int32` variable and encodes it into a byte array
func (encoder *Encoder) EncodeSigned32(value int32) {
	uintValue := uint32(value)
	encoder.EncodeU32(uintValue)
}

// EncodeSigned64 takes an `int64` variable and encodes it into a byte array
func (encoder *Encoder) EncodeSigned64(value int64) {
	uintValue := uint64(value)
	encoder.EncodeU64(uintValue)
}

// EncodeBytes takes a `[]byte` variable and encodes it into a byte array
func (encoder *Encoder) EncodeBytes(value []byte) {
	encoder.EncodeU32(uint32(len(value)))
	encoder.data = append(encoder.data, value...)
}

// EncodeString takes a `string` variable and encodes it into a byte array
func (encoder *Encoder) EncodeString(value string) {
	stringBytes := []byte(value)
	encoder.EncodeU32(uint32(len(stringBytes)))
	encoder.data = append(encoder.data, stringBytes...)
}
