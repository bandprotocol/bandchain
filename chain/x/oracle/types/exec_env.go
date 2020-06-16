package types

import (
	"github.com/bandprotocol/bandchain/go-owasm/api"
)

// BaseEnv combines shared functions used in prepare and execution Owasm program,
type BaseEnv struct {
	request Request
}

// GetCalldata implements Owasm ExecEnv interface.
func (env *BaseEnv) GetCalldata() []byte {
	return env.request.Calldata
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *BaseEnv) SetReturnData(data []byte) error {
	return api.ErrSetReturnDataWrongPeriod
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
func (env *BaseEnv) GetAnsCount() (int64, error) {
	return 0, api.ErrAnsCountWrongPeriod
}

// AskExternalData implements Owasm ExecEnv interface.
func (env *BaseEnv) AskExternalData(eid int64, did int64, data []byte) error {
	return api.ErrAskExternalDataWrongPeriod
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *PrepareEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	return 0, api.ErrGetExternalDataStatusWrongPeriod
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *PrepareEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	return nil, api.ErrGetExternalDataWrongPeriod
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
func (env *PrepareEnv) AskExternalData(eid int64, did int64, data []byte) error {
	if int64(len(data)) > MaxRawRequestDataSize {
		return api.ErrSpanExceededCapacity
	}
	if int64(len(env.rawRequests)) >= env.maxRawRequests {
		return api.ErrAskExternalDataExceed
	}
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		ExternalID(eid), DataSourceID(did), data,
	))
	return nil
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
func (env *ExecuteEnv) GetAnsCount() (int64, error) {
	return int64(len(env.reports)), nil
}

// SetReturnData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) SetReturnData(data []byte) error {
	env.Retdata = data
	return nil
}

func (env *ExecuteEnv) getExternalDataFull(eid int64, valIdx int64) ([]byte, int64, error) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, api.ErrValidatorOutOfRange
	}
	valAddr := env.request.RequestedValidators[valIdx].String()
	valReports, ok := env.reports[valAddr]
	if !ok {
		return nil, -1, nil
	}
	valReport, ok := valReports[ExternalID(eid)]
	if !ok {
		return nil, -1, api.ErrInvalidExternalID
	}
	return valReport.Data, int64(valReport.ExitCode), nil
}

// GetExternalDataStatus implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	_, status, err := env.getExternalDataFull(eid, vid)
	return status, err
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *ExecuteEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	data, _, err := env.getExternalDataFull(eid, vid)
	return data, err
}
