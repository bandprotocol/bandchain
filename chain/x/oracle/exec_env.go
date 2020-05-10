package oracle

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExecEnv struct {
	request            types.Request
	now                int64
	maxResultSize      int64
	maxCalldataSize    int64
	maxRawRequestCount int64
	rawRequests        []types.RawRequest
	reports            map[string]map[types.ExternalID]types.RawReport
}

func NewExecEnv(ctx sdk.Context, k Keeper, req types.Request) *ExecEnv {
	return &ExecEnv{
		request:            req,
		now:                ctx.BlockTime().Unix(),
		maxResultSize:      int64(k.GetParam(ctx, KeyMaxResultSize)),
		maxCalldataSize:    types.MaxCalldataSize,
		maxRawRequestCount: int64(k.GetParam(ctx, KeyMaxRawRequestCount)),
		rawRequests:        []types.RawRequest{},
		reports:            make(map[string]map[types.ExternalID]types.RawReport),
	}
}

// GetRawRequests returns the list of raw requests made during Owasm prepare run.
func (env *ExecEnv) GetRawRequests() []types.RawRequest {
	return env.rawRequests
}

// SetReports loads the reports to the environment. Must be called prior to Owasm execute run.
func (env *ExecEnv) SetReports(reports []types.Report) {
	for _, report := range reports {
		valReports := make(map[types.ExternalID]types.RawReport)
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
	return env.maxResultSize
}

// GetMaximumCalldataOfDataSourceSize implements Owasm ExecEnv interface.
func (env *ExecEnv) GetMaximumCalldataOfDataSourceSize() int64 {
	return env.maxCalldataSize
}

// RequestedValidators implements Owasm ExecEnv interface.
func (env *ExecEnv) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, types.ErrValidatorOutOfRange
	}
	return env.request.RequestedValidators[validatorIndex], nil
}

// RequestExternalData implements Owasm ExecEnv interface.
func (env *ExecEnv) RequestExternalData(did int64, eid int64, calldata []byte) error {
	if int64(len(calldata)) > env.maxCalldataSize {
		return types.ErrValidatorOutOfRange
	}
	if int64(len(env.rawRequests)) >= env.maxRawRequestCount {
		return types.ErrTooManyRawRequests
	}
	env.rawRequests = append(env.rawRequests, types.NewRawRequest(
		types.ExternalID(eid), types.DataSourceID(did), calldata,
	))
	return nil
}

// GetExternalData implements Owasm ExecEnv interface.
func (env *ExecEnv) GetExternalData(eid int64, valIdx int64) ([]byte, uint32, error) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, types.ErrValidatorOutOfRange
	}
	valAddr := env.request.RequestedValidators[valIdx].String()
	valReports, ok := env.reports[valAddr]
	if !ok {
		return nil, 0, types.ErrItemNotFound
	}
	valReport, ok := valReports[types.ExternalID(eid)]
	if !ok {
		return nil, 0, types.ErrItemNotFound
	}
	return valReport.Data, valReport.ExitCode, nil
}
