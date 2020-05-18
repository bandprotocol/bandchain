package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	GoodTestAddr    = sdk.AccAddress(make([]byte, 20))
	BadTestAddr     = sdk.AccAddress([]byte("BAD_ADDR"))
	GoodTestValAddr = sdk.ValAddress(make([]byte, 20))
	BadTestValAddr  = sdk.ValAddress([]byte("BAD_ADDR"))
)

type validateTestCase struct {
	valid bool
	msg   sdk.Msg
}

func performValidateTests(t *testing.T, cases []validateTestCase) {
	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestMsgRoute(t *testing.T) {
	require.Equal(t, "oracle", MsgCreateDataSource{}.Route())
	require.Equal(t, "oracle", MsgEditDataSource{}.Route())
	require.Equal(t, "oracle", MsgCreateOracleScript{}.Route())
	require.Equal(t, "oracle", MsgEditOracleScript{}.Route())
	require.Equal(t, "oracle", MsgRequestData{}.Route())
	require.Equal(t, "oracle", MsgReportData{}.Route())
	require.Equal(t, "oracle", MsgAddOracleAddress{}.Route())
	require.Equal(t, "oracle", MsgRemoveOracleAddress{}.Route())
}

func TestMsgType(t *testing.T) {
	require.Equal(t, "create_data_source", MsgCreateDataSource{}.Type())
	require.Equal(t, "edit_data_source", MsgEditDataSource{}.Type())
	require.Equal(t, "create_oracle_script", MsgCreateOracleScript{}.Type())
	require.Equal(t, "edit_oracle_script", MsgEditOracleScript{}.Type())
	require.Equal(t, "request", MsgRequestData{}.Type())
	require.Equal(t, "report", MsgReportData{}.Type())
	require.Equal(t, "add_oracle_address", MsgAddOracleAddress{}.Type())
	require.Equal(t, "remove_oracle_address", MsgRemoveOracleAddress{}.Type())
}

func TestMsgGetSigners(t *testing.T) {
	signerAcc := sdk.AccAddress([]byte("01234567890123456789"))
	signerVal := sdk.ValAddress([]byte("01234567890123456789"))
	anotherAcc := sdk.AccAddress([]byte("98765432109876543210"))
	anotherVal := sdk.ValAddress([]byte("98765432109876543210"))
	signers := []sdk.AccAddress{signerAcc}
	require.Equal(t, signers, NewMsgCreateDataSource(anotherAcc, "name", "desc", []byte("exec"), signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgEditDataSource(1, anotherAcc, "name", "desc", []byte("exec"), signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgCreateOracleScript(anotherAcc, "name", "desc", []byte("code"), "schema", "url", signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgEditOracleScript(1, anotherAcc, "name", "desc", []byte("code"), "schema", "url", signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, anotherVal, signerAcc).GetSigners())
	require.Equal(t, signers, NewMsgAddOracleAddress(signerVal, anotherAcc).GetSigners())
	require.Equal(t, signers, NewMsgRemoveOracleAddress(signerVal, anotherAcc).GetSigners())
}

func TestMsgGetSignBytes(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)
	require.Equal(t,
		`{"type":"oracle/CreateDataSource","value":{"description":"desc","executable":"ZXhlYw==","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
		string(NewMsgCreateDataSource(GoodTestAddr, "name", "desc", []byte("exec"), GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/EditDataSource","value":{"data_source_id":"1","description":"desc","executable":"ZXhlYw==","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
		string(NewMsgEditDataSource(1, GoodTestAddr, "name", "desc", []byte("exec"), GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/CreateOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","schema":"schema","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","source_code_url":"url"}}`,
		string(NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/EditOracleScript","value":{"code":"Y29kZQ==","description":"desc","name":"name","oracle_script_id":"1","owner":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","schema":"schema","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","source_code_url":"url"}}`,
		string(NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/Request","value":{"ask_count":"10","calldata":"Y2FsbGRhdGE=","client_id":"client-id","min_count":"5","oracle_script_id":"1","sender":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4"}}`,
		string(NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/Report","value":{"data_set":[{"data":"ZGF0YTE=","exit_code":1,"external_id":"1"},{"data":"ZGF0YTI=","exit_code":2,"external_id":"2"}],"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","request_id":"1","validator":"cosmosvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqkh52tw"}}`,
		string(NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/AddOracleAddress","value":{"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","validator":"cosmosvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqkh52tw"}}`,
		string(NewMsgAddOracleAddress(GoodTestValAddr, GoodTestAddr).GetSignBytes()),
	)
	require.Equal(t,
		`{"type":"oracle/RemoveOracleAddress","value":{"reporter":"band1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2vqal4","validator":"cosmosvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqkh52tw"}}`,
		string(NewMsgRemoveOracleAddress(GoodTestValAddr, GoodTestAddr).GetSignBytes()),
	)
}

func TestMsgCreateDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgCreateDataSource(GoodTestAddr, "name", "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgCreateDataSource(BadTestAddr, "name", "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgCreateDataSource(GoodTestAddr, strings.Repeat("x", 200), "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgCreateDataSource(GoodTestAddr, "name", strings.Repeat("x", 5000), []byte("exec"), GoodTestAddr)},
		{false, NewMsgCreateDataSource(GoodTestAddr, "name", "desc", []byte{}, GoodTestAddr)},
		{false, NewMsgCreateDataSource(GoodTestAddr, "name", "desc", []byte(strings.Repeat("x", 20000)), GoodTestAddr)},
		{false, NewMsgCreateDataSource(GoodTestAddr, "name", "desc", []byte("exec"), BadTestAddr)},
	})
}

func TestMsgEditDataSourceValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgEditDataSource(1, GoodTestAddr, "name", "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgEditDataSource(1, BadTestAddr, "name", "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgEditDataSource(1, GoodTestAddr, strings.Repeat("x", 200), "desc", []byte("exec"), GoodTestAddr)},
		{false, NewMsgEditDataSource(1, GoodTestAddr, "name", strings.Repeat("x", 5000), []byte("exec"), GoodTestAddr)},
		{false, NewMsgEditDataSource(1, GoodTestAddr, "name", "desc", []byte{}, GoodTestAddr)},
		{false, NewMsgEditDataSource(1, GoodTestAddr, "name", "desc", []byte(strings.Repeat("x", 20000)), GoodTestAddr)},
		{false, NewMsgEditDataSource(1, GoodTestAddr, "name", "desc", []byte("exec"), BadTestAddr)},
	})
}

func TestMsgCreateOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(BadTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, strings.Repeat("x", 200), "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", strings.Repeat("x", 5000), []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte("code"), strings.Repeat("x", 1000), "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte("code"), "schema", strings.Repeat("x", 200), GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte{}, "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte(strings.Repeat("x", 600000)), "schema", "url", GoodTestAddr)},
		{false, NewMsgCreateOracleScript(GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", BadTestAddr)},
	})
}

func TestMsgEditOracleScriptValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, BadTestAddr, "name", "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, strings.Repeat("x", 200), "desc", []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", strings.Repeat("x", 5000), []byte("code"), "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte("code"), strings.Repeat("x", 1000), "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte("code"), "schema", strings.Repeat("x", 200), GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte{}, "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte(strings.Repeat("x", 600000)), "schema", "url", GoodTestAddr)},
		{false, NewMsgEditOracleScript(1, GoodTestAddr, "name", "desc", []byte("code"), "schema", "url", BadTestAddr)},
	})
}

func TestMsgRequestDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte(strings.Repeat("x", 2000)), 10, 5, "client-id", GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 2, 5, "client-id", GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 0, 0, "client-id", GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, strings.Repeat("x", 300), GoodTestAddr)},
		{false, NewMsgRequestData(1, []byte("calldata"), 10, 5, "client-id", BadTestAddr)},
	})
}

func TestMsgReportDataValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte(strings.Repeat("x", 2000))}, {2, 2, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {1, 1, []byte("data2")}}, GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, BadTestValAddr, GoodTestAddr)},
		{false, NewMsgReportData(1, []RawReport{{1, 1, []byte("data1")}, {2, 2, []byte("data2")}}, GoodTestValAddr, BadTestAddr)},
	})
}

func TestMsgAddOracleAddressValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgAddOracleAddress(GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgAddOracleAddress(BadTestValAddr, GoodTestAddr)},
		{false, NewMsgAddOracleAddress(GoodTestValAddr, BadTestAddr)},
	})
}

func TestMsgRemoveOracleAddressValidation(t *testing.T) {
	performValidateTests(t, []validateTestCase{
		{true, NewMsgRemoveOracleAddress(GoodTestValAddr, GoodTestAddr)},
		{false, NewMsgRemoveOracleAddress(BadTestValAddr, GoodTestAddr)},
		{false, NewMsgRemoveOracleAddress(GoodTestValAddr, BadTestAddr)},
	})
}
