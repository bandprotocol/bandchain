package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is the name of the oracle module
const RouterKey = ModuleName

// Route implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) Type() string { return "request" }

// ValidateBasic implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if len(msg.Calldata) > MaxDataSize {
		return WrapMaxError(ErrTooLargeCalldata, len(msg.Calldata), MaxDataSize)
	}
	if msg.MinCount <= 0 {
		return sdkerrors.Wrapf(ErrInvalidMinCount, "got: %d", msg.MinCount)
	}
	if msg.AskCount < msg.MinCount {
		return sdkerrors.Wrapf(ErrInvalidAskCount, "got: %d, min count: %d", msg.AskCount, msg.MinCount)
	}
	if len(msg.ClientID) > MaxClientIDLength {
		return WrapMaxError(ErrTooLongClientID, len(msg.ClientID), MaxClientIDLength)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgRequestData.
func (msg MsgRequestData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) Type() string { return "report" }

// ValidateBasic implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s", msg.Validator)
	}
	if err := sdk.VerifyAddressFormat(msg.Reporter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "reporter: %s", msg.Reporter)
	}
	if len(msg.RawReports) == 0 {
		return ErrEmptyReport
	}
	uniqueMap := make(map[ExternalID]bool)
	for _, r := range msg.RawReports {
		if _, found := uniqueMap[r.ExternalID]; found {
			return sdkerrors.Wrapf(ErrDuplicateExternalID, "external id: %d", r.ExternalID)
		}
		uniqueMap[r.ExternalID] = true
		if len(r.Data) > MaxDataSize {
			return WrapMaxError(ErrTooLargeRawReportData, len(r.Data), MaxDataSize)
		}
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Reporter}
}

// GetSignBytes implements the sdk.Msg interface for MsgReportData.
func (msg MsgReportData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) Type() string { return "create_data_source" }

// ValidateBasic implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "owner: %s", msg.Owner)
	}
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if len(msg.Name) > MaxNameLength {
		return WrapMaxError(ErrTooLongName, len(msg.Name), MaxNameLength)
	}
	if len(msg.Description) > MaxDescriptionLength {
		return WrapMaxError(ErrTooLongDescription, len(msg.Description), MaxDescriptionLength)
	}
	if len(msg.Executable) == 0 {
		return ErrEmptyExecutable
	}
	if len(msg.Executable) > MaxExecutableSize {
		return WrapMaxError(ErrTooLargeExecutable, len(msg.Executable), MaxExecutableSize)
	}
	if bytes.Equal(msg.Executable, DoNotModifyBytes) {
		return ErrCreateWithDoNotModify
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgCreateDataSource.
func (msg MsgCreateDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) Type() string { return "edit_data_source" }

// ValidateBasic implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "owner: %s", msg.Owner)
	}
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if len(msg.Name) > MaxNameLength {
		return WrapMaxError(ErrTooLongName, len(msg.Name), MaxNameLength)
	}
	if len(msg.Description) > MaxDescriptionLength {
		return WrapMaxError(ErrTooLongDescription, len(msg.Description), MaxDescriptionLength)
	}
	if len(msg.Executable) == 0 {
		return ErrEmptyExecutable
	}
	if len(msg.Executable) > MaxExecutableSize {
		return WrapMaxError(ErrTooLargeExecutable, len(msg.Executable), MaxExecutableSize)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgEditDataSource.
func (msg MsgEditDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) Type() string { return "create_oracle_script" }

// ValidateBasic implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "owner: %s", msg.Owner)
	}
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if len(msg.Name) > MaxNameLength {
		return WrapMaxError(ErrTooLongName, len(msg.Name), MaxNameLength)
	}
	if len(msg.Description) > MaxDescriptionLength {
		return WrapMaxError(ErrTooLongDescription, len(msg.Description), MaxDescriptionLength)
	}
	if len(msg.Schema) > MaxSchemaLength {
		return WrapMaxError(ErrTooLongSchema, len(msg.Schema), MaxSchemaLength)
	}
	if len(msg.SourceCodeURL) > MaxURLLength {
		return WrapMaxError(ErrTooLongURL, len(msg.SourceCodeURL), MaxURLLength)
	}
	if len(msg.Code) == 0 {
		return ErrEmptyWasmCode
	}
	if len(msg.Code) > MaxWasmCodeSize {
		return WrapMaxError(ErrTooLargeWasmCode, len(msg.Code), MaxWasmCodeSize)
	}
	if bytes.Equal(msg.Code, DoNotModifyBytes) {
		return ErrCreateWithDoNotModify
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgCreateOracleScript.
func (msg MsgCreateOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) Type() string { return "edit_oracle_script" }

// ValidateBasic implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "owner: %s", msg.Owner)
	}
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s", msg.Sender)
	}
	if len(msg.Name) > MaxNameLength {
		return WrapMaxError(ErrTooLongName, len(msg.Name), MaxNameLength)
	}
	if len(msg.Description) > MaxDescriptionLength {
		return WrapMaxError(ErrTooLongDescription, len(msg.Description), MaxDescriptionLength)
	}
	if len(msg.Schema) > MaxSchemaLength {
		return WrapMaxError(ErrTooLongSchema, len(msg.Schema), MaxSchemaLength)
	}
	if len(msg.SourceCodeURL) > MaxURLLength {
		return WrapMaxError(ErrTooLongURL, len(msg.SourceCodeURL), MaxURLLength)
	}
	if len(msg.Code) == 0 {
		return ErrEmptyWasmCode
	}
	if len(msg.Code) > MaxWasmCodeSize {
		return WrapMaxError(ErrTooLargeWasmCode, len(msg.Code), MaxWasmCodeSize)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes implements the sdk.Msg interface for MsgEditOracleScript.
func (msg MsgEditOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgActivate.
func (msg MsgActivate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgActivate.
func (msg MsgActivate) Type() string { return "activate" }

// ValidateBasic implements the sdk.Msg interface for MsgActivate.
func (msg MsgActivate) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s", msg.Validator)
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgActivate.
func (msg MsgActivate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes implements the sdk.Msg interface for MsgActivate.
func (msg MsgActivate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgAddReporter.
func (msg MsgAddReporter) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgAddReporter.
func (msg MsgAddReporter) Type() string { return "add_reporter" }

// ValidateBasic implements the sdk.Msg interface for MsgAddReporter.
func (msg MsgAddReporter) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s", msg.Validator)
	}
	if err := sdk.VerifyAddressFormat(msg.Reporter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "reporter: %s", msg.Reporter)
	}
	if sdk.ValAddress(msg.Reporter).Equals(msg.Validator) {
		return ErrSelfReferenceAsReporter
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgAddReporter.
func (msg MsgAddReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes implements the sdk.Msg interface for MsgAddReporter.
func (msg MsgAddReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements the sdk.Msg interface for MsgRemoveReporter.
func (msg MsgRemoveReporter) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgRemoveReporter.
func (msg MsgRemoveReporter) Type() string { return "remove_reporter" }

// ValidateBasic implements the sdk.Msg interface for MsgRemoveReporter.
func (msg MsgRemoveReporter) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s", msg.Validator)
	}
	if err := sdk.VerifyAddressFormat(msg.Reporter); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "reporter: %s", msg.Reporter)
	}
	if sdk.ValAddress(msg.Reporter).Equals(msg.Validator) {
		return ErrSelfReferenceAsReporter
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgRemoveReporter.
func (msg MsgRemoveReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes implements the sdk.Msg interface for MsgRemoveReporter.
func (msg MsgRemoveReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
