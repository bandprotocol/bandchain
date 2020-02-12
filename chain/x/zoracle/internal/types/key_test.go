package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRequestStoreKey(t *testing.T) {
	requestID := int64(20)
	expectKeyByte, _ := hex.DecodeString("010000000000000014")

	require.Equal(t, expectKeyByte, RequestStoreKey(requestID))
}

func TestResultStoreKey(t *testing.T) {
	requestID := int64(20)
	codeHash := []byte("code")
	params := []byte("params")

	expectKeyByte, _ := hex.DecodeString("ff0000000000000014636f6465706172616d73")
	require.Equal(t, expectKeyByte, ResultStoreKey(requestID, codeHash, params))
}

func TestCodeHashStoreKey(t *testing.T) {
	codeHash, _ := hex.DecodeString("d765967790f8d531be31d395de5cbce145e71f4d76ed90eda01646d2638ad932")
	expectKeyByte, _ := hex.DecodeString("02d765967790f8d531be31d395de5cbce145e71f4d76ed90eda01646d2638ad932")

	require.Equal(t, expectKeyByte, CodeHashStoreKey(codeHash))
}

func TestReportStoreKey(t *testing.T) {
	requestID := int64(20)
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	expectKeyByte, _ := hex.DecodeString("030000000000000014b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, expectKeyByte, ReportStoreKey(requestID, validator))
}

func TestDataSourceStoreKey(t *testing.T) {
	dataSourceID := int64(888)
	expectKeyByte, _ := hex.DecodeString("040000000000000378")

	require.Equal(t, expectKeyByte, DataSourceStoreKey(dataSourceID))
}
