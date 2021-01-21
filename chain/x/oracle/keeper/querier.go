package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryParams:
			return queryParameters(ctx, keeper)
		case types.QueryCounts:
			return queryCounts(ctx, keeper)
		case types.QueryData:
			return queryData(ctx, path[1:], keeper)
		case types.QueryDataSources:
			return queryDataSourceByID(ctx, path[1:], keeper)
		case types.QueryOracleScripts:
			return queryOracleScriptByID(ctx, path[1:], keeper)
		case types.QueryRequests:
			return queryRequestByID(ctx, path[1:], keeper)
		case types.QueryValidatorStatus:
			return queryValidatorStatus(ctx, path[1:], keeper)
		case types.QueryReporters:
			return queryReporters(ctx, path[1:], keeper)
		case types.QueryActiveValidators:
			return queryActiveValidators(ctx, keeper)
		case types.QueryPendingRequests:
			return queryPendingRequests(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown oracle query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	return types.QueryOK(k.GetParams(ctx))
}

func queryCounts(ctx sdk.Context, k Keeper) ([]byte, error) {
	return types.QueryOK(types.QueryCountsResult{
		DataSourceCount:   k.GetDataSourceCount(ctx),
		OracleScriptCount: k.GetOracleScriptCount(ctx),
		RequestCount:      k.GetRequestCount(ctx),
	})
}

func queryData(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "data hash not specified")
	}
	return k.fileCache.GetFile(path[0])
}

func queryDataSourceByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "data source not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return types.QueryBadRequest(err.Error())
	}
	dataSource, err := k.GetDataSource(ctx, types.DataSourceID(id))
	if err != nil {
		return types.QueryNotFound(err.Error())
	}
	return types.QueryOK(dataSource)
}

func queryOracleScriptByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "oracle script not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return types.QueryBadRequest(err.Error())
	}
	oracleScript, err := k.GetOracleScript(ctx, types.OracleScriptID(id))
	if err != nil {
		return types.QueryNotFound(err.Error())
	}
	return types.QueryOK(oracleScript)
}

func queryRequestByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "request not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return types.QueryBadRequest(err.Error())
	}
	request, err := k.GetRequest(ctx, types.RequestID(id))
	if err != nil {
		return types.QueryNotFound(err.Error())
	}
	reports := k.GetReports(ctx, types.RequestID(id))
	if !k.HasResult(ctx, types.RequestID(id)) {
		return types.QueryOK(types.QueryRequestResult{
			Request: request,
			Reports: reports,
			Result:  nil,
		})
	}
	result := k.MustGetResult(ctx, types.RequestID(id))
	return types.QueryOK(types.QueryRequestResult{
		Request: request,
		Reports: reports,
		Result:  &result,
	})
}

func queryValidatorStatus(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "validator address not specified")
	}
	validatorAddress, err := sdk.ValAddressFromBech32(path[0])
	if err != nil {
		return types.QueryBadRequest(err.Error())
	}
	return types.QueryOK(k.GetValidatorStatus(ctx, validatorAddress))
}

func queryReporters(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "validator address not specified")
	}
	validatorAddress, err := sdk.ValAddressFromBech32(path[0])
	if err != nil {
		return types.QueryBadRequest(err.Error())
	}
	return types.QueryOK(k.GetReporters(ctx, validatorAddress))
}

func queryActiveValidators(ctx sdk.Context, k Keeper) ([]byte, error) {
	vals := []types.QueryActiveValidatorResult{}
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val exported.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				vals = append(vals, types.QueryActiveValidatorResult{
					Address: val.GetOperator(),
					Power:   val.GetTokens().Uint64(),
				})
			}
			return false
		})
	return types.QueryOK(vals)
}

func queryPendingRequests(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) > 1 {
		return types.QueryBadRequest("too many arguments")
	}

	var valAddress *sdk.ValAddress
	if len(path) == 1 {
		valAddress = new(sdk.ValAddress)
		address, err := sdk.ValAddressFromBech32(path[0])
		if err != nil {
			return types.QueryBadRequest(err.Error())
		}

		*valAddress = address
	}

	lastExpired := k.GetRequestLastExpired(ctx)
	requestCount := k.GetRequestCount(ctx)

	var pendingIDs []types.RequestID
	for id := lastExpired + 1; int64(id) <= requestCount; id++ {

		req := k.MustGetRequest(ctx, id)

		// If all validators reported on this request, then skip it.
		reports := k.GetReports(ctx, id)
		if len(reports) == len(req.RequestedValidators) {
			continue
		}

		// Skip if validator hasn't been assigned.
		if valAddress != nil {

			// If the validator isn't in requested validators set, then skip it.
			isValidator := false
			for _, v := range req.RequestedValidators {
				if valAddress.Equals(v) {
					isValidator = true
					break
				}
			}

			if !isValidator {
				continue
			}

			// If the validator has reported, then skip it.
			reported := false
			for _, r := range reports {
				if valAddress.Equals(r.Validator) {
					reported = true
					break
				}
			}

			if reported {
				continue
			}
		}

		pendingIDs = append(pendingIDs, id)
	}

	return types.QueryOK(pendingIDs)
}
