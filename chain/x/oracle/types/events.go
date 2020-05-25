package types

// Event types
const (
	EventTypeCreateDataSource   = "create_data_source"
	EventTypeEditDataSource     = "edit_data_source"
	EventTypeCreateOracleScript = "create_oracle_script"
	EventTypeEditOracleScript   = "edit_oracle_script"
	EventTypeRequest            = "request"
	EventTypeRawRequest         = "raw_request"
	EventTypeReport             = "report"
	EventTypeAddReporter        = "add_reporter"
	EventTypeRemoveReporter     = "remove_reporter"
	EventTypeRequestExecute     = "request_execute"

	AttributeKeyID             = "id"
	AttributeKeyRequestID      = "request_id"
	AttributeKeyDataSourceID   = "data_source_id"
	AttributeKeyOracleScriptID = "oracle_script_id"
	AttributeKeyExternalID     = "external_id"
	AttributeKeyCalldata       = "calldata"
	AttributeKeyValidator      = "validator"
	AttributeKeyReporter       = "reporter"
	AttributeKeyClientID       = "client_id"
	AttributeKeyAskCount       = "ask_count"
	AttributeKeyMinCount       = "min_count"
	AttributeKeyAnsCount       = "ans_count"
	AttributeKeyRequestTime    = "request_time"
	AttributeKeyResolveTime    = "resolve_time"
	AttributeKeyResolveStatus  = "resolve_status"
	AttributeKeyResult         = "result"
	AttributeKeyResultHash     = "result_hash"
	AttributeKeyDataSourceHash = "datasource_hash"
)
