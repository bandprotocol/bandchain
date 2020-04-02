package zoracle

import (
	"errors"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExecutionEnvironment struct {
	requestID              types.RequestID
	request                types.Request
	now                    int64
	maxResultSize          int64
	maxCalldataSize        int64
	maxRawDataRequestCount int64
	rawDataRequests        []types.RawDataRequestWithExternalID
	rawDataReports         map[string]types.RawDataReport
}

func NewExecutionEnvironment(
	ctx sdk.Context, keeper Keeper, requestID types.RequestID,
) (ExecutionEnvironment, sdk.Error) {
	request, err := keeper.GetRequest(ctx, requestID)
	if err != nil {
		return ExecutionEnvironment{}, err
	}
	return ExecutionEnvironment{
		requestID:              requestID,
		request:                request,
		now:                    ctx.BlockTime().Unix(),
		maxResultSize:          int64(keeper.GetParam(ctx, KeyMaxResultSize)),
		maxCalldataSize:        int64(keeper.GetParam(ctx, KeyMaxCalldataSize)),
		maxRawDataRequestCount: int64(keeper.GetParam(ctx, KeyMaxDataSourceCountPerRequest)),
		rawDataRequests:        []types.RawDataRequestWithExternalID{},
		rawDataReports:         make(map[string]types.RawDataReport),
	}, nil
}

func (env *ExecutionEnvironment) SaveRawDataRequests(ctx sdk.Context, keeper Keeper) sdk.Error {
	for _, r := range env.rawDataRequests {
		err := keeper.AddNewRawDataRequest(
			ctx, env.requestID, r.ExternalID,
			r.RawDataRequest.DataSourceID, r.RawDataRequest.Calldata,
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
) sdk.Error {

	for iterator := keeper.GetRawDataReportsIterator(ctx, env.requestID); iterator.Valid(); iterator.Next() {
		validatorAddress, externalID := types.GetValidatorAddressAndExternalID(iterator.Key(), env.requestID)

		rawDataReport, err := keeper.GetRawDataReport(
			ctx,
			env.requestID,
			externalID,
			validatorAddress,
		)
		if err != nil { // should never happen
			return err
		}

		key := string(types.RawDataReportStoreKey(env.requestID, externalID, validatorAddress))

		env.rawDataReports[key] = rawDataReport
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
	return int64(len(env.request.ReceivedValidators))
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
	if int64(len(env.request.ReceivedValidators)) >= env.request.SufficientValidatorCount {
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
	externalDataID int64,
	calldata []byte,
) error {
	if int64(len(calldata)) > env.maxCalldataSize {
		return errors.New("calldata size limit exceeded")
	}
	if int64(len(env.rawDataRequests)) >= env.maxRawDataRequestCount {
		return errors.New("cannot request more than maxRawDataRequestCount")
	}

	env.rawDataRequests = append(env.rawDataRequests, types.NewRawDataRequestWithExternalID(
		types.ExternalID(externalDataID),
		types.NewRawDataRequest(types.DataSourceID(dataSourceID), calldata),
	))
	return nil
}

func (env *ExecutionEnvironment) GetExternalData(
	externalDataID int64,
	validatorIndex int64,
) ([]byte, uint8, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, errors.New("validator out of range")
	}
	validatorAddress := env.request.RequestedValidators[validatorIndex]
	key := string(types.RawDataReportStoreKey(env.requestID, types.ExternalID(externalDataID), validatorAddress))

	rawDataReport, ok := env.rawDataReports[key]

	if !ok {
		return nil, 0, types.ErrItemNotFound("Unable to find raw data report with request ID (%d) external ID (%d) from (%s)", env.requestID, externalDataID, validatorAddress.String())
	}

	return rawDataReport.Data, rawDataReport.ExitCode, nil
}
