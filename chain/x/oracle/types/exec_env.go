package types

type ExecEnv struct {
	request            Request
	now                int64
	maxRawRequestCount int64
	rawRequests        []RawRequest
	reports            map[string]map[ExternalID]RawReport
}

func NewExecEnv(req Request, now, maxRawRequestCount int64) *ExecEnv {
	return &ExecEnv{
		request:            req,
		now:                now,
		maxRawRequestCount: maxRawRequestCount,
		rawRequests:        []RawRequest{},
		reports:            make(map[string]map[ExternalID]RawReport),
	}
}

// GetRawRequests returns the list of raw requests made during Owasm prepare run.
func (env *ExecEnv) GetRawRequests() []RawRequest {
	return env.rawRequests
}

// SetReports loads the reports to the environment. Must be called prior to Owasm execute run.
func (env *ExecEnv) SetReports(reports []Report) {
	for _, report := range reports {
		valReports := make(map[ExternalID]RawReport)
		for _, each := range report.RawReports {
			valReports[each.ExternalID] = each
		}
		env.reports[report.Validator.String()] = valReports
	}
}

// GetAskCount implements Owasm ExecEnv interface.
func (env *ExecEnv) GetAskCount() int64 {
	return int64(len(env.request.RequestedValidators))
}

// GetMinCount implements Owasm ExecEnv interface.
func (env *ExecEnv) GetMinCount() int64 {
	return env.request.MinCount
}

// GetAnsCount implements Owasm ExecEnv interface.
func (env *ExecEnv) GetAnsCount() int64 {
	return int64(len(env.reports))
}

// GetPrepareBlockTime implements Owasm ExecEnv interface.
func (env *ExecEnv) GetPrepareBlockTime() int64 {
	return env.request.RequestTime
}

// GetAggregateBlockTime implements Owasm ExecEnv interface.
func (env *ExecEnv) GetAggregateBlockTime() int64 {
	if len(env.reports) == 0 { // Size of reports must be zero during prepare.
		return 0
	}
	return env.now
}

// GetMaximumResultSize implements Owasm ExecEnv interface.
func (env *ExecEnv) GetMaximumResultSize() int64 {
	return MaxResultSize
}

// GetMaximumCalldataOfDataSourceSize implements Owasm ExecEnv interface.
func (env *ExecEnv) GetMaximumCalldataOfDataSourceSize() int64 {
	return MaxRawRequestDataSize
}

// RequestedValidators implements Owasm ExecEnv interface.
func (env *ExecEnv) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, ErrItemNotFound
	}
	return env.request.RequestedValidators[validatorIndex], nil
}

// RequestExternalData implements Owasm ExecEnv interface.
func (env *ExecEnv) RequestExternalData(did int64, eid int64, calldata []byte) error {
	if int64(len(calldata)) > MaxRawRequestDataSize {
		return ErrTooLargeCalldata
	}
	if int64(len(env.rawRequests)) >= env.maxRawRequestCount {
		return ErrTooManyRawRequests
	}
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		ExternalID(eid), DataSourceID(did), calldata,
	))
	return nil
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *ExecEnv) GetExternalData(eid int64, valIdx int64) ([]byte, uint32, error) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, ErrItemNotFound
	}
	valAddr := env.request.RequestedValidators[valIdx].String()
	valReports, ok := env.reports[valAddr]
	if !ok {
		return nil, 0, ErrItemNotFound
	}
	valReport, ok := valReports[ExternalID(eid)]
	if !ok {
		return nil, 0, ErrItemNotFound
	}
	return valReport.Data, valReport.ExitCode, nil
}
