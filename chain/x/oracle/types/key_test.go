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
	oracleScriptID := OracleScriptID(3)
	calldata := []byte("calldata")
	expectKeyByte, _ := hex.DecodeString("ff0000000000000014000000000000000363616c6c64617461")
	require.Equal(t, expectKeyByte, ResultStoreKey(requestID, oracleScriptID, calldata))
}

func TestRawDataRequestStoreKey(t *testing.T) {
	requestID := RequestID(20)
	externalID := ExternalID(947313)
	expectKeyByte, _ := hex.DecodeString("02000000000000001400000000000E7471")
	require.Equal(t, expectKeyByte, RawDataRequestStoreKey(requestID, externalID))
}

func TestRawDataReportStoreKey(t *testing.T) {
	requestID := RequestID(20)
	externalID := ExternalID(6)
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	expectKeyByte, _ := hex.DecodeString("0300000000000000140000000000000006b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, expectKeyByte, RawDataReportStoreKey(requestID, externalID, validator))
}

func TestDataSourceStoreKey(t *testing.T) {
	dataSourceID := DataSourceID(888)
	expectKeyByte, _ := hex.DecodeString("040000000000000378")

	require.Equal(t, expectKeyByte, DataSourceStoreKey(dataSourceID))
}

func TestOracleScriptStoreKey(t *testing.T) {
	oracleScriptID := OracleScriptID(123)
	expectKeyByte, _ := hex.DecodeString("05000000000000007b")

	require.Equal(t, expectKeyByte, OracleScriptStoreKey(oracleScriptID))
}

func TestGetExternalIDFromRawDataRequestKey(t *testing.T) {
	key, _ := hex.DecodeString("02000000000000001400000000000E7471")
	externalID := GetExternalIDFromRawDataRequestKey(key)

	require.Equal(t, ExternalID(947313), externalID)
}

func TestGetValidatorAddressAndExternalID(t *testing.T) {
	key, _ := hex.DecodeString("0300000000000000140000000000000006b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	valAddress, externalID := GetValidatorAddressAndExternalID(key, 1)

	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	require.Equal(t, valAddress, validator)
	require.Equal(t, ExternalID(6), externalID)
}
