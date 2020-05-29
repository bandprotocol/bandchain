package obi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ID uint8

type Inner struct {
	A ID    `obi:"a"`
	B uint8 `obi:"b"`
}

type ExampleData struct {
	Symbol string  `obi:"symbol"`
	Px     uint64  `obi:"px"`
	In     Inner   `obi:"in"`
	Arr    []int16 `obi:"arr"`
}

func TestEncodeBytes(t *testing.T) {
	require.Equal(t, MustEncode(ExampleData{
		Symbol: "BTC",
		Px:     9000,
		In: Inner{
			A: 1,
			B: 2,
		},
		Arr: []int16{10, 11},
	}), []byte{0x0, 0x0, 0x0, 0x3, 0x42, 0x54, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x23, 0x28, 0x1, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xa, 0x0, 0xb})
}

func TestDecode1(t *testing.T) {
	var n ExampleData
	err := Decode([]byte{0x0, 0x0, 0x0, 0x3, 0x42, 0x54, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x23, 0x28, 0x1, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xa, 0x0, 0xb}, &n)
	require.Nil(t, err)
	require.Equal(t, n, ExampleData{
		Symbol: "BTC",
		Px:     9000,
		In: Inner{
			A: 1,
			B: 2,
		},
		Arr: []int16{10, 11},
	})
}

func TestSchema(t *testing.T) {
	require.Equal(t, MustGetSchema(ExampleData{}), "{symbol:string,px:u64,in:{a:u8,b:u8},arr:[i16]}")
}
