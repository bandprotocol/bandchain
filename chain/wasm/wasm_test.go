package wasm

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func loadWasmFile() ([]byte, wasm.Instance) {
	file, err := os.Open("./res/result.wasm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		panic(statsErr)
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	if err != nil {
		panic(err)
	}
	instance, err := wasm.NewInstance(bytes)
	if err != nil {
		panic(err)
	}
	return bytes, instance
}

func TestParamsInfo(t *testing.T) {
	code, _ := loadWasmFile()
	info, err := ParamsInfo(code)
	require.Nil(t, err)
	expect := `[["symbol","coins::Coins"]]`
	require.Equal(t, expect, string(info))
}

func TestParseParams(t *testing.T) {
	code, _ := loadWasmFile()
	params, _ := hex.DecodeString("00000001")
	paramsByte, err := ParseParams(code, params)
	require.Nil(t, err)
	expect := `{"symbol":"ETH"}`
	require.Equal(t, expect, string(paramsByte))
}

func TestRawDataInfo(t *testing.T) {
	code, _ := loadWasmFile()
	info, err := RawDataInfo(code)
	require.Nil(t, err)
	expect := `[["coin_gecko","f32"],["crypto_compare","f32"]]`
	require.Equal(t, expect, string(info))
}

func TestParseRawData(t *testing.T) {
	code, _ := loadWasmFile()
	params, _ := hex.DecodeString("00000000")
	data, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")
	dataByte, err := ParseRawData(code, params, data)
	require.Nil(t, err)
	expect := `{"coin_gecko":7194.25,"crypto_compare":7214.12}`
	require.Equal(t, expect, string(dataByte))
}

func TestResultInfo(t *testing.T) {
	code, _ := loadWasmFile()
	info, err := ResultInfo(code)
	require.Nil(t, err)
	expect := `[["price_in_usd","u64"]]`
	require.Equal(t, expect, string(info))
}

func TestParseResult(t *testing.T) {
	code, _ := loadWasmFile()
	result, _ := hex.DecodeString("00000000000d0e72")
	resultByte, err := ParseResult(code, result)
	require.Nil(t, err)
	expect := `{"price_in_usd":855666}`
	require.Equal(t, expect, string(resultByte))
}

func TestParseEmptyResult(t *testing.T) {
	code, _ := loadWasmFile()
	result := []byte(nil)
	_, err := ParseResult(code, result)
	require.EqualError(t, err, "Failed to call the `__parse_result` exported function.")
}

func TestAllocateInner(t *testing.T) {
	_, instance := loadWasmFile()
	// Small data
	ptr, err := allocateInner(instance, []byte("test"))
	require.Nil(t, err)

	result, err := parseOutput(instance, ptr)
	require.Nil(t, err)
	require.Equal(t, "test", string(result))

	// Big data
	bigBytes := make([]byte, 0)
	for i := 0; i < 1000; i++ {
		bigBytes = append(bigBytes, byte(0xa))
	}

	ptr, err = allocateInner(instance, bigBytes)
	require.Nil(t, err)

	result, err = parseOutput(instance, ptr)
	require.Nil(t, err)
	require.Equal(t, bigBytes, result)
}

func TestAllocate(t *testing.T) {
	_, instance := loadWasmFile()
	data := [][]byte{[]byte("test1"), []byte("test2"), []byte("test3"), []byte("test4")}
	ptr, err := allocate(instance, data)
	require.Nil(t, err)

	// Size must be 4
	sz, pointer := int(ptr>>32), (ptr & ((1 << 32) - 1))
	require.Equal(t, 4, sz)
	mem := instance.Memory.Data()[pointer:]
	for idx := 0; idx < sz; idx++ {
		pt := int64(binary.LittleEndian.Uint64(mem[8*idx : 8*idx+8]))
		out, err := parseOutput(instance, pt)
		require.Nil(t, err)
		require.Equal(t, data[idx], out)
	}
}

func TestPrepare(t *testing.T) {
	code, _ := loadWasmFile()
	params, _ := hex.DecodeString("00000001")
	prepare, err := Prepare(code, params)
	require.Nil(t, err)
	expect := `[{"cmd":"curl","args":["https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"]},{"cmd":"curl","args":["https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD"]}]`
	require.Equal(t, expect, string(prepare))
}

func TestExecute(t *testing.T) {
	code, _ := loadWasmFile()
	params, _ := hex.DecodeString("00000000")
	data, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")
	expect, _ := hex.DecodeString("00000000000afe22")
	result, err := Execute(code, params, [][]byte{data, data})
	require.Nil(t, err)
	require.Equal(t, expect, result)
}
