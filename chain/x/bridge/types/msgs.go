package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
	tmtypes "github.com/tendermint/tendermint/types"

	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// RouterKey is the name of the bridge module
const RouterKey = ModuleName

// MsgUpdateChainID is a message for update chain ID of BandChain.
type MsgUpdateChainID struct {
	ChainID string
	Sender  sdk.AccAddress
}

// NewMsgUpdateChainID creates a new MsgUpdateChainID instance.
func NewMsgUpdateChainID(
	chainID string,
	sender sdk.AccAddress,
) MsgUpdateChainID {
	return MsgUpdateChainID{
		ChainID: chainID,
		Sender:  sender,
	}
}

// Route implements the sdk.Msg interface for MsgUpdateChainID.
func (msg MsgUpdateChainID) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgUpdateChainID.
func (msg MsgUpdateChainID) Type() string { return "update_chain_id" }

// ValidateBasic implements the sdk.Msg interface for MsgUpdateChainID.
func (msg MsgUpdateChainID) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgUpdateChainID.
func (msg MsgUpdateChainID) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgUpdateChainID.
func (msg MsgUpdateChainID) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgUpdateValidators is a message for update chain ID of BandChain.
type MsgUpdateValidators struct {
	Validators []tmtypes.Validator
	Sender     sdk.AccAddress
}

// NewMsgUpdateValidators creates a new MsgUpdateValidators instance.
func NewMsgUpdateValidators(
	validators []tmtypes.Validator,
	sender sdk.AccAddress,
) MsgUpdateValidators {
	return MsgUpdateValidators{
		Validators: validators,
		Sender:     sender,
	}
}

// Route implements the sdk.Msg interface for MsgUpdateValidators.
func (msg MsgUpdateValidators) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgUpdateValidators.
func (msg MsgUpdateValidators) Type() string { return "update_validators" }

// ValidateBasic implements the sdk.Msg interface for MsgUpdateValidators.
func (msg MsgUpdateValidators) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgUpdateValidators.
func (msg MsgUpdateValidators) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgUpdateValidators.
func (msg MsgUpdateValidators) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgRelay is a message for relay block on BandChain.
type MsgRelay struct {
	SignedHeader tmtypes.SignedHeader
	Sender       sdk.AccAddress
}

// NewMsgRelay creates a new MsgRelay instance.
func NewMsgRelay(
	signedHeader tmtypes.SignedHeader,
	sender sdk.AccAddress,
) MsgRelay {
	return MsgRelay{
		SignedHeader: signedHeader,
		Sender:       sender,
	}
}

// Route implements the sdk.Msg interface for MsgRelay.
func (msg MsgRelay) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRelay.
func (msg MsgRelay) Type() string { return "relay" }

// ValidateBasic implements the sdk.Msg interface for MsgRelay.
func (msg MsgRelay) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgRelay.
func (msg MsgRelay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRelay.
func (msg MsgRelay) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgVerifyProof is a message for verify proof
type MsgVerifyProof struct {
	Height         int64
	Proof          tmmerkle.Proof
	RequestPacket  otypes.OracleRequestPacketData
	ResponsePacket otypes.OracleResponsePacketData
	Sender         sdk.AccAddress
}

// NewMsgVerifyProof creates a new MsgVerifyProof instance.
func NewMsgVerifyProof(
	height int64,
	proof tmmerkle.Proof,
	requestPacket otypes.OracleRequestPacketData,
	responsePacket otypes.OracleResponsePacketData,
	sender sdk.AccAddress,
) MsgVerifyProof {
	return MsgVerifyProof{
		Height:         height,
		Proof:          proof,
		RequestPacket:  requestPacket,
		ResponsePacket: responsePacket,
		Sender:         sender,
	}
}

// Route implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) Type() string { return "verify_proof" }

// ValidateBasic implements the sdk.Msg interface for MsgVerifyProof.
func (msg MsgVerifyProof) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
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
