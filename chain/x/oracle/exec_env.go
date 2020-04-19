package oracle

import (
	"errors"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExecutionEnvironment struct {
	request                types.Request
	now                    int64
	maxResultSize          int64
	maxCalldataSize        int64
	maxRawDataRequestCount int64
	rawDataRequests        []types.RawRequest
	rawDataReports         map[string]map[types.EID]types.RawDataReport
}

func NewExecutionEnvironment(ctx sdk.Context, k Keeper, req types.Request) *ExecutionEnvironment {
	return &ExecutionEnvironment{
		request:                req,
		now:                    ctx.BlockTime().Unix(),
		maxResultSize:          int64(k.GetParam(ctx, KeyMaxResultSize)),
		maxCalldataSize:        int64(k.GetParam(ctx, KeyMaxCalldataSize)),
		maxRawDataRequestCount: int64(k.GetParam(ctx, KeyMaxDataSourceCountPerRequest)),
		rawDataRequests:        []types.RawRequest{},
		rawDataReports:         make(map[string]map[types.EID]types.RawDataReport),
	}
}

func (env *ExecutionEnvironment) GetRawRequests() []types.RawRequest {
	return env.rawDataRequests
}

func (env *ExecutionEnvironment) SetReports(reports []types.Report) {
	for _, report := range reports {
		valReports := make(map[types.EID]types.RawDataReport)
		for _, each := range report.RawDataReports {
			valReports[each.ExternalID] = types.NewRawDataReport(each.ExitCode, each.Data)
		}
		env.rawDataReports[report.Validator.String()] = valReports
	}
}

func (env *ExecutionEnvironment) GetRequestedValidatorCount() int64 {
	return int64(len(env.request.RequestedValidators))
}

func (env *ExecutionEnvironment) GetSufficientValidatorCount() int64 {
	return env.request.SufficientValidatorCount
}

func (env *ExecutionEnvironment) GetReceivedValidatorCount() int64 {
	return int64(len(env.rawDataReports))
}

func (env *ExecutionEnvironment) GetPrepareBlockTime() int64 {
	return env.request.RequestTime
}

func (env *ExecutionEnvironment) GetMaximumResultSize() int64 {
	return env.maxResultSize
}

func (env *ExecutionEnvironment) GetMaximumCalldataOfDataSourceSize() int64 {
	return env.maxCalldataSize
}
func (env *ExecutionEnvironment) GetAggregateBlockTime() int64 {
	if len(env.rawDataReports) == 0 { // Size of reports must be zero during prepare.
		return 0
	}
	return env.now
}

func (env *ExecutionEnvironment) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, errors.New("validator out of range")
	}
	return env.request.RequestedValidators[validatorIndex], nil
}

func (env *ExecutionEnvironment) RequestExternalData(did int64, eid int64, calldata []byte) error {
	if int64(len(calldata)) > env.maxCalldataSize {
		return errors.New("calldata size limit exceeded")
	}
	if int64(len(env.rawDataRequests)) >= env.maxRawDataRequestCount {
		return errors.New("cannot request more than maxRawDataRequestCount")
	}

	env.rawDataRequests = append(env.rawDataRequests, types.NewRawRequest(
		types.ExternalID(eid), types.DataSourceID(did), calldata,
	))
	return nil
}

func (env *ExecutionEnvironment) GetExternalData(eid int64, valIdx int64) ([]byte, uint8, error) {
	if valIdx < 0 || valIdx >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, errors.New("validator out of range")
	}
	valAddr := env.request.RequestedValidators[valIdx].String()
	valReports, ok := env.rawDataReports[valAddr]
	if !ok {
		return nil, 0, types.ErrItemNotFound
	}
	valReport, ok := valReports[types.EID(eid)]
	if !ok {
		return nil, 0, types.ErrItemNotFound
	}
	return valReport.Data, valReport.ExitCode, nil
}
