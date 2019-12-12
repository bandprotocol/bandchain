package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgRequest defines a Request message
type MsgRequest struct {
	Code         []byte         `json:"code"`
	ReportPeriod uint64         `json:"reportPeriod"`
	Sender       sdk.AccAddress `json:"sender"`
}

// NewMsgRequest is a constructor function for MsgRequest
func NewMsgRequest(
	code []byte,
	reportPeriod uint64,
	sender sdk.AccAddress,
) MsgRequest {
	return MsgRequest{
		Code:         code,
		ReportPeriod: reportPeriod,
		Sender:       sender,
	}
}

// Route should return the name of the module
func (msg MsgRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequest) Type() string { return "request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRequest) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if msg.Code == nil || len(msg.Code) == 0 {
		return sdk.ErrUnknownRequest("Code must not be empty bytes")
	}
	if msg.ReportPeriod <= 0 {
		return sdk.ErrInternal("Report period must be greater than zero")
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg MsgRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes encodes the message for signing
func (msg MsgRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgReport defines a Report message
type MsgReport struct {
	RequestID uint64         `json:"requestID"`
	Data      []byte         `json:"data"`
	Validator sdk.ValAddress `json:"validator"`
}

// NewMsgReport is a constructor function for MsgReport
func NewMsgReport(
	requestID uint64,
	data []byte,
	validator sdk.ValAddress,
) MsgReport {
	return MsgReport{
		RequestID: requestID,
		Data:      data,
		Validator: validator,
	}
}

// Route should return the name of the module
func (msg MsgReport) Route() string { return RouterKey }

// Type should return the action
func (msg MsgReport) Type() string { return "report" }

// ValidateBasic runs stateless checks on the message
func (msg MsgReport) ValidateBasic() sdk.Error {
	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress(msg.Validator.String())
	}

	if msg.Data == nil || len(msg.Data) == 0 {
		return sdk.ErrUnknownRequest("Data must not be empty bytes")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgReport) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator.Bytes()}
}

// GetSignBytes encodes the message for signing
func (msg MsgReport) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
