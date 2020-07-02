package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
)

// MsgVerifyProof is a message for relay and verify proof
type MsgVerifyProof struct {
	Proof  tmmerkle.Proof
	Sender sdk.AccAddress
}

// NewMsgVerifyProof creates a new MsgVerifyProof instance.
func NewMsgVerifyProof(
	proof tmmerkle.Proof,
	sender sdk.AccAddress,
) MsgVerifyProof {
	return MsgVerifyProof{
		Proof:  proof,
		Sender: sender,
	}
}

// RouterKey is the name of the oracle module
const RouterKey = ModuleName

// Route implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) Type() string { return "verifyProof" }

// ValidateBasic implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) ValidateBasic() error {
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
