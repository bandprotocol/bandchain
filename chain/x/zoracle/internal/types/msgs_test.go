package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRequestData(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 100, sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "request", msg.Type())
	require.Equal(t, int64(1), msg.OracleScriptID)
	require.Equal(t, []byte("calldata"), msg.Calldata)
	require.Equal(t, int64(10), msg.RequestedValidatorCount)
	require.Equal(t, int64(5), msg.SufficientValidatorCount)
	require.Equal(t, int64(100), msg.Expiration)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgRequestDataValidation(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	requestedValidatorCount := int64(10)
	sufficientValidatorCount := int64(5)
	expiration := int64(100)
	cases := []struct {
		valid bool
		tx    MsgRequestData
	}{
		{
			true, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, sender,
			),
		},
		{
			false, NewMsgRequestData(
				0, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, sender,
			),
		},
		{
			true, NewMsgRequestData(
				1, nil, requestedValidatorCount,
				sufficientValidatorCount, expiration, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 0,
				sufficientValidatorCount, expiration, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				-1, expiration, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 6,
				8, expiration, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, -10, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, nil,
			),
		},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgRequestDataGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 100, sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Request","value":{"calldata":"Y2FsbGRhdGE=","expiration":"100","oracleScriptID":"1","requestedValidatorCount":"10","sender":"band1wdjkuer9wgvz7c4y","sufficientValidatorCount":"5"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgReportData(t *testing.T) {
	requestID := int64(3)
	data := []ExternalData{NewExternalData(1, []byte("data1")), NewExternalData(2, []byte("data2"))}
	provider, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReportData(requestID, data, provider)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "report", msg.Type())
}

func TestMsgReportDataValidation(t *testing.T) {
	requestID := int64(3)
	data := []ExternalData{NewExternalData(1, []byte("data1")), NewExternalData(2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failValidator, _ := sdk.ValAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgReportData
	}{
		{true, NewMsgReportData(requestID, data, validator)},
		{false, NewMsgReportData(-1, data, validator)},
		{false, NewMsgReportData(requestID, []ExternalData{}, validator)},
		{false, NewMsgReportData(requestID, nil, validator)},
		{false, NewMsgReportData(requestID, data, failValidator)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgReportDataGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForValidator("band"+sdk.PrefixValidator+sdk.PrefixOperator, "band"+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)

	requestID := int64(3)
	data := []ExternalData{NewExternalData(1, []byte("data1")), NewExternalData(2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReportData(requestID, data, validator)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Report","value":{"dataSet":[{"data":"ZGF0YTE=","externalDataID":"1"},{"data":"ZGF0YTI=","externalDataID":"2"}],"requestID":"3","sender":"bandvaloper1hq8j5h0h64csk9tz95df783cxr0dt0dgay2kyy"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgStoreCode(t *testing.T) {
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	name := "script name"
	msg := NewMsgStoreCode([]byte("Code"), name, owner)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "store", msg.Type())
}

func TestMsgStoreCodeValidation(t *testing.T) {
	code := []byte("Code")
	name := "script name"
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failOwner, _ := sdk.AccAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgStoreCode
	}{
		{true, NewMsgStoreCode(code, name, owner)},
		{false, NewMsgStoreCode([]byte(""), name, owner)},
		{false, NewMsgStoreCode(code, "", owner)},
		{false, NewMsgStoreCode(nil, name, owner)},
		{false, NewMsgStoreCode(code, name, failOwner)},
		{false, NewMsgStoreCode(code, name, nil)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgStoreCodeGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	name := "script name"
	code := []byte("Code")
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgStoreCode(code, name, owner)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Store","value":{"code":"Q29kZQ==","name":"script name","owner":"band1hq8j5h0h64csk9tz95df783cxr0dt0dg3jw4p0"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgDeleteCode(t *testing.T) {
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	msg := NewMsgDeleteCode(codeHash, owner)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "delete", msg.Type())
}

func TestMsgDeleteCodeValidation(t *testing.T) {
	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failOwner, _ := sdk.AccAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgDeleteCode
	}{
		{true, NewMsgDeleteCode(codeHash, owner)},
		{false, NewMsgDeleteCode([]byte(""), owner)},
		{false, NewMsgDeleteCode(nil, owner)},
		{false, NewMsgDeleteCode([]byte("code"), failOwner)},
		{false, NewMsgDeleteCode(codeHash, failOwner)},
		{false, NewMsgDeleteCode(codeHash, nil)},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgDeleteCodeGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")

	msg := NewMsgDeleteCode(codeHash, owner)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Delete","value":{"codeHash":"VpTQii5T/8rgwxA+Wtb2B2q9lg6x+KVldwQLwQKPcCs=","owner":"band1hq8j5h0h64csk9tz95df783cxr0dt0dg3jw4p0"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgCreateDataSource(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1000))
	msg := NewMsgCreateDataSource(owner, "data_source_1", fee, []byte("executable"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "createDataSource", msg.Type())
	require.Equal(t, owner, msg.Owner)
	require.Equal(t, "data_source_1", msg.Name)
	require.Equal(t, fee, msg.Fee)
	require.Equal(t, []byte("executable"), msg.Executable)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgCreateDataSourceValidation(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	name := "data_source_1"
	executable := []byte("executable")
	feeUband10 := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	feeUband0 := sdk.NewCoins(sdk.NewInt64Coin("uband", 0))

	cases := []struct {
		valid bool
		tx    MsgCreateDataSource
	}{
		{
			true, NewMsgCreateDataSource(owner, name, feeUband10, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(nil, name, feeUband10, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, "", feeUband10, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, feeUband0, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, feeUband10, []byte(""), sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, feeUband10, nil, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, feeUband10, executable, nil),
		},
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgCreateDataSourceGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	msg := NewMsgCreateDataSource(owner, "data_source_1", fee, []byte("executable"), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/CreateDataSource","value":{"executable":"ZXhlY3V0YWJsZQ==","fee":[{"amount":"10","denom":"uband"}],"name":"data_source_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}
