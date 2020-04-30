package borsh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeBytes(t *testing.T) {
	// Declare bytes to use to encode/decode
	firstBytes := []byte{0}
	secondBytes := []byte{0, 5}
	thirdBytes := []byte{123, 0, 0, 43, 12, 0, 123}
	fourthBytes := []byte{}

	encoder := NewEncoder()

	// Test bytes returned from encode functions
	encoder.EncodeBytes(firstBytes)
	encoder.EncodeBytes(secondBytes)
	encoder.EncodeBytes(thirdBytes)
	encoder.EncodeBytes(fourthBytes)

	decoder := NewDecoder(encoder.data)

	// Test bytes returned when passing encoded data into decoder
	decodedBytes, err := decoder.DecodeBytes()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, firstBytes, decodedBytes, "Incorrect decoded first bytes")

	decodedBytes, err = decoder.DecodeBytes()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, secondBytes, decodedBytes, "Incorrect decoded second bytes")

	decodedBytes, err = decoder.DecodeBytes()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, thirdBytes, decodedBytes, "Incorrect decoded third bytes")

	decodedBytes, err = decoder.DecodeBytes()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, fourthBytes, decodedBytes, "Incorrect decoded fourth bytes")
}

func TestEncodeDecodeString(t *testing.T) {
	// Declare strings to use to encode/decode
	firstString := "BAND"
	secondString := "BTC"
	thirdString := "ETH"
	fourthString := ""
	fifthString := "!@#$%^&*(ERTEYRUTIKJ!$#@%^&*&("

	encoder := NewEncoder()

	// Test strings returned from encode functions
	encoder.EncodeString(firstString)
	encoder.EncodeString(secondString)
	encoder.EncodeString(thirdString)
	encoder.EncodeString(fourthString)
	encoder.EncodeString(fifthString)

	decoder := NewDecoder(encoder.data)

	// Test strings returned when passing encoded data into decoder
	decodedString, err := decoder.DecodeString()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, firstString, decodedString, "Incorrect decoded first string")

	decodedString, err = decoder.DecodeString()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, secondString, decodedString, "Incorrect decoded second string")

	decodedString, err = decoder.DecodeString()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, thirdString, decodedString, "Incorrect decoded third string")

	decodedString, err = decoder.DecodeString()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, fourthString, decodedString, "Incorrect decoded fourth string")

	decodedString, err = decoder.DecodeString()
	require.Nil(t, err, "Incorrect decoded error type")
	require.Equal(t, fifthString, decodedString, "Incorrect decoded fifth string")
}
