package borsh

import (
	"encoding/binary"
	"errors"
)

// Decoder stores the information necessary for the decoding functions
type Decoder struct {
	data   []byte
	offset uint32
}

// NewDecoder returns an empty Decoder struct
func NewDecoder(data []byte) Decoder {
	return Decoder{
		data:   data,
		offset: 0,
	}
}

// Finished returns a bool on whether the decoding has finished
func (decoder *Decoder) Finished() bool {
	return decoder.offset == uint32(len(decoder.data))
}

// DecodeU8 decodes the input bytes and returns the corresponding
// `uint8` value and any errors
func (decoder *Decoder) DecodeU8() (uint8, error) {
	if uint32(len(decoder.data)) < decoder.offset+1 {
		return 0, errors.New("Borsh: out of range")
	}
	val := uint8(decoder.data[decoder.offset])
	decoder.offset++
	return val, nil
}

// DecodeU16 decodes the input bytes and returns the corresponding
// `uint8` value and any errors
func (decoder *Decoder) DecodeU16() (uint16, error) {
	if uint32(len(decoder.data)) < decoder.offset+2 {
		return 0, errors.New("Borsh: out of range")
	}
	val := binary.LittleEndian.Uint16(decoder.data[decoder.offset : decoder.offset+2])
	decoder.offset += 2
	return val, nil
}

// DecodeU32 decodes the input bytes and returns the corresponding
// `uint32` value and any errors
func (decoder *Decoder) DecodeU32() (uint32, error) {
	if uint32(len(decoder.data)) < decoder.offset+4 {
		return 0, errors.New("Borsh: out of range")
	}
	val := binary.LittleEndian.Uint32(decoder.data[decoder.offset : decoder.offset+4])
	decoder.offset += 4
	return val, nil
}

// DecodeU64 decodes the input bytes and returns the corresponding
// `uint64` value and any errors
func (decoder *Decoder) DecodeU64() (uint64, error) {
	if uint32(len(decoder.data)) < decoder.offset+8 {
		return 0, errors.New("Borsh: out of range")
	}
	val := binary.LittleEndian.Uint64(decoder.data[decoder.offset : decoder.offset+8])
	decoder.offset += 8
	return val, nil
}

// DecodeI8 decodes the input bytes and returns the corresponding signed
// `int8` value and any errors
func (decoder *Decoder) DecodeI8() (int8, error) {
	unsigned, err := decoder.DecodeU8()
	if err != nil {
		return 0, err
	}
	return int8(unsigned), nil
}

// DecodeI16 decodes the input bytes and returns the corresponding signed
// `int16` value and any errors
func (decoder *Decoder) DecodeI16() (int16, error) {
	unsigned, err := decoder.DecodeU16()
	if err != nil {
		return 0, err
	}
	return int16(unsigned), nil
}

// DecodeI32 decodes the input bytes and returns the corresponding signed
// `int32` value and any errors
func (decoder *Decoder) DecodeI32() (int32, error) {
	unsigned, err := decoder.DecodeU32()
	if err != nil {
		return 0, err
	}
	return int32(unsigned), nil
}

// DecodeI64 decodes the input bytes and returns the corresponding signed
// `int64` value and any errors
func (decoder *Decoder) DecodeI64() (int64, error) {
	unsigned, err := decoder.DecodeU64()
	if err != nil {
		return 0, err
	}
	return int64(unsigned), nil
}

// DecodeBytes decodes the input bytes and returns the corresponding
// `[]bytes` slice and any errors
func (decoder *Decoder) DecodeBytes() ([]byte, error) {
	length, err := decoder.DecodeU32()
	if err != nil {
		return nil, err
	}
	if uint32(len(decoder.data)) < decoder.offset+length {
		return nil, errors.New("Borsh: out of range")
	}
	val := decoder.data[decoder.offset : decoder.offset+uint32(length)]
	decoder.offset += length
	return val, nil
}

// DecodeString decodes the input bytes and returns the corresponding
// `string` and any errors
func (decoder *Decoder) DecodeString() (string, error) {
	length, err := decoder.DecodeU32()
	if err != nil {
		return "", err
	}
	if uint32(len(decoder.data)) < decoder.offset+length {
		return "", errors.New("Borsh: out of range")
	}
	val := string(decoder.data[decoder.offset : decoder.offset+uint32(length)])
	decoder.offset += length
	return val, nil
}
