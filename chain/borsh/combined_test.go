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

func TestEncodeDecodeString(t *testing.T) {
	// Declare strings to use to encode/decode
	firstString := "BAND"
	secondString := "BTC"
	thirdString := "ETH"
	fourthString := ""
	fifthString := "!@#$%^&*(ERTEYRUTIKJ!$#@%^&*&("

	encoder := NewBorshEncoder()

	// Test strings returned from encode functions
	returnString := encoder.EncodeString(firstString)
	if bytes.Compare(returnString, []byte(firstString)) == 1 {
		t.Errorf(`Incorrect return first string, expected "%s" got "%s" `, firstString, returnString)
	}
	returnString = encoder.EncodeString(secondString)
	if bytes.Compare(returnString, []byte(secondString)) == 1 {
		t.Errorf(`Incorrect return second string, expected "%s" got "%s" `, firstString, returnString)
	}
	returnString = encoder.EncodeString(thirdString)
	if bytes.Compare(returnString, []byte(thirdString)) == 1 {
		t.Errorf(`Incorrect return third string, expected "%s" got "%s" `, firstString, returnString)
	}
	returnString = encoder.EncodeString(fourthString)
	if bytes.Compare(returnString, []byte(thirdString)) == 1 {
		t.Errorf(`Incorrect return fourth string, expected "%s" got "%s" `, firstString, returnString)
	}
	returnString = encoder.EncodeString(fifthString)
	if bytes.Compare(returnString, []byte(fifthString)) == 1 {
		t.Errorf(`Incorrect return fifth string, expected "%s" got "%s" `, firstString, returnString)
	}

	decoder := NewBorshDecoder(encoder.data)

	// Test strings returned when passing encoded data into decoder
	decodedString, err :=decoder.DecodeString()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if decodedString != firstString {
		t.Errorf(`Incorrect decoded first string, expected "%s" got "%s" `, firstString, decodedString)
	}
	decodedString, err =decoder.DecodeString()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if decodedString != secondString {
		t.Errorf(`Incorrect decoded second string, expected "%s" got "%s" `, secondString, decodedString)
	}
	decodedString, err =decoder.DecodeString()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if decodedString != thirdString {
		t.Errorf(`Incorrect decoded third string, expected "%s" got "%s" `, thirdString, decodedString)
	}
	decodedString, err =decoder.DecodeString()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if decodedString != fourthString {
		t.Errorf(`Incorrect decoded fourth string, expected "%s" got "%s" `, fourthString, decodedString)
	}
	decodedString, err =decoder.DecodeString()
	if err != nil {
		t.Errorf(`Incorrect decoded error type, expected "%v" got "%v" `, nil, err)
	}
	if decodedString != fifthString {
		t.Errorf(`Incorrect decoded fourth string, expected "%s" got "%s" `, fifthString, decodedString)
	}
}
