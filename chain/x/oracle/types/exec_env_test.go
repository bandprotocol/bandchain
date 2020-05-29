package types_test

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

func mockExecutionEnv() *types.ExecEnv {
	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")
	rawRequestID := []types.ExternalID{1, 2, 3}
	request := types.NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := types.NewExecEnv(request, int64(1581589770), 0)
	rawReport1 := types.NewRawReport(1, 0, []byte("DATA1"))
	rawReport2 := types.NewRawReport(2, 1, []byte("DATA2"))
	rawReport3 := types.NewRawReport(3, 0, []byte("DATA3"))

	report1 := types.NewReport(validatorAddress1, []types.RawReport{rawReport1, rawReport2})
	report2 := types.NewReport(validatorAddress2, []types.RawReport{rawReport3})

	env.SetReports([]types.Report{report1, report2})
	return env
}

func mockPreparedExecEnv() *types.ExecEnv {
	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")
	rawRequestID := []types.ExternalID{1, 2, 3}
	request := types.NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := types.NewExecEnv(request, int64(1581589770), 3)
	env.RequestExternalData(1, 0, []byte("CALLDATA1"))
	env.RequestExternalData(2, 1, []byte("CALLDATA2"))
	env.RequestExternalData(3, 0, []byte("CALLDATA3"))
	return env
}

func mockFreshEnv() *types.ExecEnv {
	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := int64(1581589700)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")
	rawRequestID := []types.ExternalID{1, 2, 3}
	request := types.NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, &ibcInfo, rawRequestID)
	env := types.NewExecEnv(request, int64(1581589770), 0)
	return env
}
func TestGetMaxRawRequestDataSize(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(types.MaxCalldataSize), env.GetMaxRawRequestDataSize())
}

func TestGetMaxResultSize(t *testing.T) {
	env := mockExecutionEnv()
	require.Equal(t, int64(types.MaxResultSize), env.GetMaxResultSize())
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

	rawReport1 := types.NewRawReport(1, 0, []byte("DATA1"))
	rawReport2 := types.NewRawReport(2, 1, []byte("DATA2"))
	rawReport3 := types.NewRawReport(3, 0, []byte("DATA3"))

	report1 := types.NewReport(validatorAddress1, []types.RawReport{rawReport1, rawReport2})
	report2 := types.NewReport(validatorAddress2, []types.RawReport{rawReport3})

	require.Panics(t, func() { env.SetReports([]types.Report{report1, report2}) })
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
	expect := []types.RawRequest{
		types.NewRawRequest(0, 1, []byte("CALLDATA1")),
		types.NewRawRequest(1, 2, []byte("CALLDATA2")),
		types.NewRawRequest(0, 3, []byte("CALLDATA3")),
	}
	require.Equal(t, expect, env.GetRawRequests())

	env = mockExecutionEnv()
	require.Panics(t, func() {
		env.GetRawRequests()
	})
}
