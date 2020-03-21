package zoracle

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
)

type ExecutionEnvironment struct {
	ctx    sdk.Context
	keeper Keeper
	//
	requestID       types.RequestID
	request         types.Request
	now             int64
	maxResultSize   int64
	maxCalldataSize int64
	//
	rawDataRequests []types.RawDataRequestWithExternalID
	rawDataReports  map[string]types.RawDataReport
}

func NewExecutionEnvironment(
	ctx sdk.Context, keeper Keeper, requestID types.RequestID,
) (ExecutionEnvironment, sdk.Error) {
	request, err := keeper.GetRequest(ctx, requestID)
	if err != nil {
		return ExecutionEnvironment{}, err
	}
	return ExecutionEnvironment{
		ctx:             ctx,
		keeper:          keeper,
		requestID:       requestID,
		request:         request,
		now:             ctx.BlockTime().Unix(),
		maxResultSize:   keeper.MaxResultSize(ctx),
		maxCalldataSize: keeper.MaxCalldataSize(ctx),
		rawDataRequests: []types.RawDataRequestWithExternalID{},
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

func (env *ExecutionEnvironment) GetMaximumResultSize() int64 {
	return env.maxResultSize
}

func (env *ExecutionEnvironment) GetMaximumCalldataOfDataSourceSize() int64 {
	return env.maxCalldataSize
}

func (env *ExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalDataID int64,
	calldata []byte,
) {
	env.rawDataRequests = append(env.rawDataRequests, types.NewRawDataRequestWithExternalID(
		types.ExternalID(externalDataID),
		types.NewRawDataRequest(types.DataSourceID(dataSourceID), calldata),
	))
}

func (env *ExecutionEnvironment) GetExternalData(
	externalDataID int64,
	validatorIndex int64,
) ([]byte, uint8, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return nil, 0, errors.New("validator out of range")
	}
	validatorAddress := env.request.RequestedValidators[validatorIndex]
	rawReport, err := env.keeper.GetRawDataReport(
		env.ctx,
		env.requestID,
		types.ExternalID(externalDataID),
		validatorAddress,
	)
	if err != nil {
		return nil, 0, err
	}
	return rawReport.Data, rawReport.ExitCode, nil
}
