package keeper

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	owasm "github.com/bandprotocol/bandchain/go-owasm/api"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// nolint
const (
	PrepareFunc = "prepare"
	ExecuteFunc = "execute"
)

// PrepareRequest takes an request specification object, performs the prepare call, and saves
// the request object to store. Also emits events related to the request.
func (k Keeper) PrepareRequest(ctx sdk.Context, r types.RequestSpec, ibcInfo *types.IBCInfo) error {
	askCount := r.GetAskCount()
	if askCount > k.GetParam(ctx, types.KeyMaxAskCount) {
		return sdkerrors.Wrapf(types.ErrInvalidAskCount, "got: %d, max: %d", askCount, k.GetParam(ctx, types.KeyMaxAskCount))
	}
	// Consume gas for data requests. We trust that we have reasonable params that don't cause overflow.
	ctx.GasMeter().ConsumeGas(k.GetParam(ctx, types.KeyBaseRequestGas), "BASE_REQUEST_FEE")
	ctx.GasMeter().ConsumeGas(askCount*k.GetParam(ctx, types.KeyPerValidatorRequestGas), "PER_VALIDATOR_REQUEST_FEE")
	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(askCount), k.GetRequestCount(ctx)+1)
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
	exitCode := owasm.Prepare(code, env) // TODO: Don't forget about prepare gas!
	if exitCode != 0 {
		k.Logger(ctx).Info(fmt.Sprintf("failed to prepare request with code: %d", exitCode))
		return types.ErrBadWasmExecution
	}
	// Preparation complete! It's time to collect raw request ids.
	for _, rawReq := range env.GetRawRequests() {
		req.RawRequestIDs = append(req.RawRequestIDs, rawReq.ExternalID)
	}
	// We now have everything we need to the request, so let's add it to the store.
	id := k.AddRequest(ctx, req)
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
	)
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
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", rawReq.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(rawReq.Calldata)),
		))
	}
	return nil
}

// ResolveRequest resolves the given request, sends response packet out (if applicable),
// and saves result hash to the store. Assumes that the given request is in a resolvable state.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecEnv(req, ctx.BlockTime().Unix(), 0)
	env.SetReports(k.GetReports(ctx, reqID))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	exitCode := owasm.Execute(code, env) // TODO: Don't forget about gas!
	var res types.OracleResponsePacketData
	if exitCode != 0 {
		k.Logger(ctx).Info(fmt.Sprintf(
			"failed to execute request id: %d with code: %d", reqID, exitCode,
		))
		res = k.SaveResult(ctx, reqID, types.ResolveStatus_Failure, nil)
	} else {
		res = k.SaveResult(ctx, reqID, types.ResolveStatus_Success, env.Retdata)
	}
	if req.IBCInfo != nil {
		k.SendOracleResponse(ctx, req.IBCInfo.SourcePort, req.IBCInfo.SourceChannel, res)
	}
}
