package keeper

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	PrepareFunc = "prepare"
	ExecuteFunc = "execute"
)

// PrepareRequest takes an request specification object, performs the prepare call, and saves
// the request object to store. Also emits events related to the request.
func (k Keeper) PrepareRequest(ctx sdk.Context, r types.RequestSpec, ibcInfo *types.IBCInfo) error {
	// TODO: FIX ME! Consume a fixed gas amount for processing oracle request.
	nextID := k.GetRequestCount(ctx) + 1
	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(r.GetAskCount()), nextID)
	if err != nil {
		return err
	}
	// Create a request object. Note that RawRequestIDs will be populated after preparation is done.
	req := types.NewRequest(
		r.GetOracleScriptID(), r.GetCalldata(), validators, r.GetMinCount(),
		ctx.BlockHeight(), ctx.BlockTime().Unix(), r.GetClientID(), ibcInfo, nil,
	)
	// Create an execution environment and call Owasm prepare function.
	env := types.NewExecEnv(req, ctx.BlockTime().Unix(), int64(k.GetParam(ctx, types.KeyMaxRawRequestCount)))
	script, err := k.GetOracleScript(ctx, req.OracleScriptID)
	if err != nil {
		return err
	}
	code := k.GetFile(script.Filename)
	_, _, err = k.OwasmExecute(env, code, PrepareFunc, req.Calldata, types.WasmPrepareGas)
	if err != nil {
		k.Logger(ctx).Info(fmt.Sprintf("failed to prepare request with error: %s", err.Error()))
		return types.ErrBadWasmExecution
	}
	// Preparation complete! It's time to collect raw request ids and ask for more gas.
	for _, rawReq := range env.GetRawRequests() {
		// TODO: FIX ME! Consume more gas for each raw request
		req.RawRequestIDs = append(req.RawRequestIDs, rawReq.ExternalID)
	}
	// We now have everything we need to the request, so let's add it to the store.
	id := k.AddRequest(ctx, req)
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)))
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val.String()))
	}
	ctx.EventManager().EmitEvent(event)
	// Emit an event for each of the raw data requests.
	for _, rawReq := range env.GetRawRequests() {
		ds, err := k.GetDataSource(ctx, rawReq.DataSourceID)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, fmt.Sprintf("%d", rawReq.DataSourceID)),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", rawReq.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(rawReq.Calldata)),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds.Filename),
		))
	}
	return nil
}

// ResolveRequest resolves the given request, sends response packet out (if applicable),
// and saves result hash to the store. Assumes that the given request is in a resolvable state.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecEnv(req, ctx.BlockTime().Unix(), int64(k.GetParam(ctx, types.KeyMaxRawRequestCount)))
	env.SetReports(k.GetReports(ctx, reqID))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	result, _, err := k.OwasmExecute(env, code, ExecuteFunc, req.Calldata, types.WasmExecuteGas)
	var res types.OracleResponsePacketData
	if err != nil {
		k.Logger(ctx).Info(fmt.Sprintf(
			"failed to execute request id: %d with error: %s", reqID, err.Error(),
		))
		res = k.SaveResult(ctx, reqID, types.ResolveStatus_Failure, nil)
	} else {
		res = k.SaveResult(ctx, reqID, types.ResolveStatus_Success, result)
	}
	if req.IBCInfo != nil {
		k.SendOracleResponse(ctx, req.IBCInfo.SourcePort, req.IBCInfo.SourceChannel, res)
	}
}
