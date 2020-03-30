package zoracle

import (
	"errors"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExecutionEnvironment struct {
	ctx       sdk.Context
	keeper    Keeper
	requestID types.RequestID
	request   types.Request
}

func NewExecutionEnvironment(
	ctx sdk.Context, keeper Keeper, requestID types.RequestID,
) (ExecutionEnvironment, sdk.Error) {
	request, err := keeper.GetRequest(ctx, requestID)
	if err != nil {
		return ExecutionEnvironment{}, err
	}
	return ExecutionEnvironment{
		ctx:       ctx,
		keeper:    keeper,
		requestID: requestID,
		request:   request,
	}, nil
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
		return env.ctx.BlockTime().Unix()
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
	return int64(env.keeper.GetParam(env.ctx, KeyMaxResultSize))
}

func (env *ExecutionEnvironment) GetMaximumCalldataOfDataSourceSize() int64 {
	return int64(env.keeper.GetParam(env.ctx, KeyMaxCalldataSize))
}

func (env *ExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalDataID int64,
	calldata []byte,
) error {
	return env.keeper.AddNewRawDataRequest(env.ctx, env.requestID, types.ExternalID(externalDataID), types.DataSourceID(dataSourceID), calldata)
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
