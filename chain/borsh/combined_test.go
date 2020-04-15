package goborsh

import (
	"bytes"
	"testing"
)

func TestEncodeDecodeBytes(t *testing.T) {
	// Declare bytes to use to encode/decode
	firstBytes := []byte{0}
	secondBytes := []byte{0,5}
	thirdBytes := []byte{123,0,0,43,12,0,123}
	fourthBytes := []byte{}
	var fifthBytes []byte = nil

	encoder := NewBorshEncoder()

	// Test bytes returned from encode functions
	returnBytes := encoder.EncodeBytes(firstBytes)
	if bytes.Compare(returnBytes, []byte(firstBytes)) == 1 {
		t.Errorf(`Incorrect return first byte, expected "%s" got "%s" `, firstBytes, returnBytes)
	}
	returnBytes = encoder.EncodeBytes(secondBytes)
	if bytes.Compare(returnBytes, []byte(secondBytes)) == 1 {
		t.Errorf(`Incorrect return second byte, expected "%s" got "%s" `, firstBytes, returnBytes)
	}
	returnBytes = encoder.EncodeBytes(thirdBytes)
	if bytes.Compare(returnBytes, []byte(thirdBytes)) == 1 {
		t.Errorf(`Incorrect return third byte, expected "%s" got "%s" `, firstBytes, returnBytes)
	}
	returnBytes = encoder.EncodeBytes(fourthBytes)
	if bytes.Compare(returnBytes, []byte(fourthBytes)) == 1 {
		t.Errorf(`Incorrect return fourth byte, expected "%s" got "%s" `, firstBytes, returnBytes)
	}
	returnBytes = encoder.EncodeBytes(fifthBytes)
	if bytes.Compare(returnBytes, []byte(fifthBytes)) == 1 {
		t.Errorf(`Incorrect return fifth byte, expected "%s" got "%s" `, firstBytes, returnBytes)
	}

	decoder := NewBorshDecoder(encoder.data)

	// Test bytes returned when passing encoded data into decoder
	decodedBytes, err :=decoder.DecodeBytes()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if bytes.Compare(decodedBytes, firstBytes ) == 1{
		t.Errorf(`Incorrect decoded first bytes, expected "%s" got "%s" `, firstBytes, decodedBytes)
	}
	decodedBytes, err =decoder.DecodeBytes()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if bytes.Compare(decodedBytes, secondBytes) == 1 {
		t.Errorf(`Incorrect decoded second bytes, expected "%s" got "%s" `, secondBytes, decodedBytes)
	}
	decodedBytes, err =decoder.DecodeBytes()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if bytes.Compare(decodedBytes, thirdBytes ) == 1{
		t.Errorf(`Incorrect decoded third bytes, expected "%s" got "%s" `, thirdBytes, decodedBytes)
	}
	decodedBytes, err =decoder.DecodeBytes()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if bytes.Compare(decodedBytes, fourthBytes ) == 1{
		t.Errorf(`Incorrect decoded fourth bytes, expected "%s" got "%s" `, fourthBytes, decodedBytes)
	}
	decodedBytes, err =decoder.DecodeBytes()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if bytes.Compare(decodedBytes, fifthBytes ) == 1{
		t.Errorf(`Incorrect decoded fifth bytes, expected "%s" got "%s" `, fifthBytes, decodedBytes)
	}
}
