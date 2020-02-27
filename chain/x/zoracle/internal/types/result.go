package types

// Result is a data structure that stores the detail of a result of a specific request.
type Result struct {
	RequestTime              int64  `json:"requestTime"`
	AggregationTime          int64  `json:"aggregationTime"`
	RequestedValidatorsCount int64  `json:"requestedValidatorsCount"`
	SufficientValidatorCount int64  `json:"sufficientValidatorCount"`
	ReportedValidatorsCount  int64  `json:"reportedValidatorsCount"`
	Data                     []byte `json:"data"`
}

// NewResult creates a new Result instance.
func NewResult(
	requestTime int64,
	aggregationTime int64,
	requestedValidatorsCount int64,
	sufficientValidatorCount int64,
	reportedValidatorsCount int64,
	data []byte,
) Result {
	return Result{
		RequestTime:              requestTime,
		AggregationTime:          aggregationTime,
		RequestedValidatorsCount: requestedValidatorsCount,
		SufficientValidatorCount: sufficientValidatorCount,
		ReportedValidatorsCount:  reportedValidatorsCount,
		Data:                     data,
	}
}
