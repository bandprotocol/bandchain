package types

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgRequestData is a message for requesting a new data request to an existing oracle script.
type MsgRequestData struct {
	OracleScriptID oracle.OracleScriptID `json:"oracleScriptID"`
	SourceChannel  string                `json:"sourceChannel"`
	ClientID       string                `json:"clientID"`
	Calldata       []byte                `json:"calldata"`
	AskCount       int64                 `json:"askCount"`
	MinCount       int64                 `json:"minCount"`
	Sender         sdk.AccAddress        `json:"sender"`
}

// NewMsgRequestData creates a new MsgRequestData instance.
func NewMsgRequestData(
	oracleScriptID oracle.OracleScriptID,
	sourceChannel string,
	clientID string,
	calldata []byte,
	askCount int64,
	minCount int64,
	sender sdk.AccAddress,
) MsgRequestData {
	return MsgRequestData{
		OracleScriptID: oracleScriptID,
		SourceChannel:  sourceChannel,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
		ClientID:       clientID,
		Sender:         sender,
	}
}

// Route implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Type() string { return "consuming" }

// ValidateBasic implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgRequestData: Sender address must not be empty.")
	}
	if msg.OracleScriptID <= 0 {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgRequestData: Oracle script id (%d) must be positive.", msg.OracleScriptID)
	}
	if msg.AskCount <= 0 {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg,
			"MsgRequestData: Ask validator count (%d) must be positive.",
			msg.AskCount,
		)
	}
	if msg.AskCount < msg.MinCount {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg,
			"MsgRequestData: Request validator count (%d) must not be less than minimum validator count (%d).",
			msg.AskCount,
			msg.MinCount,
		)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
