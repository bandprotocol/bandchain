package oracle

import (
	"errors"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type ExecutionEnvironment struct {
	// TODO: SWIT CLEAN UP
	isPrepare     bool
	receivedCount int64

	requestID              types.RequestID
	request                types.Request
	now                    int64
	maxResultSize          int64
	maxCalldataSize        int64
	maxRawDataRequestCount int64
	rawDataRequests        []types.RawRequest
	rawDataReports         map[string]types.RawDataReport
}

func NewExecutionEnvironment(
	ctx sdk.Context, keeper Keeper, requestID types.RequestID, isPrepare bool, receivedCount int64,
) ExecutionEnvironment {
	return ExecutionEnvironment{
		isPrepare:              isPrepare,
		receivedCount:          receivedCount,
		requestID:              requestID,
		request:                keeper.MustGetRequest(ctx, requestID),
		now:                    ctx.BlockTime().Unix(),
		maxResultSize:          int64(keeper.GetParam(ctx, KeyMaxResultSize)),
		maxCalldataSize:        int64(keeper.GetParam(ctx, KeyMaxCalldataSize)),
		maxRawDataRequestCount: int64(keeper.GetParam(ctx, KeyMaxDataSourceCountPerRequest)),
		rawDataRequests:        []types.RawRequest{},
		rawDataReports:         make(map[string]types.RawDataReport),
	}
}

func (env *ExecutionEnvironment) SaveRawDataRequests(ctx sdk.Context, keeper Keeper) error {
	for _, r := range env.rawDataRequests {
		err := keeper.AddRawRequest(
			ctx, env.requestID, types.NewRawRequest(r.ExternalID, r.DataSourceID, r.Calldata),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *ExecutionEnvironment) LoadRawDataReports(
	ctx sdk.Context,
	keeper Keeper,
) error {
	for _, report := range keeper.GetReports(ctx, env.requestID) {
		for _, reportv1 := range report.RawDataReports {

			key := string(types.RawDataReportStoreKeyUnique(env.requestID, reportv1.ExternalID, report.Validator))
			env.rawDataReports[key] = types.NewRawDataReport(reportv1.ExitCode, reportv1.Data)
		}
	}
	return nil
}

func (env *ExecutionEnvironment) GetCurrentRequestID() int64 {
	return int64(env.requestID)
}

func (env *ExecutionEnvironment) GetRequestedValidatorCount() int64 {
	return int64(len(env.request.RequestedValidators))
}

func (env *ExecutionEnvironment) GetSufficientValidatorCount() int64 {
	return env.request.SufficientValidatorCount
}

func (env *ExecutionEnvironment) GetReceivedValidatorCount() int64 {
	return env.receivedCount
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
	if !env.isPrepare {
		return env.now
	}
	return 0
}

func (env *ExecutionEnvironment) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, errors.New("validator out of range")
	}
	return env.request.RequestedValidators[validatorIndex], nil
}

func (env *ExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalID int64,
	calldata []byte,
) error {
	if int64(len(calldata)) > env.maxCalldataSize {
		return errors.New("calldata size limit exceeded")
	}
	if int64(len(env.rawDataRequests)) >= env.maxRawDataRequestCount {
		return errors.New("cannot request more than maxRawDataRequestCount")
	}

	env.rawDataRequests = append(env.rawDataRequests, types.NewRawRequest(
		types.ExternalID(externalID),
		types.DataSourceID(dataSourceID),
		calldata,
	))
	return nil
}

func (env *ExecutionEnvironment) GetExternalData(
	externalID int64,
	validatorIndex int64,
) ([]byte, uint8, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, errors.New("validator out of range")
	}
	validatorAddress := env.request.RequestedValidators[validatorIndex]

	key := string(types.RawDataReportStoreKeyUnique(env.requestID, types.EID(externalID), validatorAddress))

	rawDataReport, ok := env.rawDataReports[key]

	if !ok {
		return nil, 0, sdkerrors.Wrapf(types.ErrItemNotFound, "Unable to find raw data report with request ID (%d) external ID (%d) from (%s)", env.requestID, externalID, validatorAddress.String())
	}

	return rawDataReport.Data, rawDataReport.ExitCode, nil
}
