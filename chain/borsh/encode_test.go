package goborsh

import (
	"bytes"
	"testing"
)

func TestEncodeU8(t *testing.T) {
	encoder := NewBorshEncoder()
	expectedReturnByte := []byte{5}
	expectedEncoderData := []byte{5}
	returnBytes := encoder.EncodeU8(5)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}

	expectedReturnByte = []byte{6}
	expectedEncoderData = []byte{5, 6}
	returnBytes = encoder.EncodeU8(6)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}

	expectedReturnByte = []byte{7}
	expectedEncoderData = []byte{5, 6, 7}
	returnBytes = encoder.EncodeU8(7)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}
}

func TestEncodeU16(t *testing.T) {
	encoder := NewBorshEncoder()
	expectedReturnByte := []byte{5, 0}
	expectedEncoderData := []byte{5}
	returnBytes := encoder.EncodeU16(5)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}

	expectedReturnByte = []byte{6}
	expectedEncoderData = []byte{5, 0, 6, 0}
	returnBytes = encoder.EncodeU16(6)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}

	expectedReturnByte = []byte{132, 3}
	expectedEncoderData = []byte{5, 0, 6, 0, 132, 3}
	returnBytes = encoder.EncodeU16(900)
	if bytes.Compare(expectedReturnByte, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedReturnByte)
	}
	if bytes.Compare(expectedEncoderData, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedEncoderData)
	}
}

func TestEncodeU32Single(t *testing.T) {
	expectedBytes := []byte{57, 48, 128, 0}
	encoder := NewBorshEncoder()
	returnBytes := encoder.EncodeU32(8400953)
	if bytes.Compare(expectedBytes, returnBytes) == 1 {
		t.Errorf("Incorrect return bytes: got %v expected %v", returnBytes, expectedBytes)
	}
	if bytes.Compare(expectedBytes, encoder.data) == 1 {
		t.Errorf("Incorrect encoder data bytes: got %v expected %v", encoder.data, expectedBytes)
	}
}
