package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDataPointStoreKey(t *testing.T) {
	requestID := uint64(20)
	expectKeyByte, _ := hex.DecodeString("010000000000000014")

	require.Equal(t, expectKeyByte, DataPointStoreKey(requestID))
}

func TestCodeHashStoreKey(t *testing.T) {
	codeHash, _ := hex.DecodeString("d765967790f8d531be31d395de5cbce145e71f4d76ed90eda01646d2638ad932")
	expectKeyByte, _ := hex.DecodeString("02d765967790f8d531be31d395de5cbce145e71f4d76ed90eda01646d2638ad932")

	require.Equal(t, expectKeyByte, CodeHashStoreKey(codeHash))
}

func TestReportStoreKey(t *testing.T) {
	requestID := uint64(20)
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	expectKeyByte, _ := hex.DecodeString("030000000000000014b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, expectKeyByte, ReportStoreKey(requestID, validator))
}
