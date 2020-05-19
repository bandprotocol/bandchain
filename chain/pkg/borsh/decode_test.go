package borsh

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func getDecoder(inputString string) Decoder {
	data, _ := hex.DecodeString(inputString)
	return NewDecoder(data)
}

func TestFinishedExpectFalse(t *testing.T) {
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewDecoder(data)
	decodeFinished := decoder.Finished()

	require.False(t, decodeFinished, `Incorrect result for Finished boolean: expected "false"`)
}

func TestFinishedExpectTrue(t *testing.T) {
	data, _ := hex.DecodeString("10")
	decoder := Decoder{data, 1}
	decodeFinished := decoder.Finished()

	require.True(t, decodeFinished, `Incorrect result for Finished boolean: expected "true"`)
}

func TestDecodeU8(t *testing.T) {
	var expectedU8 uint8 = uint8(112)
	decoder := getDecoder("70")
	resultU8, err := decoder.DecodeU8()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedU8, resultU8, `Incorrect result for decoded uint8: expected "112"`)
}

func TestDecodeU8Error(t *testing.T) {
	data, _ := hex.DecodeString("70")
	decoder := Decoder{data, 8}
	_, err := decoder.DecodeU8()

	require.Error(t, err, `Incorrect result for result message`)
}

func TestDecodeU16(t *testing.T) {
	var expectedU16 uint16 = uint16(2345)
	decoder := getDecoder("2909")
	resultU16, err := decoder.DecodeU16()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedU16, resultU16, `Incorrect result for decoded uint16: expected "2345"`)
}

func TestDecodeU16Error(t *testing.T) {
	data, _ := hex.DecodeString("2909")
	decoder := Decoder{data, 16}
	_, err := decoder.DecodeU16()

	require.Error(t, err, `Incorrect result for result message`)
}

func TestDecodeU32(t *testing.T) {
	var expectedU32 uint32 = uint32(12345)
	decoder := getDecoder("39300000")
	resultU32, err := decoder.DecodeU32()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedU32, resultU32, `Incorrect result for decoded uint32: expected "12345"`)
}

func TestDecodeU32Error(t *testing.T) {
	data, _ := hex.DecodeString("39300000")
	decoder := Decoder{data, 32}
	_, err := decoder.DecodeU32()

	require.Error(t, err, `Incorrect result for result message`)
}

func TestDecodeU64(t *testing.T) {
	var expectedU64 uint64 = uint64(50)
	decoder := getDecoder("3200000000000000")
	resultU64, err := decoder.DecodeU64()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedU64, resultU64, `Incorrect result for decoded uint64: expected "50"`)
}

func TestDecodeU64Error(t *testing.T) {
	data, _ := hex.DecodeString("3200000000000000")
	decoder := Decoder{data, 64}
	_, err := decoder.DecodeU64()

	require.Error(t, err, `Incorrect result for result message`)
}

func TestDecodeI8Positive(t *testing.T) {
	var expectedI8 int8 = int8(112)
	decoder := getDecoder("70")
	resultI8, err := decoder.DecodeI8()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI8, resultI8, `Incorrect result for decoded int8: expected "112"`)
}

func TestDecodeI8Negative(t *testing.T) {
	var expectedI8 int8 = int8(-2)
	decoder := getDecoder("FE")

	resultI8, err := decoder.DecodeI8()
	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI8, resultI8, `Incorrect result for decoded int8: expected "-2"`)
}

func TestDecodeI16Positive(t *testing.T) {
	var expectedI16 int16 = int16(2345)
	decoder := getDecoder("2909")
	resultI16, err := decoder.DecodeI16()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI16, resultI16, `Incorrect result for decoded int16: expected "2345"`)
}

func TestDecodeI16Negative(t *testing.T) {
	var expectedI16 int16 = int16(-2)
	decoder := getDecoder("FEFF")
	resultI16, err := decoder.DecodeI16()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI16, resultI16, `Incorrect result for decoded int16: expected "-2"`)
}

func TestDecodeI32Positive(t *testing.T) {
	var expectedI32 int32 = int32(12345)
	decoder := getDecoder("39300000")
	resultI32, err := decoder.DecodeI32()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI32, resultI32, `Incorrect result for decoded int32: expected "12345"`)
}

func TestDecodeI32Negative(t *testing.T) {
	var expectedI32 int32 = int32(-70)
	decoder := getDecoder("BAFFFFFF")
	resultI32, err := decoder.DecodeI32()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI32, resultI32, `Incorrect result for decoded int32: expected "-70"`)
}

func TestDecodeI64Positive(t *testing.T) {
	var expectedI64 int64 = int64(50)
	decoder := getDecoder("3200000000000000")
	resultI64, err := decoder.DecodeI64()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI64, resultI64, `Incorrect result for decoded int364: expected "50"`)
}

func TestDecodeI64Negative(t *testing.T) {
	var expectedI64 int64 = int64(-20486)
	decoder := getDecoder("FAAFFFFFFFFFFFFF")
	resultI64, err := decoder.DecodeI64()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedI64, resultI64, `Incorrect result for decoded int64: expected "-20486"`)
}

func TestDecodeBytes(t *testing.T) {
	var expectedBytes []byte = []byte{66, 84, 67}
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewDecoder(data)
	resultBytes, err := decoder.DecodeBytes()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedBytes, resultBytes, `Incorrect result for decoded bytes: expected "{66, 84, 67}"`)
}

func TestDecodeBytesErrorDecodeU32(t *testing.T) {
	var expectedBytes []byte
	data, _ := hex.DecodeString("03000000425443")
	decoder := Decoder{data, 64}
	resultBytes, err := decoder.DecodeBytes()

	require.Error(t, err, `Incorrect result for return error: expected "Borsh: out of range"`)
	require.Equal(t, expectedBytes, resultBytes, `Incorrect result for decoded bytes: expected "{}"`)
}

func TestDecodeString(t *testing.T) {
	expectedString := "BTC"
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewDecoder(data)
	resultString, err := decoder.DecodeString()

	require.Nil(t, err, `Incorrect result for return error: expected "nil"`)
	require.Equal(t, expectedString, resultString, `Incorrect result for decoded bytes: expected "BTC"`)
}

func TestDecodeStringErrorDecodeU32(t *testing.T) {
	expectedString := ""
	data, _ := hex.DecodeString("03000000425443")
	decoder := Decoder{data, 64}
	resultString, err := decoder.DecodeString()

	require.Error(t, err, `Incorrect result for return error:`)
	require.Equal(t, expectedString, resultString, `Incorrect result for decoded bytes: expected ""`)
}
