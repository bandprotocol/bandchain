package types

// Event types
const (
	EventTypeCreateDataSource    = "create_data_source"
	EventTypeEditDataSource      = "edit_data_source"
	EventTypeCreateOracleScript  = "create_oracle_script"
	EventTypeEditOracleScript    = "edit_oracle_script"
	EventTypeRawRequest          = "raw_request"
	EventTypeRequest             = "request"
	EventTypeReport              = "report"
	EventTypeAddOracleAddress    = "add_oracle_address"
	EventTypeRemoveOracleAddress = "remove_oracle_address"
	EventTypeRequestExecute      = "request_execute"

	AttributeKeyID            = "id"
	AttributeKeyRequestID     = "request_id"
	AttributeKeyDataSourceID  = "data_source_id"
	AttributeKeyExternalID    = "external_id"
	AttributeKeyCalldata      = "calldata"
	AttributeKeyValidator     = "validator"
	AttributeKeyReporter      = "reporter"
	AttributeKeyResolveStatus = "resolve_status"
	AttributeKeyResult        = "result"
	AttributeKeyRequestTime   = "request_time"
	AttributeKeyResolvedTime  = "resolved_time"
)
