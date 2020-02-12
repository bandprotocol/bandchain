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

// NewMsgRequest creates a new MsgRequest instance.
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

// Route implements the sdk.Msg interface for MsgRequest.
func (msg MsgRequest) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRequest.
func (msg MsgRequest) Type() string { return "request" }

// ValidateBasic implements the sdk.Msg interface for MsgRequest.
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

// GetSigners implements the sdk.Msg interface for MsgRequest.
func (msg MsgRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRequest.
func (msg MsgRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgRequestData is a message for requesting a new data request to an existing oracle script.
type MsgRequestData struct {
	OracleScriptID           int64          `json:"oracleScriptID"`
	Calldata                 []byte         `json:"calldata"`
	RequestedValidatorCount  int64          `json:"requestedValidatorCount"`
	SufficientValidatorCount int64          `json:"sufficientValidatorCount"`
	Expiration               int64          `json:"expiration"`
	Sender                   sdk.AccAddress `json:"sender"`
}

// NewMsgRequestData creates a new MsgRequestData instance.
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

// Route implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Type() string { return "request" }

// ValidateBasic implements the sdk.Msg interface for MsgRequestData.
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

// GetSigners implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRequestData.
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

// NewMsgReport creates a new MsgReport instance.
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

// Route implements the sdk.Msg interface for MsgReport.
func (msg MsgReport) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgReport.
func (msg MsgReport) Type() string { return "report" }

// ValidateBasic implements the sdk.Msg interface for MsgReport.
func (msg MsgReport) ValidateBasic() sdk.Error {
	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress(msg.Validator.String())
	}

	if msg.Data == nil || len(msg.Data) == 0 {
		return sdk.ErrUnknownRequest("Data must not be empty bytes")
	}

	return nil
}

// GetSigners implements the sdk.Msg interface for MsgReport.
func (msg MsgReport) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface for MsgReport.
func (msg MsgReport) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgReportData is a message sent by each of the block validators to respond to a data request.
type MsgReportData struct {
	RequestID int64          `json:"requestID"`
	DataSet   []ExternalData `json:"dataSet"`
	Sender    sdk.ValAddress `json:"sender"`
}

// NewMsgReportData creates a new MsgReportData instance.
func NewMsgReportData(
	requestID int64,
	dataSet []ExternalData,
	sender sdk.ValAddress,
) MsgReportData {
	return MsgReportData{
		RequestID: requestID,
		DataSet:   dataSet,
		Sender:    sender,
	}
}

// Route implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) Type() string { return "report" }

// ValidateBasic implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) ValidateBasic() sdk.Error {
	if msg.RequestID <= 0 {
		return sdk.ErrUnknownRequest("Request id must be greater than zero")
	}
	if msg.DataSet == nil || len(msg.DataSet) == 0 {
		return sdk.ErrUnknownRequest("Data must not be empty struct")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Sender)}
}

// GetSignBytes implements the sdk.Msg interface for MsgReportData.
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

// NewMsgStoreCode creates a new MsgStoreCode instance.
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

// Route implements the sdk.Msg interface for MsgStoreCode.
func (msg MsgStoreCode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgStoreCode.
func (msg MsgStoreCode) Type() string { return "store" }

// ValidateBasic implements the sdk.Msg interface for MsgStoreCode.
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

// GetSigners implements the sdk.Msg interface for MsgStoreCode.
func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSignBytes implements the sdk.Msg interface for MsgStoreCode.
func (msg MsgStoreCode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

type MsgDeleteCode struct {
	CodeHash []byte         `json:"codeHash"`
	Owner    sdk.AccAddress `json:"owner"`
}

// NewMsgDeleteCode creates a new MsgDeleteCode instance.
func NewMsgDeleteCode(
	codeHash []byte,
	owner sdk.AccAddress,
) MsgDeleteCode {
	return MsgDeleteCode{
		CodeHash: codeHash,
		Owner:    owner,
	}
}

// Route implements the sdk.Msg interface for MsgDeleteCode.
func (msg MsgDeleteCode) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgDeleteCode.
func (msg MsgDeleteCode) Type() string { return "delete" }

// ValidateBasic implements the sdk.Msg interface for MsgDeleteCode.
func (msg MsgDeleteCode) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	if msg.CodeHash == nil || len(msg.CodeHash) != 32 {
		return sdk.ErrUnknownRequest("CodeHash must contain 32 bytes")
	}

	return nil
}

// GetSigners implements the sdk.Msg interface for MsgDeleteCode.
func (msg MsgDeleteCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSignBytes implements the sdk.Msg interface for MsgDeleteCode.
func (msg MsgDeleteCode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
