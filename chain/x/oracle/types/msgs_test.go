package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRequestData(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 5000, 10000, "clientID", sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "request", msg.Type())
	require.Equal(t, OracleScriptID(1), msg.OracleScriptID)
	require.Equal(t, []byte("calldata"), msg.Calldata)
	require.Equal(t, int64(10), msg.RequestedValidatorCount)
	require.Equal(t, int64(5), msg.SufficientValidatorCount)
	require.Equal(t, uint64(5000), msg.PrepareGas)
	require.Equal(t, uint64(10000), msg.ExecuteGas)
	require.Equal(t, sender, msg.Sender)
}

func TestMsgRequestDataValidation(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	requestedValidatorCount := int64(10)
	sufficientValidatorCount := int64(5)
	prepareGas := uint64(5000)
	executeGas := uint64(10000)
	clientID := "clientID"
	cases := []struct {
		valid bool
		tx    MsgRequestData
	}{
		{
			true, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				0, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			true, NewMsgRequestData(
				1, nil, requestedValidatorCount,
				sufficientValidatorCount, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 0,
				sufficientValidatorCount, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				-1, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), 6,
				8, prepareGas, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, 0, executeGas, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, prepareGas, 0, clientID, sender,
			),
		},
		{
			false, NewMsgRequestData(
				1, []byte("calldata"), requestedValidatorCount,
				sufficientValidatorCount, prepareGas, executeGas, clientID, nil,
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
	msg := NewMsgRequestData(1, []byte("calldata"), 10, 5, 5000, 10000, "clientID", sender)
	res := msg.GetSignBytes()

	expected := `{"type":"oracle/Request","value":{"calldata":"Y2FsbGRhdGE=","clientID":"clientID","executeGas":"10000","oracleScriptID":"1","prepareGas":"5000","requestedValidatorCount":"10","sender":"band1wdjkuer9wgvz7c4y","sufficientValidatorCount":"5"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgReportData(t *testing.T) {
	requestID := RequestID(3)
	data := []RawDataReportWithID{NewRawDataReportWithID(1, 1, []byte("data1")), NewRawDataReportWithID(2, 2, []byte("data2"))}
	provider, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	reporter := sdk.AccAddress(provider)

	msg := NewMsgReportData(requestID, data, provider, reporter)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "report", msg.Type())
}

func TestMsgReportDataValidation(t *testing.T) {
	requestID := RequestID(3)
	data := []RawDataReportWithID{NewRawDataReportWithID(1, 1, []byte("data1")), NewRawDataReportWithID(2, 2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	reporter := sdk.AccAddress(validator)
	failValidator, _ := sdk.ValAddressFromHex("")

	cases := []struct {
		valid bool
		tx    MsgReportData
	}{
		{true, NewMsgReportData(requestID, data, validator, reporter)},
		{false, NewMsgReportData(-1, data, validator, reporter)},
		{false, NewMsgReportData(requestID, []RawDataReportWithID{}, validator, reporter)},
		{false, NewMsgReportData(requestID, nil, validator, reporter)},
		{false, NewMsgReportData(requestID, data, failValidator, reporter)},
		{false, NewMsgReportData(requestID, data, failValidator, nil)},
		{false, NewMsgReportData(requestID, []RawDataReportWithID{NewRawDataReportWithID(1, 1, []byte("data1")), NewRawDataReportWithID(1, 2, []byte("data2"))}, validator, reporter)},
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

	requestID := RequestID(3)
	data := []RawDataReportWithID{NewRawDataReportWithID(1, 1, []byte("data1")), NewRawDataReportWithID(2, 2, []byte("data2"))}
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	reporter := sdk.AccAddress(validator)
	msg := NewMsgReportData(requestID, data, validator, reporter)
	res := msg.GetSignBytes()

	expected := `{"type":"oracle/Report","value":{"dataSet":[{"data":"ZGF0YTE=","exitCode":1,"externalDataID":"1"},{"data":"ZGF0YTI=","exitCode":2,"externalDataID":"2"}],"reporter":"band1hq8j5h0h64csk9tz95df783cxr0dt0dg3jw4p0","requestID":"3","validator":"bandvaloper1hq8j5h0h64csk9tz95df783cxr0dt0dgay2kyy"}}`

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

	expected := `{"type":"oracle/CreateDataSource","value":{"description":"description","executable":"ZXhlY3V0YWJsZQ==","fee":[{"amount":"10","denom":"uband"}],"name":"data_source_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgEditDataSource(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 100))
	msg := NewMsgEditDataSource(1, owner, "data_source_1", "description", fee, []byte("executable"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "edit_data_source", msg.Type())
	require.Equal(t, DataSourceID(1), msg.DataSourceID)
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

	expected := `{"type":"oracle/EditDataSource","value":{"dataSourceID":"1","description":"description","executable":"ZXhlY3V0YWJsZQ==","fee":[{"amount":"10","denom":"uband"}],"name":"data_source_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

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
	res := msg.GetSignBytes()

	expected := `{"type":"oracle/CreateOracleScript","value":{"code":"Y29kZQ==","description":"description","name":"oracle_script_1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgEditOracleScript(t *testing.T) {
	owner := sdk.AccAddress([]byte("owner"))
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgEditOracleScript(1, owner, "oracle_script_1", "description", []byte("code"), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "edit_oracle_script", msg.Type())
	require.Equal(t, OracleScriptID(1), msg.OracleScriptID)
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

	expected := `{"type":"oracle/EditOracleScript","value":{"code":"Y29kZQ==","description":"description","name":"oracle_script_1","oracleScriptID":"1","owner":"band1damkuetjcw3c0d","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgAddOracleAddress(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))
	msg := NewMsgAddOracleAddress(validator, reporter)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "add_oracle_address", msg.Type())
	require.Equal(t, validator, msg.Validator)
	require.Equal(t, reporter, msg.Reporter)
}

func TestMsgAddOracleAddressValidation(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))

	cases := []struct {
		valid bool
		tx    MsgAddOracleAddress
	}{
		{
			true, NewMsgAddOracleAddress(validator, reporter),
		},
		{
			false, NewMsgAddOracleAddress(nil, reporter),
		},
		{
			false, NewMsgAddOracleAddress(validator, nil),
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
func TestMsgAddOracleAddressGetSignBytes(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))
	msg := NewMsgAddOracleAddress(validator, reporter)
	res := msg.GetSignBytes()

	expected := `{"type":"oracle/AddOracleAddress","value":{"reporter":"band1wfjhqmmjw3jhyy3as6w","validator":"bandvaloper1weskc6tyv96x7usd82k92"}}`

	require.Equal(t, expected, string(res))

}

func TestMsgRemoveOracleAddress(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))
	msg := NewMsgRemoveOracleAddress(validator, reporter)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "remove_oracle_address", msg.Type())
	require.Equal(t, validator, msg.Validator)
	require.Equal(t, reporter, msg.Reporter)
}

func TestMsgRemoveOracleAddressValidation(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))

	cases := []struct {
		valid bool
		tx    MsgRemoveOracleAddress
	}{
		{
			true, NewMsgRemoveOracleAddress(validator, reporter),
		},
		{
			false, NewMsgRemoveOracleAddress(nil, reporter),
		},
		{
			false, NewMsgRemoveOracleAddress(validator, nil),
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

func TestMsgRemoveOracleAddressGetSignBytes(t *testing.T) {
	validator := sdk.ValAddress([]byte("validator"))
	reporter := sdk.AccAddress([]byte("reporter"))
	msg := NewMsgRemoveOracleAddress(validator, reporter)
	res := msg.GetSignBytes()

	expected := `{"type":"oracle/RemoveOracleAddress","value":{"reporter":"band1wfjhqmmjw3jhyy3as6w","validator":"bandvaloper1weskc6tyv96x7usd82k92"}}`

	require.Equal(t, expected, string(res))

}
