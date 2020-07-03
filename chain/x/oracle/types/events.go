package types

// nolint
const (
	EventTypeCreateDataSource   = "create_data_source"
	EventTypeEditDataSource     = "edit_data_source"
	EventTypeCreateOracleScript = "create_oracle_script"
	EventTypeEditOracleScript   = "edit_oracle_script"
	EventTypeRequest            = "request"
	EventTypeRawRequest         = "raw_request"
	EventTypeReport             = "report"
	EventTypeActivate           = "activate"
	EventTypeDeactivate         = "deactivate"
	EventTypeAddReporter        = "add_reporter"
	EventTypeRemoveReporter     = "remove_reporter"
	EventTypeResolve            = "resolve"

	AttributeKeyID             = "id"
	AttributeKeyDataSourceID   = "data_source_id"
	AttributeKeyOracleScriptID = "oracle_script_id"
	AttributeKeyExternalID     = "external_id"
	AttributeKeyDataSourceHash = "data_source_hash"
	AttributeKeyCalldata       = "calldata"
	AttributeKeyValidator      = "validator"
	AttributeKeyReporter       = "reporter"
	AttributeKeyClientID       = "client_id"
	AttributeKeyAskCount       = "ask_count"
	AttributeKeyMinCount       = "min_count"
	AttributeKeyResolveStatus  = "resolve_status"
	AttributeKeyResult         = "result"
	AttributeKeyReason         = "reason"
)
