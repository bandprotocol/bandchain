package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is the name of the oracle module
const RouterKey = ModuleName

// Route returns the route of MsgRequestData - "oracle" (sdk.Msg interface).
func (msg MsgRequestData) Route() string { return RouterKey }

// Type returns the message type of MsgRequestData (sdk.Msg interface).
func (msg MsgRequestData) Type() string { return "request" }

// ValidateBasic checks whether the given MsgRequestData instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgRequestData (sdk.Msg interface).
func (msg MsgRequestData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgRequestData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgReportData - "oracle" (sdk.Msg interface).
func (msg MsgReportData) Route() string { return RouterKey }

// Type returns the message type of MsgReportData (sdk.Msg interface).
func (msg MsgReportData) Type() string { return "report" }

// ValidateBasic checks whether the given MsgReportData instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgReportData (sdk.Msg interface).
func (msg MsgReportData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Reporter}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgReportData) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgCreateDataSource - "oracle" (sdk.Msg interface).
func (msg MsgCreateDataSource) Route() string { return RouterKey }

// Type returns the message type of MsgCreateDataSource (sdk.Msg interface).
func (msg MsgCreateDataSource) Type() string { return "create_data_source" }

// ValidateBasic checks whether the given MsgCreateDataSource instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgCreateDataSource (sdk.Msg interface).
func (msg MsgCreateDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgCreateDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgEditDataSource - "oracle" (sdk.Msg interface).
func (msg MsgEditDataSource) Route() string { return RouterKey }

// Type returns the message type of MsgEditDataSource (sdk.Msg interface).
func (msg MsgEditDataSource) Type() string { return "edit_data_source" }

// ValidateBasic checks whether the given MsgEditDataSource instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgEditDataSource (sdk.Msg interface).
func (msg MsgEditDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgEditDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgCreateOracleScript - "oracle" (sdk.Msg interface).
func (msg MsgCreateOracleScript) Route() string { return RouterKey }

// Type returns the message type of MsgCreateOracleScript (sdk.Msg interface).
func (msg MsgCreateOracleScript) Type() string { return "create_oracle_script" }

// ValidateBasic checks whether the given MsgCreateOracleScript instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgCreateOracleScript (sdk.Msg interface).
func (msg MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgCreateOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgEditOracleScript - "oracle" (sdk.Msg interface).
func (msg MsgEditOracleScript) Route() string { return RouterKey }

// Type returns the message type of MsgEditOracleScript (sdk.Msg interface).
func (msg MsgEditOracleScript) Type() string { return "edit_oracle_script" }

// ValidateBasic checks whether the given MsgEditOracleScript instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgEditOracleScript (sdk.Msg interface).
func (msg MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgEditOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgActivate - "oracle" (sdk.Msg interface).
func (msg MsgActivate) Route() string { return RouterKey }

// Type returns the message type of MsgActivate (sdk.Msg interface).
func (msg MsgActivate) Type() string { return "activate" }

// ValidateBasic checks whether the given MsgActivate instance (sdk.Msg interface).
func (msg MsgActivate) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s", msg.Validator)
	}
	return nil
}

// GetSigners returns the required signers for the given MsgActivate (sdk.Msg interface).
func (msg MsgActivate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgActivate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgAddReporter - "oracle" (sdk.Msg interface).
func (msg MsgAddReporter) Route() string { return RouterKey }

// Type returns the message type of MsgAddReporter (sdk.Msg interface).
func (msg MsgAddReporter) Type() string { return "add_reporter" }

// ValidateBasic checks whether the given MsgAddReporter instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgAddReporter (sdk.Msg interface).
func (msg MsgAddReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgAddReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route returns the route of MsgRemoveReporter - "oracle" (sdk.Msg interface).
func (msg MsgRemoveReporter) Route() string { return RouterKey }

// Type returns the message type of MsgRemoveReporter (sdk.Msg interface).
func (msg MsgRemoveReporter) Type() string { return "remove_reporter" }

// ValidateBasic checks whether the given MsgRemoveReporter instance (sdk.Msg interface).
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

// GetSigners returns the required signers for the given MsgRemoveReporter (sdk.Msg interface).
func (msg MsgRemoveReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Validator)}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (msg MsgRemoveReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
