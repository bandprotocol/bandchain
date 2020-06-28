package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRequestStoreKey(t *testing.T) {
	requestID := RequestID(20)
	expectKeyByte, _ := hex.DecodeString("010000000000000014")

	require.Equal(t, expectKeyByte, RequestStoreKey(requestID))
}

func TestResultStoreKey(t *testing.T) {
	requestID := RequestID(20)
	expectKeyByte, _ := hex.DecodeString("ff0000000000000014")
	require.Equal(t, expectKeyByte, ResultStoreKey(requestID))
}

func TestReportStoreKeyPerValidator(t *testing.T) {
	requestID := RequestID(20)
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	expectKeyByte, _ := hex.DecodeString("020000000000000014b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, expectKeyByte, ReportStoreKeyPerValidator(requestID, validator))
}

func TestDataSourceStoreKey(t *testing.T) {
	dataSourceID := DataSourceID(888)
	expectKeyByte, _ := hex.DecodeString("030000000000000378")

	require.Equal(t, expectKeyByte, DataSourceStoreKey(dataSourceID))
}

func TestOracleScriptStoreKey(t *testing.T) {
	oracleScriptID := OracleScriptID(123)
	expectKeyByte, _ := hex.DecodeString("04000000000000007b")

	require.Equal(t, expectKeyByte, OracleScriptStoreKey(oracleScriptID))
}

func TestReportInfoStoreKey(t *testing.T) {
	v, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	expectKeyByte, _ := hex.DecodeString("06b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	require.Equal(t, expectKeyByte, ReportInfoStoreKey(v))
}

func TestReportStoreKey(t *testing.T) {
	requestID := RequestID(12)
	expectKeyByte, _ := hex.DecodeString("02000000000000000c")
	require.Equal(t, expectKeyByte, ReportStoreKey(requestID))
}

func TestReporterStoreKey(t *testing.T) {
	validatorAddress, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	reporterAddress, _ := sdk.AccAddressFromHex("ba11d00c5f74255f56a5e366f4f77f5a186d7f55")
	expectKeyByte, _ := hex.DecodeString("05b80f2a5df7d5710b15622d1a9f1e3830ded5bda8ba11d00c5f74255f56a5e366f4f77f5a186d7f55")

	require.Equal(t, expectKeyByte, ReporterStoreKey(validatorAddress, reporterAddress))
}

func TestValidatorReporterPrefixKey(t *testing.T) {
	validatorAddress, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	expectKeyByte, _ := hex.DecodeString("05b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, expectKeyByte, ValidatorReporterPrefixKey(validatorAddress))
}
