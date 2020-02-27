package zoracle

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

type ExecutionEnvironment struct {
	ctx       sdk.Context
	keeper    Keeper
	requestID int64
	request   types.Request
}

func NewExecutionEnvironment(
	ctx sdk.Context, keeper Keeper, requestID int64,
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
	return env.requestID
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
		return []byte{}, errors.New("validator out of range")
	}
	return env.request.RequestedValidators[validatorIndex], nil
}

func (env *ExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalDataID int64,
	calldata []byte,
) error {
	return env.keeper.AddNewRawDataRequest(env.ctx, env.requestID, externalDataID, dataSourceID, calldata)
}

func (env *ExecutionEnvironment) GetExternalData(
	externalDataID int64,
	validatorIndex int64,
) ([]byte, error) {
	if validatorIndex < 0 || validatorIndex >= int64(len(env.request.RequestedValidators)) {
		return []byte{}, errors.New("validator out of range")
	}
	validatorAddress := env.request.RequestedValidators[validatorIndex]

	return env.keeper.GetRawDataReport(env.ctx, env.requestID, externalDataID, validatorAddress)
}
