package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRequest(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequest([]byte("code"), uint64(10), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "request", msg.Type())
}

func TestMsgRequestValidation(t *testing.T) {
	code := []byte("Code")
	sender := sdk.AccAddress([]byte("sender"))
	reportPeriod := uint64(10)
	cases := []struct {
		valid bool
		tx    MsgRequest
	}{
		{true, NewMsgRequest(code, reportPeriod, sender)},
		{false, NewMsgRequest([]byte{}, reportPeriod, sender)},
		{false, NewMsgRequest(nil, reportPeriod, sender)},
		{false, NewMsgRequest(code, reportPeriod, sdk.AccAddress([]byte("")))},
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

func TestMsgRequestGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("band", "band"+sdk.PrefixPublic)

	code := []byte("Code")
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequest(code, uint64(10), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Request","value":{"code":"Q29kZQ==","reportPeriod":"10","sender":"band1wdjkuer9wgvz7c4y"}}`

	require.Equal(t, expected, string(res))
}

func TestMsgReport(t *testing.T) {
	requestID := uint64(3)
	data := []byte("Data")
	provider, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReport(requestID, data, provider)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "report", msg.Type())
}

func TestMsgReportValidation(t *testing.T) {
	requestID := uint64(3)
	data := []byte("Data")
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failValidator, _ := sdk.ValAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgReport
	}{
		{true, NewMsgReport(requestID, data, validator)},
		{false, NewMsgReport(requestID, []byte(""), validator)},
		{false, NewMsgReport(requestID, nil, validator)},
		{false, NewMsgReport(requestID, data, failValidator)},
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

func TestMsgReportGetSignBytes(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForValidator("band"+sdk.PrefixValidator+sdk.PrefixOperator, "band"+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)

	requestID := uint64(3)
	data := []byte("Data")
	validator, _ := sdk.ValAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgReport(requestID, data, validator)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Report","value":{"data":"RGF0YQ==","requestID":"3","validator":"bandvaloper1hq8j5h0h64csk9tz95df783cxr0dt0dgay2kyy"}}`

	require.Equal(t, expected, string(res))
}
