package goborsh

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

func getDecoder(inputString string) BorshDecoder {
	data, _ := hex.DecodeString(inputString)
	return NewBorshDecoder(data)
}

func TestFinishedExpectFalse(t *testing.T) {
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewBorshDecoder(data)
	decodeFinished := decoder.Finished()

	if decodeFinished {
		t.Errorf(`Incorrect result for Finished bool: expected "%t" got "%t"`, false, decodeFinished)
	}
}

func TestFinishedExpectTrue(t *testing.T) {
	data, _ := hex.DecodeString("10")
	decoder := BorshDecoder{data, 1}
	fmt.Println(len(decoder.data))
	decodeFinished := decoder.Finished()

	if !decodeFinished {
		t.Errorf(`Incorrect result for Finished bool: expected "%t" got "%t"`, true, decodeFinished)
	}
}

func TestDecodeU8(t *testing.T) {
	var expectedU8 uint8 = 112
	decoder := getDecoder("70")
	resultU8, err := decoder.DecodeU8()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultU8 != expectedU8 {
		t.Errorf(`Incorrect decoded uint8: expected "%d" got "%d"`, expectedU8, resultU8)
	}

}

func TestDecodeU8Error(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	data, _ := hex.DecodeString("70")
	decoder := BorshDecoder{data, 8}
	_, err := decoder.DecodeU8()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}

func TestDecodeU16(t *testing.T) {
	var expectedU16 uint16 = 2345
	decoder := getDecoder("2909")
	resultU16, err := decoder.DecodeU16()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultU16 != expectedU16 {
		t.Errorf(`Incorrect decoded uint16: expected "%d" got "%d"`, expectedU16, resultU16)
	}
}

func TestDecodeU16Error(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	data, _ := hex.DecodeString("2909")
	decoder := BorshDecoder{data, 16}
	_, err := decoder.DecodeU16()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}

func TestDecodeU32(t *testing.T) {
	var expectedU32 uint32 = 12345
	decoder := getDecoder("39300000")
	resultU32, err := decoder.DecodeU32()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultU32 != expectedU32 {
		t.Errorf(`Incorrect decoded uint32: expected "%d" got "%d"`, expectedU32, resultU32)
	}
}

func TestDecodeU32Error(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	data, _ := hex.DecodeString("39300000")
	decoder := BorshDecoder{data, 32}
	_, err := decoder.DecodeU32()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}

func TestDecodeU64(t *testing.T) {
	var expectedU64 uint64 = 50
	decoder := getDecoder("3200000000000000")
	resultU64, err := decoder.DecodeU64()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultU64 != expectedU64 {
		t.Errorf(`Incorrect decoded uint64: expected "%d" got "%d"`, expectedU64, resultU64)
	}
}

func TestDecodeU64Error(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	data, _ := hex.DecodeString("3200000000000000")
	decoder := BorshDecoder{data, 64}
	_, err := decoder.DecodeU64()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}

func TestDecodeSigned8Positive(t *testing.T) {
	var expectedSigned64 int8 = 112
	decoder := getDecoder("70")
	resultSigned8, err := decoder.DecodeSigned8()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned8 != expectedSigned64 {
		t.Errorf(`Incorrect decoded int8: expected "%d" got "%d"`, expectedSigned64, resultSigned8)
	}
}

func TestDecodeSigned8Negative(t *testing.T) {
	var expectedSigned64 int8 = -2
	decoder := getDecoder("FE")
	resultSigned8, err := decoder.DecodeSigned8()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned8 != expectedSigned64 {
		t.Errorf(`Incorrect decoded int8: expected "%d" got "%d"`, expectedSigned64, resultSigned8)
	}
}

func TestDecodeSigned16Positive(t *testing.T) {
	var expectedSigned16 int16 = 2345
	decoder := getDecoder("2909")
	resultSigned16, err := decoder.DecodeSigned16()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned16 != expectedSigned16 {
		t.Errorf(`Incorrect decoded int16: expected "%d" got "%d"`, expectedSigned16, resultSigned16)
	}

}

func TestDecodeSigned16Negative(t *testing.T) {
	var expectedSigned16 int16 = -2
	decoder := getDecoder("FEFF")
	resultSigned16, err := decoder.DecodeSigned16()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned16 != expectedSigned16 {
		t.Errorf(`Incorrect decoded int16: expected "%d" got "%d"`, expectedSigned16, resultSigned16)
	}
}

func TestDecodeSigned32Positive(t *testing.T) {
	var expectedSigned32 int32 = 12345
	decoder := getDecoder("39300000")
	resultSigned32, err := decoder.DecodeSigned32()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned32 != expectedSigned32 {
		t.Errorf(`Incorrect decoded int32: expected "%d" got "%d"`, expectedSigned32, resultSigned32)
	}
}

func TestDecodeSigned32Negative(t *testing.T) {
	var expectedSigned32 int32 = -70
	decoder := getDecoder("BAFFFFFF")
	resultSigned32, err := decoder.DecodeSigned32()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned32 != expectedSigned32 {
		t.Errorf(`Incorrect decoded int32: expected "%d" got "%d"`, expectedSigned32, resultSigned32)
	}
}

func TestDecodeSigned64Positive(t *testing.T) {
	var expectedSigned64 int64 = 50
	decoder := getDecoder("3200000000000000")
	resultSigned64, err := decoder.DecodeSigned64()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned64 != expectedSigned64 {
		t.Errorf(`Incorrect decoded int64: expected "%d" got "%d"`, expectedSigned64, resultSigned64)
	}
}

func TestDecodeSigned64Negative(t *testing.T) {
	var expectedSigned64 int64 = -20486
	decoder := getDecoder("FAAFFFFFFFFFFFFF")
	resultSigned64, err := decoder.DecodeSigned64()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultSigned64 != expectedSigned64 {
		t.Errorf(`Incorrect decoded int64: expected "%d" got "%d"`, expectedSigned64, resultSigned64)
	}
}

func TestDecodeBytes(t *testing.T) {
	var expectedBytes []byte = []byte{66, 84, 67}
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewBorshDecoder(data)
	resultBytes, err := decoder.DecodeBytes()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if bytes.Compare(expectedBytes, resultBytes) == 1 {
		t.Errorf(`Incorrect decoded bytes: expected "%v" got "%v"`, expectedBytes, resultBytes)
	}
}

func TestDecodeBytesErrorDecodeU32(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	expectedBytes := []byte{}
	data, _ := hex.DecodeString("03000000425443")
	decoder := BorshDecoder{data, 64}
	resultBytes, err := decoder.DecodeBytes()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if bytes.Compare(expectedBytes, resultBytes) == 1 {
		t.Errorf(`Incorrect decoded bytes: expected "%v" got "%v"`, expectedBytes, resultBytes)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}

func TestDecodeString(t *testing.T) {
	expectedString := "BTC"
	data, _ := hex.DecodeString("03000000425443")
	decoder := NewBorshDecoder(data)
	resultString, err := decoder.DecodeString()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultString != expectedString {
		t.Errorf(`Incorrect decoded string: expected "%s" got "%s"`, expectedString, resultString)
	}
}

func TestDecodeStringErrorDecodeU32(t *testing.T) {
	expectedErr := errors.New("Borsh: out of range")
	expectedString := ""
	data, _ := hex.DecodeString("03000000425443")
	decoder := BorshDecoder{data, 64}
	resultString, err := decoder.DecodeString()

	if err == nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultString != expectedString {
		t.Errorf(`Incorrect decoded string: expected "%s" got "%s"`, expectedString, resultString)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf(`Incorrect return error message: expected "%v" got "%v"`, err.Error(), expectedErr.Error())
	}
}
