package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/GeoDB-Limited/odincore/go-owasm/api"
)

var (
	pk1               = secp256k1.GenPrivKey().PubKey()
	pk2               = secp256k1.GenPrivKey().PubKey()
	pk3               = secp256k1.GenPrivKey().PubKey()
	addr1             = pk1.Address()
	addr2             = pk2.Address()
	addr3             = pk3.Address()
	validatorAddress1 = sdk.ValAddress(addr1)
	validatorAddress2 = sdk.ValAddress(addr2)
	validatorAddress3 = sdk.ValAddress(addr3)
)

func mockExecEnv() *ExecuteEnv {
	oracleScriptID := OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := time.Unix(1581589700, 0)
	clientID := "beeb"
	request := NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, nil)
	rawReport1 := NewRawReport(1, 0, []byte("DATA1"))
	rawReport2 := NewRawReport(2, 1, []byte("DATA2"))
	rawReport3 := NewRawReport(3, 0, []byte("DATA3"))
	report1 := NewReport(validatorAddress1, true, []RawReport{rawReport1, rawReport2})
	report2 := NewReport(validatorAddress2, true, []RawReport{rawReport3})
	env := NewExecuteEnv(request, []Report{report1, report2})
	return env
}

func mockFreshPrepareEnv() *PrepareEnv {
	oracleScriptID := OracleScriptID(1)
	calldata := []byte("CALLDATA")
	valAddresses := []sdk.ValAddress{validatorAddress1, validatorAddress2, validatorAddress3}
	minCount := uint64(1)
	requestHeight := int64(999)
	requestTime := time.Unix(1581589700, 0)
	clientID := "beeb"
	request := NewRequest(oracleScriptID, calldata, valAddresses, minCount, requestHeight, requestTime, clientID, nil)
	env := NewPrepareEnv(request, 3)
	return env
}

func mockAlreadyPreparedEnv() *PrepareEnv {
	env := mockFreshPrepareEnv()
	env.AskExternalData(1, 1, []byte("CALLDATA1"))
	env.AskExternalData(2, 2, []byte("CALLDATA2"))
	env.AskExternalData(3, 3, []byte("CALLDATA3"))
	return env
}

func TestGetCalldata(t *testing.T) {
	calldata := []byte("CALLDATA")

	penv := mockFreshPrepareEnv()
	require.Equal(t, calldata, penv.GetCalldata())

	eenv := mockExecEnv()
	require.Equal(t, calldata, eenv.GetCalldata())
}

func TestSetReturnData(t *testing.T) {
	result := []byte("RESULT")

	penv := mockFreshPrepareEnv()
	err := penv.SetReturnData(result)
	require.Equal(t, api.ErrWrongPeriodAction, err)

	eenv := mockExecEnv()
	eenv.SetReturnData(result)
	require.Equal(t, result, eenv.Retdata)

}
func TestGetAskCount(t *testing.T) {
	// Can call on both environment
	penv := mockFreshPrepareEnv()
	require.Equal(t, int64(3), penv.GetAskCount())

	eenv := mockExecEnv()
	require.Equal(t, int64(3), eenv.GetAskCount())
}

func TestGetMinCount(t *testing.T) {
	// Can call on both environment
	penv := mockFreshPrepareEnv()
	require.Equal(t, int64(1), penv.GetMinCount())

	eenv := mockExecEnv()
	require.Equal(t, int64(1), eenv.GetMinCount())
}

func TestGetAnsCount(t *testing.T) {
	// Should return error if call on prepare environment.
	penv := mockFreshPrepareEnv()
	_, err := penv.GetAnsCount()
	require.Equal(t, api.ErrWrongPeriodAction, err)

	eenv := mockExecEnv()
	v, err := eenv.GetAnsCount()
	require.NoError(t, err)
	require.Equal(t, int64(2), v)
}

func TestGetExternalData(t *testing.T) {
	env := mockExecEnv()

	data, err := env.GetExternalData(1, 0)
	require.NoError(t, err)
	require.Equal(t, []byte("DATA1"), data)
	status, err := env.GetExternalDataStatus(1, 0)
	require.NoError(t, err)
	require.Equal(t, int64(0), status)

	data, err = env.GetExternalData(0, 2)
	require.NoError(t, err)
	require.Nil(t, data)
	status, err = env.GetExternalDataStatus(0, 2)
	require.NoError(t, err)
	require.Equal(t, int64(-1), status)

	_, err = env.GetExternalData(1, 100)
	require.Equal(t, api.ErrBadValidatorIndex, err)
	_, err = env.GetExternalDataStatus(1, 100)
	require.Equal(t, api.ErrBadValidatorIndex, err)

	_, err = env.GetExternalData(1, -1)
	require.Error(t, err)
	_, err = env.GetExternalDataStatus(1, -1)
	require.Error(t, err)

	_, err = env.GetExternalData(100, 0)
	require.Error(t, err)
	_, err = env.GetExternalDataStatus(100, 0)
	require.Error(t, err)
}

func TestFailedGetExternalData(t *testing.T) {
	penv := mockAlreadyPreparedEnv()

	_, err := penv.GetExternalData(1, 1)
	require.Equal(t, api.ErrWrongPeriodAction, err)
	_, err = penv.GetExternalDataStatus(1, 1)
	require.Equal(t, api.ErrWrongPeriodAction, err)
}

func TestAskExternalData(t *testing.T) {
	env := mockFreshPrepareEnv()
	env.AskExternalData(1, 1, []byte("CALLDATA1"))
	env.AskExternalData(42, 2, []byte("CALLDATA2"))
	env.AskExternalData(3, 4, []byte("CALLDATA3"))

	rawReq := env.GetRawRequests()
	expectRawReq := []RawRequest{
		NewRawRequest(1, 1, []byte("CALLDATA1")),
		NewRawRequest(42, 2, []byte("CALLDATA2")),
		NewRawRequest(3, 4, []byte("CALLDATA3")),
	}
	require.Equal(t, expectRawReq, rawReq)
}

func TestAskExternalDataOnTooSmallSpan(t *testing.T) {
	penv := mockFreshPrepareEnv()

	err := penv.AskExternalData(1, 3, make([]byte, MaxDataSize+1))
	require.Equal(t, api.ErrSpanTooSmall, err)
	require.Equal(t, []RawRequest(nil), penv.GetRawRequests())
}
func TestAskTooManyExternalData(t *testing.T) {
	penv := mockFreshPrepareEnv()

	err := penv.AskExternalData(1, 1, []byte("CALLDATA1"))
	require.NoError(t, err)
	err = penv.AskExternalData(2, 2, []byte("CALLDATA2"))
	require.NoError(t, err)
	err = penv.AskExternalData(3, 3, []byte("CALLDATA3"))
	require.NoError(t, err)

	expectRawReq := []RawRequest{
		NewRawRequest(1, 1, []byte("CALLDATA1")),
		NewRawRequest(2, 2, []byte("CALLDATA2")),
		NewRawRequest(3, 3, []byte("CALLDATA3")),
	}
	require.Equal(t, expectRawReq, penv.GetRawRequests())

	err = penv.AskExternalData(4, 4, []byte("CALLDATA4"))
	require.Equal(t, api.ErrTooManyExternalData, err)
	require.Equal(t, expectRawReq, penv.GetRawRequests())
}

func TestAskDuplicateExternalID(t *testing.T) {
	penv := mockFreshPrepareEnv()

	err := penv.AskExternalData(1, 1, []byte("CALLDATA1"))
	require.NoError(t, err)
	err = penv.AskExternalData(2, 1, []byte("CALLDATA2"))
	require.NoError(t, err)

	err = penv.AskExternalData(1, 3, []byte("CALLDATA3"))
	require.Equal(t, api.ErrDuplicateExternalID, err)
}

func TestAskExternalDataOnExecEnv(t *testing.T) {
	env := mockExecEnv()
	calldata := []byte("CALLDATA")
	err := env.AskExternalData(2, 2, calldata)
	require.Equal(t, api.ErrWrongPeriodAction, err)
}

func TestGetRawRequests(t *testing.T) {
	env := mockAlreadyPreparedEnv()
	expect := []RawRequest{
		NewRawRequest(1, 1, []byte("CALLDATA1")),
		NewRawRequest(2, 2, []byte("CALLDATA2")),
		NewRawRequest(3, 3, []byte("CALLDATA3")),
	}
	require.Equal(t, expect, env.GetRawRequests())
}
