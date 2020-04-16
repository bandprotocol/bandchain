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

func TestDecodeI8Positive(t *testing.T) {
	var expectedI64 int8 = 112
	decoder := getDecoder("70")
	resultI8, err := decoder.DecodeI8()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI8 != expectedI64 {
		t.Errorf(`Incorrect decoded int8: expected "%d" got "%d"`, expectedI64, resultI8)
	}
}

func TestDecodeI8Negative(t *testing.T) {
	var expectedI64 int8 = -2
	decoder := getDecoder("FE")
	resultI8, err := decoder.DecodeI8()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI8 != expectedI64 {
		t.Errorf(`Incorrect decoded int8: expected "%d" got "%d"`, expectedI64, resultI8)
	}
}

func TestDecodeI16Positive(t *testing.T) {
	var expectedI16 int16 = 2345
	decoder := getDecoder("2909")
	resultI16, err := decoder.DecodeI16()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI16 != expectedI16 {
		t.Errorf(`Incorrect decoded int16: expected "%d" got "%d"`, expectedI16, resultI16)
	}

}

func TestDecodeI16Negative(t *testing.T) {
	var expectedI16 int16 = -2
	decoder := getDecoder("FEFF")
	resultI16, err := decoder.DecodeI16()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI16 != expectedI16 {
		t.Errorf(`Incorrect decoded int16: expected "%d" got "%d"`, expectedI16, resultI16)
	}
}

func TestDecodeI32Positive(t *testing.T) {
	var expectedI32 int32 = 12345
	decoder := getDecoder("39300000")
	resultI32, err := decoder.DecodeI32()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI32 != expectedI32 {
		t.Errorf(`Incorrect decoded int32: expected "%d" got "%d"`, expectedI32, resultI32)
	}
}

func TestDecodeI32Negative(t *testing.T) {
	var expectedI32 int32 = -70
	decoder := getDecoder("BAFFFFFF")
	resultI32, err := decoder.DecodeI32()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI32 != expectedI32 {
		t.Errorf(`Incorrect decoded int32: expected "%d" got "%d"`, expectedI32, resultI32)
	}
}

func TestDecodeI64Positive(t *testing.T) {
	var expectedI64 int64 = 50
	decoder := getDecoder("3200000000000000")
	resultI64, err := decoder.DecodeI64()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI64 != expectedI64 {
		t.Errorf(`Incorrect decoded int64: expected "%d" got "%d"`, expectedI64, resultI64)
	}
}

func TestDecodeI64Negative(t *testing.T) {
	var expectedI64 int64 = -20486
	decoder := getDecoder("FAAFFFFFFFFFFFFF")
	resultI64, err := decoder.DecodeI64()

	if err != nil {
		t.Errorf(`Incorrect return error type: expected "%v" got "%v"`, nil, err)
	}
	if resultI64 != expectedI64 {
		t.Errorf(`Incorrect decoded int64: expected "%d" got "%d"`, expectedI64, resultI64)
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
