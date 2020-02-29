package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the bank module
const RouterKey = ModuleName

// MsgRequestData is a message for requesting a new data request to an existing oracle script.
type MsgRequestData struct {
	OracleScriptID           int64          `json:"oracleScriptID"`
	Calldata                 []byte         `json:"calldata"`
	RequestedValidatorCount  int64          `json:"requestedValidatorCount"`
	SufficientValidatorCount int64          `json:"sufficientValidatorCount"`
	Expiration               int64          `json:"expiration"`
	PrepareGas               uint64         `json:"prepareGas"`
	ExecuteGas               uint64         `json:"executeGas"`
	Sender                   sdk.AccAddress `json:"sender"`
}

// NewMsgRequestData creates a new MsgRequestData instance.
func NewMsgRequestData(
	oracleScriptID int64,
	calldata []byte,
	requestedValidatorCount int64,
	sufficientValidatorCount int64,
	expiration int64,
	prepareGas uint64,
	executeGas uint64,
	sender sdk.AccAddress,
) MsgRequestData {
	return MsgRequestData{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidatorCount:  requestedValidatorCount,
		SufficientValidatorCount: sufficientValidatorCount,
		Expiration:               expiration,
		PrepareGas:               prepareGas,
		ExecuteGas:               executeGas,
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
		return ErrInvalidBasicMsg("MsgRequestData: Sender address must not be empty.")
	}
	if msg.OracleScriptID <= 0 {
		return ErrInvalidBasicMsg("MsgRequestData: Oracle script id (%d) must be positive.", msg.OracleScriptID)
	}
	if msg.SufficientValidatorCount <= 0 {
		return ErrInvalidBasicMsg(
			"MsgRequestData: Sufficient validator count (%d) must be positive.",
			msg.SufficientValidatorCount,
		)
	}
	if msg.RequestedValidatorCount < msg.SufficientValidatorCount {
		return ErrInvalidBasicMsg(
			"MsgRequestData: Request validator count (%d) must not be less than sufficient validator count (%d).",
			msg.RequestedValidatorCount,
			msg.SufficientValidatorCount,
		)
	}
	if msg.Expiration <= 0 {
		return ErrInvalidBasicMsg("MsgRequestData: Expiration period (%d) must be positive.",
			msg.Expiration,
		)
	}
	if msg.PrepareGas <= 0 {
		return ErrInvalidBasicMsg("MsgRequestData: Prepare gas (%d) must be positive.",
			msg.PrepareGas,
		)
	}
	if msg.ExecuteGas <= 0 {
		return ErrInvalidBasicMsg("MsgRequestData: Execute gas (%d) must be positive.",
			msg.ExecuteGas,
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

// MsgReportData is a message sent by each of the block validators to respond to a data request.
type MsgReportData struct {
	RequestID int64           `json:"requestID"`
	DataSet   []RawDataReport `json:"dataSet"`
	Sender    sdk.ValAddress  `json:"sender"`
}

// NewMsgReportData creates a new MsgReportData instance.
func NewMsgReportData(
	requestID int64,
	dataSet []RawDataReport,
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
		return ErrInvalidBasicMsg("MsgReportData: Request id (%d) must be positive.", msg.RequestID)
	}
	if msg.DataSet == nil || len(msg.DataSet) == 0 {
		return ErrInvalidBasicMsg("MsgReportData: Data set must not be empty.")
	}
	if msg.Sender.Empty() {
		return ErrInvalidBasicMsg("MsgReportData: Sender address must not be empty.")
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

// MsgCreateDataSource is a message for creating a new data source.
type MsgCreateDataSource struct {
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Fee        sdk.Coins      `json:"fee"`
	Executable []byte         `json:"executable"`
	Sender     sdk.AccAddress `json:"sender"`
}

// NewMsgCreateDataSource creates a new MsgCreateDataSource instance.
func NewMsgCreateDataSource(
	owner sdk.AccAddress,
	name string,
	fee sdk.Coins,
	executable []byte,
	sender sdk.AccAddress,
) MsgCreateDataSource {
	return MsgCreateDataSource{
		Owner:      owner,
		Name:       name,
		Fee:        fee,
		Executable: executable,
		Sender:     sender,
	}
}

// Route implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) Type() string { return "create_data_source" }

// ValidateBasic implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return ErrInvalidBasicMsg("MsgCreateDataSource: Owner address must not be empty.")
	}
	if msg.Name == "" {
		return ErrInvalidBasicMsg("MsgCreateDataSource: Name must not be empty.")
	}
	if !msg.Fee.IsValid() {
		return ErrInvalidBasicMsg("MsgCreateDataSource: Fee must be valid (%s)", msg.Fee.String())
	}
	if msg.Executable == nil || len(msg.Executable) == 0 {
		return ErrInvalidBasicMsg("MsgCreateDataSource: Executable must not be empty.")
	}
	if msg.Sender.Empty() {
		return ErrInvalidBasicMsg("MsgCreateDataSource: Sender address must not be empty.")
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgEditDataSource is a message for editing an existing data source.
type MsgEditDataSource struct {
	DataSourceID int64          `json:"dataSourceID"`
	Owner        sdk.AccAddress `json:"owner"`
	Name         string         `json:"name"`
	Fee          sdk.Coins      `json:"fee"`
	Executable   []byte         `json:"executable"`
	Sender       sdk.AccAddress `json:"sender"`
}

// NewMsgEditDataSource creates a new MsgEditDataSource instance.
func NewMsgEditDataSource(
	dataSourceID int64,
	owner sdk.AccAddress,
	name string,
	fee sdk.Coins,
	executable []byte,
	sender sdk.AccAddress,
) MsgEditDataSource {
	return MsgEditDataSource{
		DataSourceID: dataSourceID,
		Owner:        owner,
		Name:         name,
		Fee:          fee,
		Executable:   executable,
		Sender:       sender,
	}
}

// Route implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) Type() string { return "edit_data_source" }

// ValidateBasic implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) ValidateBasic() sdk.Error {
	if msg.DataSourceID <= 0 {
		return ErrInvalidBasicMsg("MsgEditDataSource: Data source id (%d) must be positive.", msg.DataSourceID)
	}
	if msg.Owner.Empty() {
		return ErrInvalidBasicMsg("MsgEditDataSource: Owner address must not be empty.")
	}
	if msg.Name == "" {
		return ErrInvalidBasicMsg("MsgEditDataSource: Name must not be empty.")
	}
	if !msg.Fee.IsValid() {
		return ErrInvalidBasicMsg("MsgEditDataSource: Fee must be valid (%s)", msg.Fee.String())
	}
	if msg.Executable == nil || len(msg.Executable) == 0 {
		return ErrInvalidBasicMsg("MsgEditDataSource: Executable must not be empty.")
	}
	if msg.Sender.Empty() {
		return ErrInvalidBasicMsg("MsgEditDataSource: Sender address must not be empty.")
	}

	return nil
}

// GetSigners implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgCreateOracleScript is a message for creating an oracle script.
type MsgCreateOracleScript struct {
	Owner  sdk.AccAddress `json:"owner"`
	Name   string         `json:"name"`
	Code   []byte         `json:"code"`
	Sender sdk.AccAddress `json:"sender"`
}

// NewMsgCreateOracleScript creates a new MsgCreateOracleScript instance.
func NewMsgCreateOracleScript(
	owner sdk.AccAddress,
	name string,
	code []byte,
	sender sdk.AccAddress,
) MsgCreateOracleScript {
	return MsgCreateOracleScript{
		Owner:  owner,
		Name:   name,
		Code:   code,
		Sender: sender,
	}
}

// Route implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) Type() string { return "create_oracle_script" }

// ValidateBasic implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return ErrInvalidBasicMsg("MsgCreateOracleScript: Owner address must not be empty.")
	}
	if msg.Sender.Empty() {
		return ErrInvalidBasicMsg("MsgCreateOracleScript: Sender address must not be empty.")
	}
	if msg.Name == "" {
		return ErrInvalidBasicMsg("MsgCreateOracleScript: Name must not be empty.")
	}
	if msg.Code == nil || len(msg.Code) == 0 {
		return ErrInvalidBasicMsg("MsgCreateOracleScript: Code must not be empty.")
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgEditOracleScript is a message for editing an existing oracle script.
type MsgEditOracleScript struct {
	OracleScriptID int64          `json:"oracleScriptID"`
	Owner          sdk.AccAddress `json:"owner"`
	Name           string         `json:"name"`
	Code           []byte         `json:"code"`
	Sender         sdk.AccAddress `json:"sender"`
}

// NewMsgEditOracleScript creates a new MsgEditOracleScript instance.
func NewMsgEditOracleScript(
	oracleScriptID int64,
	owner sdk.AccAddress,
	name string,
	code []byte,
	sender sdk.AccAddress,
) MsgEditOracleScript {
	return MsgEditOracleScript{
		OracleScriptID: oracleScriptID,
		Owner:          owner,
		Name:           name,
		Code:           code,
		Sender:         sender,
	}
}

// Route implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) Type() string { return "edit_oracle_script" }

// ValidateBasic implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) ValidateBasic() sdk.Error {
	if msg.OracleScriptID <= 0 {
		return ErrInvalidBasicMsg("MsgEditOracleScript: Oracle script id (%d) must be positive.", msg.OracleScriptID)
	}
	if msg.Owner.Empty() {
		return ErrInvalidBasicMsg("MsgEditOracleScript: Owner address must not be empty.")
	}
	if msg.Sender.Empty() {
		return ErrInvalidBasicMsg("MsgEditOracleScript: Sender address must not be empty.")
	}
	if msg.Name == "" {
		return ErrInvalidBasicMsg("MsgEditOracleScript: Name must not be empty.")
	}
	if msg.Code == nil || len(msg.Code) == 0 {
		return ErrInvalidBasicMsg("MsgEditOracleScript: Code must not be empty.")
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
