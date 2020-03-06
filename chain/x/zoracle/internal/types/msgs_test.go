package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRequestData(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 100, 5000, 10000, sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "request", msg.Type())
	require.Equal(t, int64(1), msg.OracleScriptID)
	require.Equal(t, []byte("calldata"), msg.Calldata)
	require.Equal(t, int64(10), msg.RequestedValidatorCount)
	require.Equal(t, int64(5), msg.SufficientValidatorCount)
	require.Equal(t, int64(100), msg.Expiration)
	require.Equal(t, uint64(5000), msg.PrepareGas)
	require.Equal(t, uint64(10000), msg.ExecuteGas)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgRequestDataValidation(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	requestedValidatorCount := int64(10)
	sufficientValidatorCount := int64(5)
	expiration := int64(100)
	prepareGas := uint64(5000)
	executeGas := uint64(10000)
	cases := []struct {
		valid bool
		tx    MsgRequestData
	}{
		{
			true, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				0, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			true, NewMsgRequestData(
				1, nil, requestedValidatorCount,
				sufficientValidatorCount, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 0,
				sufficientValidatorCount, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				-1, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 6,
				8, expiration, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, -10, prepareGas, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, 0, executeGas, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, prepareGas, 0, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, expiration, prepareGas, executeGas, nil,
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
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 100, 5000, 10000, sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Request","value":{"calldata":"Y2FsbGRhdGE=","executeGas":"10000","expiration":"100","oracleScriptID":"1","prepareGas":"5000","requestedValidatorCount":"10","sender":"band1wdjkuer9wgvz7c4y","sufficientValidatorCount":"5"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgReportData(t *testing.T) {
	requestID := int64(3)
	data := []RawDataReport{NewRawDataReport(1, []byte("data1")), NewRawDataReport(2, []byte("data2"))}
	provider, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReportData(requestID, data, provider)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "report", msg.Type())
}

func TestMsgReportDataValidation(t *testing.T) {
	requestID := int64(3)
	data := []RawDataReport{NewRawDataReport(1, []byte("data1")), NewRawDataReport(2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failValidator, _ := sdk.ValAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgReportData
	}{
		{true, NewMsgReportData(requestID, data, validator)},
		{false, NewMsgReportData(-1, data, validator)},
		{false, NewMsgReportData(requestID, []RawDataReport{}, validator)},
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
	data := []RawDataReport{NewRawDataReport(1, []byte("data1")), NewRawDataReport(2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReportData(requestID, data, validator)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Report","value":{"dataSet":[{"data":"ZGF0YTE=","externalDataID":"1"},{"data":"ZGF0YTI=","externalDataID":"2"}],"requestID":"3","sender":"bandvaloper1hq8j5h0h64csk9tz95df783cxr0dt0dgay2kyy"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgCreateDataSource(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1000))
	msg := NewMsgCreateDataSource(owner, "data_source_1", "description", fee, []byte("executable"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "create_data_source", msg.Type())
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
	description := "description"
	executable := []byte("executable")
	feeUband10 := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	feeUband0 := sdk.NewCoins(sdk.NewInt64Coin("uband", 0))

	cases := []struct {
		valid bool
		tx    MsgCreateDataSource
	}{
		{
			true, NewMsgCreateDataSource(owner, name, description, feeUband10, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(nil, name, description, feeUband10, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, "", description, feeUband10, executable, sender),
		},
		{
			true, NewMsgCreateDataSource(owner, name, description, feeUband0, executable, sender),
		},
		{
			true, NewMsgCreateDataSource(owner, name, description, sdk.Coins{}, executable, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, description, feeUband10, []byte(""), sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, description, feeUband10, nil, sender),
		},
		{
			false, NewMsgCreateDataSource(owner, name, description, feeUband10, executable, nil),
		},
		{
			false, NewMsgCreateDataSource(owner, name, "", feeUband10, executable, nil),
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
	msg := NewMsgCreateDataSource(owner, "data_source_1", "description", fee, []byte("executable"), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/CreateDataSource","value":{"description":"description","executable":"ZXhlY3V0YWJsZQ==","fee":[{"amount":"10","denom":"uband"}],"name":"data_source_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgEditDataSource(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 100))
	msg := NewMsgEditDataSource(1, owner, "data_source_1", "description", fee, []byte("executable"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "edit_data_source", msg.Type())
	require.Equal(t, int64(1), msg.DataSourceID)
	require.Equal(t, owner, msg.Owner)
	require.Equal(t, "data_source_1", msg.Name)
	require.Equal(t, fee, msg.Fee)
	require.Equal(t, []byte("executable"), msg.Executable)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgEditDataSourceValidation(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	feeUband10 := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	feeUband0 := sdk.NewCoins(sdk.NewInt64Coin("uband", 0))

	cases := []struct {
		valid bool
		tx    MsgEditDataSource
	}{
		{
			true, NewMsgEditDataSource(1, owner, name, description, feeUband10, executable, sender),
		},
		{
			false, NewMsgEditDataSource(0, owner, name, description, feeUband10, executable, sender),
		},
		{
			false, NewMsgEditDataSource(1, nil, name, description, feeUband10, executable, sender),
		},
		{
			false, NewMsgEditDataSource(1, owner, "", description, feeUband10, executable, sender),
		},
		{
			true, NewMsgEditDataSource(1, owner, name, description, feeUband0, executable, sender),
		},
		{
			true, NewMsgEditDataSource(1, owner, name, description, sdk.Coins{}, executable, sender),
		},
		{
			false, NewMsgEditDataSource(1, owner, name, description, feeUband10, []byte(""), sender),
		},
		{
			false, NewMsgEditDataSource(1, owner, name, description, feeUband10, nil, sender),
		},
		{
			false, NewMsgEditDataSource(1, owner, name, description, feeUband10, executable, nil),
		},
		{
			false, NewMsgEditDataSource(1, owner, name, "", feeUband10, executable, sender),
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

func TestMsgEditDataSourceGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source_1"
	description := "description"
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	msg := NewMsgEditDataSource(1, owner, name, description, fee, []byte("executable"), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/EditDataSource","value":{"dataSourceID":"1","description":"description","executable":"ZXhlY3V0YWJsZQ==","fee":[{"amount":"10","denom":"uband"}],"name":"data_source_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgCreateOracleScript(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgCreateOracleScript(owner, "oracle_script_1", "description", []byte("code"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "create_oracle_script", msg.Type())
	require.Equal(t, owner, msg.Owner)
	require.Equal(t, "oracle_script_1", msg.Name)
	require.Equal(t, []byte("code"), msg.Code)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgCreateOracleScriptValidation(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	description := "description"
	sender := sdk.AccAddress([]byte("sender"))
	name := "oracle_script_1"
	code := []byte("code")

	cases := []struct {
		valid bool
		tx    MsgCreateOracleScript
	}{
		{
			true, NewMsgCreateOracleScript(owner, name, description, code, sender),
		},
		{
			false, NewMsgCreateOracleScript(nil, name, description, code, sender),
		},
		{
			false, NewMsgCreateOracleScript(owner, "", description, code, sender),
		},
		{
			false, NewMsgCreateOracleScript(owner, name, description, []byte{}, sender),
		},
		{
			false, NewMsgCreateOracleScript(owner, name, description, nil, sender),
		},
		{
			false, NewMsgCreateOracleScript(owner, name, description, code, nil),
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

func TestMsgCreateOracleScriptGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgCreateOracleScript(owner, "oracle_script_1", "description", []byte("code"), sender)
	fmt.Println(msg)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/CreateOracleScript","value":{"code":"Y29kZQ==","description":"description","name":"oracle_script_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgEditOracleScript(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgEditOracleScript(1, owner, "oracle_script_1", "description", []byte("code"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "edit_oracle_script", msg.Type())
	require.Equal(t, int64(1), msg.OracleScriptID)
	require.Equal(t, owner, msg.Owner)
	require.Equal(t, "oracle_script_1", msg.Name)
	require.Equal(t, []byte("code"), msg.Code)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgEditOracleScriptValidation(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	name := "oracle_script_1"
	description := "description"
	code := []byte("code")

	cases := []struct {
		valid bool
		tx    MsgEditOracleScript
	}{
		{
			true, NewMsgEditOracleScript(1, owner, name, description, code, sender),
		},
		{
			false, NewMsgEditOracleScript(0, nil, name, description, code, sender),
		},
		{
			false, NewMsgEditOracleScript(1, nil, name, description, code, sender),
		},
		{
			false, NewMsgEditOracleScript(1, owner, "", description, code, sender),
		},
		{
			false, NewMsgEditOracleScript(1, owner, name, description, []byte{}, sender),
		},
		{
			false, NewMsgEditOracleScript(1, owner, name, description, nil, sender),
		},
		{
			false, NewMsgEditOracleScript(1, owner, name, description, code, nil),
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

func TestMsgEditOracleScriptGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgEditOracleScript(1, owner, "oracle_script_1", "description", []byte("code"), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/EditOracleScript","value":{"code":"Y29kZQ==","description":"description","name":"oracle_script_1","oracleScriptID":"1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}
