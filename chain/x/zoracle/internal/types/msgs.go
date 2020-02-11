package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgRequest defines a Request message
type MsgRequest struct {
	CodeHash     []byte         `json:"codeHash"`
	Params       []byte         `json:"params"`
	ReportPeriod uint64         `json:"reportPeriod"`
	Sender       sdk.AccAddress `json:"sender"`
}

// NewMsgRequest is a constructor function for MsgRequest
func NewMsgRequest(
	codeHash []byte,
	params []byte,
	reportPeriod uint64,
	sender sdk.AccAddress,
) MsgRequest {
	return MsgRequest{
		CodeHash:     codeHash,
		Params:       params,
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
	if msg.CodeHash == nil || len(msg.CodeHash) != 32 {
		return sdk.ErrUnknownRequest("CodeHash must contain 32 bytes")
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

// MsgRequestData is a new version of request message.
type MsgRequestData struct {
	OracleScriptID           int64          `json:"oracleScriptID"`
	Calldata                 []byte         `json:"calldata"`
	RequestedValidatorCount  int64          `json:"requestedValidatorCount"`
	SufficientValidatorCount int64          `json:"sufficientValidatorCount"`
	Expiration               int64          `json:"expiration"`
	Sender                   sdk.AccAddress `json:"sender"`
}

// NewMsgRequestData is a constructor function for MsgRequestData
func NewMsgRequestData(
	oracleScriptID int64,
	calldata []byte,
	requestedValidatorCount int64,
	sufficientValidatorCount int64,
	expiration int64,
	sender sdk.AccAddress,
) MsgRequestData {
	return MsgRequestData{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidatorCount:  requestedValidatorCount,
		SufficientValidatorCount: sufficientValidatorCount,
		Expiration:               expiration,
		Sender:                   sender,
	}
}

// Route should return the name of the module
func (msg MsgRequestData) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestData) Type() string { return "request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRequestData) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	if msg.OracleScriptID <= 0 {
		return sdk.ErrUnknownRequest("Oracle script id must be greater than zero")
	}
	if msg.SufficientValidatorCount <= 0 {
		return sdk.ErrUnknownRequest("Sufficient validator count must be greater than zero")
	}
	if msg.RequestedValidatorCount < msg.SufficientValidatorCount {
		return sdk.ErrUnknownRequest("Request validator count must be greater than sufficient validator")
	}
	if msg.Expiration <= 0 {
		return sdk.ErrUnknownRequest("Expiration period must be greater than zero")
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg MsgRequestData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestData) GetSignBytes() []byte {
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

// MsgReportData defines a new version Report message
type MsgReportData struct {
	RequestID int64          `json:"requestID"`
	Data      []ExternalData `json:"data"`
	Sender    sdk.ValAddress `json:"sender"`
}

// NewMsgReportData is a constructor function for MsgReportData
func NewMsgReportData(
	requestID int64,
	data []ExternalData,
	sender sdk.ValAddress,
) MsgReportData {
	return MsgReportData{
		RequestID: requestID,
		Data:      data,
		Sender:    sender,
	}
}

// Route should return the name of the module
func (msg MsgReportData) Route() string { return RouterKey }

// Type should return the action
func (msg MsgReportData) Type() string { return "report" }

// ValidateBasic runs stateless checks on the message
func (msg MsgReportData) ValidateBasic() sdk.Error {
	if msg.RequestID <= 0 {
		return sdk.ErrUnknownRequest("Request id must be greater than zero")
	}
	if msg.Data == nil {
		return sdk.ErrUnknownRequest("Data must not be empty struct")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg MsgReportData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

// GetSignBytes encodes the message for signing
func (msg MsgReportData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgStoreCode defines a Code and owner of this
type MsgStoreCode struct {
	Code  []byte         `json:"code"`
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgStoreCode is a constructor function for MsgReport
func NewMsgStoreCode(
	code []byte,
	name string,
	owner sdk.AccAddress,
) MsgStoreCode {
	return MsgStoreCode{
		Code:  code,
		Name:  name,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgStoreCode) Route() string { return RouterKey }

// Type should return the action
func (msg MsgStoreCode) Type() string { return "store" }

// ValidateBasic runs stateless checks on the message
func (msg MsgStoreCode) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name must not be empty")
	}
	if msg.Code == nil || len(msg.Code) == 0 {
		return sdk.ErrUnknownRequest("Code must not be empty")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSignBytes encodes the message for signing
func (msg MsgStoreCode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

type MsgDeleteCode struct {
	CodeHash []byte         `json:"codeHash"`
	Owner    sdk.AccAddress `json:"owner"`
}

// NewMsgDeleteCode is a constructor function for MsgReport
func NewMsgDeleteCode(
	codeHash []byte,
	owner sdk.AccAddress,
) MsgDeleteCode {
	return MsgDeleteCode{
		CodeHash: codeHash,
		Owner:    owner,
	}
}

// Route should return the name of the module
func (msg MsgDeleteCode) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeleteCode) Type() string { return "delete" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeleteCode) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.CodeHash == nil || len(msg.CodeHash) != 32 {
		return sdk.ErrUnknownRequest("CodeHash must contain 32 bytes")
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgDeleteCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSignBytes encodes the message for signing
func (msg MsgDeleteCode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
