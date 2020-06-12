package types

// Query endpoints supported by the oracle Querier.
const (
	QueryParams        = "params"
	QueryCounts        = "counts"
	QueryData          = "data"
	QueryDataSources   = "data_sources"
	QueryOracleScripts = "oracle_scripts"
	QueryRequests      = "requests"
	QueryReporters     = "reporters"
	QueryReportInfo    = "report_info"
)

// QueryCountsResult is the struct for the result of query counts.
type QueryCountsResult struct {
	DataSourceCount   int64 `json:"data_source_count"`
	OracleScriptCount int64 `json:"oracle_script_count"`
	RequestCount      int64 `json:"request_count"`
}

// QueryRequestResult is the struct for the result of request query.
type QueryRequestResult struct {
	Request Request  `json:"request"`
	Reports []Report `json:"reports"`
	Result  *Result  `json:"result"`
}
