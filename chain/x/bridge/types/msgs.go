package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
)

// MsgRelayAndVerify is a message for relay and verify proof
type MsgRelayAndVerify struct {
	Proof  tmmerkle.Proof
	Sender sdk.AccAddress
}

// NewMsgRelayAndVerify creates a new MsgRelayAndVerify instance.
func NewMsgRelayAndVerify(
	proof tmmerkle.Proof,
	sender sdk.AccAddress,
) MsgRelayAndVerify {
	return MsgRelayAndVerify{
		Proof:  proof,
		Sender: sender,
	}
}

// RouterKey is the name of the oracle module
const RouterKey = ModuleName

// Route implements the sdk.Msg interface for MsgRelayAndVerify.
func (msg MsgRelayAndVerify) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRelayAndVerify.
func (msg MsgRelayAndVerify) Type() string { return "relayAndVerify" }

// ValidateBasic implements the sdk.Msg interface for MsgRelayAndVerify.
func (msg MsgRelayAndVerify) ValidateBasic() error {
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgRelayAndVerify.
func (msg MsgRelayAndVerify) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRelayAndVerify.
func (msg MsgRelayAndVerify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
