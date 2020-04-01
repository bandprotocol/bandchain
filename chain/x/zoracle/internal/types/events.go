package types

// Event types
const (
	EventTypeCreateDataSource    = "create_data_source"
	EventTypeEditDataSource      = "edit_data_source"
	EventTypeCreateOracleScript  = "create_oracle_script"
	EventTypeEditOracleScript    = "edit_oracle_script"
	EventTypeRequest             = "request"
	EventTypeReport              = "report"
	EventTypeAddOracleAddress    = "add_oracle_address"
	EventTypeRemoveOracleAddress = "remove_oracle_address"
	EventTypeRequestExecute      = "request_execute"

	AttributeKeyID            = "id"
	AttributeKeyRequestID     = "request_id"
	AttributeKeyValidator     = "validator"
	AttributeKeyReporter      = "reporter"
	AttributeKeyResolveStatus = "resolve_status"
)
