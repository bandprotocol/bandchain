package types

// BaseEnv combines shared functions used in prepare and execution Owasm program,
type BaseEnv struct {
	request Request
}

// GetCalldata implements Owasm ExecEnv interface.
func (env *BaseEnv) GetCalldata() []byte {
	return env.request.Calldata
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *BaseEnv) SetReturnData(data []byte) {
	panic("Cannot set return data on non-execution period")
}

// GetAskCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetAskCount() int64 {
	return int64(len(env.request.RequestedValidators))
}

// GetMinCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetMinCount() int64 {
	return int64(env.request.MinCount)
}

// GetAnsCount implements Owasm ExecEnv interface.
func (env *BaseEnv) GetAnsCount() int64 {
	panic("Cannot get ans count on non-execution period")
}

// AskExternalData implements Owasm ExecEnv interface.
func (env *BaseEnv) AskExternalData(eid int64, did int64, data []byte) {
	panic("Cannot ask external data on non-prepare period")
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *PrepareEnv) GetExternalDataStatus(eid int64, vid int64) int64 {
	panic("Cannot get external data status on non-execution period")
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *PrepareEnv) GetExternalData(eid int64, vid int64) []byte {
	panic("Cannot get external data on non-execution period")
}

// PrepareEnv implements ExecEnv interface only expected function and panic on non-prepare functions.
type PrepareEnv struct {
	BaseEnv
	maxRawRequests int64
	rawRequests    []RawRequest
}

// NewPrepareEnv creates a new environment instance for prepare period.
func NewPrepareEnv(req Request, maxRawRequests int64) *PrepareEnv {
	return &PrepareEnv{
		BaseEnv: BaseEnv{
			request: req,
		},
		maxRawRequests: maxRawRequests,
	}
}

// AskExternalData implements Owasm ExecEnv interface.
func (env *PrepareEnv) AskExternalData(eid int64, did int64, data []byte) {
	if int64(len(data)) > MaxRawRequestDataSize {
		return
	}
	if int64(len(env.rawRequests)) >= env.maxRawRequests {
		return
	}
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		ExternalID(eid), DataSourceID(did), data,
	))
}

// GetRawRequests returns the list of raw requests made during Owasm prepare run.
func (env *PrepareEnv) GetRawRequests() []RawRequest {
	return env.rawRequests
}

// ExecuteEnv implements ExecEnv interface only expected function and panic on prepare related functions.
type ExecuteEnv struct {
	BaseEnv
	reports map[string]map[ExternalID]RawReport
	Retdata []byte
}

// NewExecuteEnv creates a new environment instance for execution period.
func NewExecuteEnv(req Request) *ExecuteEnv {
	return &ExecuteEnv{
		BaseEnv: BaseEnv{
			request: req,
		},
		reports: make(map[string]map[ExternalID]RawReport),
	}
}

// SetReports loads the reports to the environment.
func (env *ExecuteEnv) SetReports(reports []Report) {
	for _, report := range reports {
		valReports := make(map[ExternalID]RawReport)
		for _, each := range report.RawReports {
			valReports[each.ExternalID] = each
		}
		env.reports[report.Validator.String()] = valReports
	}
}

// GetAnsCount implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetAnsCount() int64 {
	return int64(len(env.reports))
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) SetReturnData(data []byte) {
	env.Retdata = data
}

func (env *ExecuteEnv) getExternalDataFull(eid int64, valIdx int64) ([]byte, int64) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, -1
	}
	valAddr := env.request.RequestedValidators[valIdx].String()
	valReports, ok := env.reports[valAddr]
	if !ok {
		return nil, -1
	}
	valReport, ok := valReports[ExternalID(eid)]
	if !ok {
		return nil, -1
	}
	return valReport.Data, int64(valReport.ExitCode)
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalDataStatus(eid int64, vid int64) int64 {
	_, status := env.getExternalDataFull(eid, vid)
	return status
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalData(eid int64, vid int64) []byte {
	data, _ := env.getExternalDataFull(eid, vid)
	return data
}
