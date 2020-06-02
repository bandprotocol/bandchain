package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	pk1               = ed25519.GenPrivKey().PubKey()
	pk2               = ed25519.GenPrivKey().PubKey()
	pk3               = ed25519.GenPrivKey().PubKey()
	addr1             = pk1.Address()
	addr2             = pk2.Address()
	addr3             = pk3.Address()
	validatorAddress1 = sdk.ValAddress(addr1)
	validatorAddress2 = sdk.ValAddress(addr2)
	validatorAddress3 = sdk.ValAddress(addr3)
)

func mockExecutionEnv() *ExecEnv {
	oracleScriptID := OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := NewIBCInfo("source_port", "source_channel")
	rawRequestID := []ExternalID{1, 2, 3}
	request := NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := NewExecEnv(request, int64(1581589770), 0)
	rawReport1 := NewRawReport(1, 0, []byte("DATA1"))
	rawReport2 := NewRawReport(2, 1, []byte("DATA2"))
	rawReport3 := NewRawReport(3, 0, []byte("DATA3"))

	report1 := NewReport(validatorAddress1, []RawReport{rawReport1, rawReport2})
	report2 := NewReport(validatorAddress2, []RawReport{rawReport3})

	env.SetReports([]Report{report1, report2})
	return env
}

func mockPreparedExecEnv() *ExecEnv {
	oracleScriptID := OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := NewIBCInfo("source_port", "source_channel")
	rawRequestID := []ExternalID{1, 2, 3}
	request := NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := NewExecEnv(request, int64(1581589770), 3)
	env.RequestExternalData(1, 0, []byte("CALLDATA1"))
	env.RequestExternalData(2, 1, []byte("CALLDATA2"))
	env.RequestExternalData(3, 0, []byte("CALLDATA3"))
	return env
}

func mockFreshEnv() *ExecEnv {
	oracleScriptID := OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := NewIBCInfo("source_port", "source_channel")
	rawRequestID := []ExternalID{1, 2, 3}
	request := NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := NewExecEnv(request, int64(1581589770), 3)
	return env
}
func TestGetMaxRawRequestDataSize(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(MaxCalldataSize), env.GetMaxRawRequestDataSize())
}

func TestGetMaxResultSize(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(MaxResultSize), env.GetMaxResultSize())
}

func TestGetAskCount(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(3), env.GetAskCount())
}

func TestGetMinCount(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(1), env.GetMinCount())
}

func TestGetAnsCount(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(2), env.GetAnsCount())
}

func TestGetPrepareBlockTime(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(1581589700), env.GetPrepareBlockTime())
}

func TestGetAggregateBlockTime(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(1581589770), env.GetAggregateBlockTime())

	env = mockPreparedExecEnv()
	require.Equal(t, int64(0), env.GetAggregateBlockTime())

}

func TestGetValidatorAddress(t *testing.T) {
	env := mockExecutionEnv()

	valAddrs, err := env.GetValidatorAddress(0)
	require.NoError(t, err)
	require.Equal(t, []byte(validatorAddress1), valAddrs)

	valAddrs, err = env.GetValidatorAddress(1)
	require.NoError(t, err)
	require.Equal(t, []byte(validatorAddress2), valAddrs)

	_, err = env.GetValidatorAddress(100)
	require.Error(t, err)

	_, err = env.GetValidatorAddress(-1)
	require.Error(t, err)
}

func TestSetReportFail(t *testing.T) {
	env := mockPreparedExecEnv()

	rawReport1 := NewRawReport(1, 0, []byte("DATA1"))
	rawReport2 := NewRawReport(2, 1, []byte("DATA2"))
	rawReport3 := NewRawReport(3, 0, []byte("DATA3"))

	report1 := NewReport(validatorAddress1, []RawReport{rawReport1, rawReport2})
	report2 := NewReport(validatorAddress2, []RawReport{rawReport3})

	require.Panics(t, func() { env.SetReports([]Report{report1, report2}) })
}

func TestGetExternalData(t *testing.T) {
	env := mockExecutionEnv()
	data, exitCode, err := env.GetExternalData(1, 0)
	require.NoError(t, err)
	require.Equal(t, uint32(0), exitCode)
	require.Equal(t, []byte("DATA1"), data)

	data, exitCode, err = env.GetExternalData(2, 0)
	require.NoError(t, err)
	require.Equal(t, uint32(1), exitCode)
	require.Equal(t, []byte("DATA2"), data)

	_, _, err = env.GetExternalData(3, 0)
	require.Error(t, err)

	_, _, err = env.GetExternalData(1, 1)
	require.Error(t, err)

	_, _, err = env.GetExternalData(1, 2)
	require.Error(t, err)

	data, exitCode, err = env.GetExternalData(3, 1)
	require.NoError(t, err)
	require.Equal(t, uint32(0), exitCode)
	require.Equal(t, []byte("DATA3"), data)

	_, _, err = env.GetExternalData(1, 100)
	require.Error(t, err)

	_, _, err = env.GetExternalData(1, 3)
	require.Error(t, err)
}

func TestRequestExternalData(t *testing.T) {
	env := mockFreshEnv()
	err := env.RequestExternalData(1, 0, []byte("CALLDATA1"))
	require.NoError(t, err)
	err = env.RequestExternalData(2, 1, []byte("CALLDATA2"))
	require.NoError(t, err)
	err = env.RequestExternalData(3, 0, []byte("CALLDATA3"))
	require.NoError(t, err)

	rawReq := env.GetRawRequests()
	expectRawReq := []RawRequest{
		NewRawRequest(0, 1, []byte("CALLDATA1")),
		NewRawRequest(1, 2, []byte("CALLDATA2")),
		NewRawRequest(0, 3, []byte("CALLDATA3")),
	}
	require.Equal(t, expectRawReq, rawReq)
}

func TestRequestExternalDataFail(t *testing.T) {
	env := mockExecutionEnv()
	calldata := []byte("CALLDATA")

	err := env.RequestExternalData(2, 2, calldata)
	require.Error(t, err)

	calldata = make([]byte, 2000)
	err = env.RequestExternalData(1, 1, calldata)
	require.Error(t, err)
}

func TestGetRawRequests(t *testing.T) {
	env := mockPreparedExecEnv()
	expect := []RawRequest{
		NewRawRequest(0, 1, []byte("CALLDATA1")),
		NewRawRequest(1, 2, []byte("CALLDATA2")),
		NewRawRequest(0, 3, []byte("CALLDATA3")),
	}
	require.Equal(t, expect, env.GetRawRequests())

	env = mockExecutionEnv()
	require.Panics(t, func() {
		env.GetRawRequests()
	})
}
