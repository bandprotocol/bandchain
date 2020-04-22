package borsh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeU8(t *testing.T) {
	encoder := NewEncoder()

	expectedEncoderData := []byte{5}
	encoder.EncodeU8(5)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")

	expectedEncoderData = []byte{5, 6}
	encoder.EncodeU8(6)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")

	expectedEncoderData = []byte{5, 6, 7}
	encoder.EncodeU8(7)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")
}

func TestEncodeU16(t *testing.T) {
	encoder := NewEncoder()

	expectedEncoderData := []byte{5, 0}
	encoder.EncodeU16(5)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")

	expectedEncoderData = []byte{5, 0, 6, 0}
	encoder.EncodeU16(6)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")

	expectedEncoderData = []byte{5, 0, 6, 0, 132, 3}
	encoder.EncodeU16(900)
	require.Equal(t, expectedEncoderData, encoder.GetEncodedData(), "Incorrect encoder data bytes")
}

func TestEncodeU32Single(t *testing.T) {
	expectedBytes := []byte{57, 48, 128, 0}
	encoder := NewEncoder()
	encoder.EncodeU32(8400953)
	require.Equal(t, expectedBytes, encoder.GetEncodedData(), "Incorrect encoder data bytes")
}
