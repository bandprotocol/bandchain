package obi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Inner struct {
	A uint8
	B uint8
}

type ExampleData struct {
	Symbol string
	Px     uint64
	In     Inner
	Arr    []int16
}

func TestEncodeBytes(t *testing.T) {
	require.Nil(t, MustEncode(ExampleData{
		Symbol: "BTC",
		Px:     9000,
		In: Inner{
			A: 1,
			B: 2,
		},
		Arr: []int16{10, 11},
	}))
}

func TestDecode1(t *testing.T) {
	var n ExampleData
	err := Decode([]byte{0x3, 0x0, 0x0, 0x0, 0x42, 0x54, 0x43, 0x28, 0x23, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x2, 0x0, 0x0, 0x0, 0xa, 0x0, 0xb, 0x0}, &n)
	require.Nil(t, err)
	require.Nil(t, n)
}
