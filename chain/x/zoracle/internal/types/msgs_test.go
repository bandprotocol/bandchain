package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRequest(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	msg := NewMsgRequest(codeHash, []byte("params"), uint64(10), sender)
	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "request", msg.Type())
}

func TestMsgRequestValidation(t *testing.T) {
	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	sender := sdk.AccAddress([]byte("sender"))
	reportPeriod := uint64(10)
	cases := []struct {
		valid bool
		tx    MsgRequest
	}{
		{true, NewMsgRequest(codeHash, []byte("params"), reportPeriod, sender)},
		{false, NewMsgRequest([]byte{}, []byte("params"), reportPeriod, sender)},
		{false, NewMsgRequest(nil, []byte("params"), reportPeriod, sender)},
		{false, NewMsgRequest(codeHash, []byte("params"), reportPeriod, sdk.AccAddress([]byte("")))},
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

	codeHash, _ := hex.DecodeString("5694d08a2e53ffcae0c3103e5ad6f6076abd960eb1f8a56577040bc1028f702b")
	sender := sdk.AccAddress([]byte("sender"))
	msg := NewMsgRequest(codeHash, []byte("params"), uint64(10), sender)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Request","value":{"codeHash":"VpTQii5T/8rgwxA+Wtb2B2q9lg6x+KVldwQLwQKPcCs=","params":"cGFyYW1z","reportPeriod":"10","sender":"band1wdjkuer9wgvz7c4y"}}`

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

func TestMsgStoreCode(t *testing.T) {
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgStoreCode([]byte("Code"), owner)

	require.Equal(t, RouterKey, msg.Route())
	require.Equal(t, "store", msg.Type())
}

func TestMsgStoreCodeValidation(t *testing.T) {
	code := []byte("Code")
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	failOwner, _ := sdk.AccAddressFromHex("")
	cases := []struct {
		valid bool
		tx    MsgStoreCode
	}{
		{true, NewMsgStoreCode(code, owner)},
		{false, NewMsgStoreCode([]byte(""), owner)},
		{false, NewMsgStoreCode(nil, owner)},
		{false, NewMsgStoreCode(code, failOwner)},
		{false, NewMsgStoreCode(code, nil)},
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

	code := []byte("Code")
	owner, _ := sdk.AccAddressFromHex("b80f2a5df7d5710b15622d1a9f1e3830ded5bda8")
	msg := NewMsgStoreCode(code, owner)
	res := msg.GetSignBytes()

	expected := `{"type":"zoracle/Store","value":{"code":"Q29kZQ==","owner":"band1hq8j5h0h64csk9tz95df783cxr0dt0dg3jw4p0"}}`

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
