package oracle

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/bandprotocol/bandchain/chain/owasm"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

type ResolveContext struct {
	endBlockExecuteGasLimit uint64
	gasConsumed             uint64
}

func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func newRequestExecuteEvent(requestID RequestID, resolveStatus types.ResolveStatus) sdk.Event {
	return sdk.NewEvent(
		types.EventTypeRequestExecute,
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", requestID)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", resolveStatus)),
	)
}

func handleResolveRequest(
	ctx sdk.Context,
	keeper Keeper,
	requestID types.RequestID,
	resolveContext *ResolveContext,
) (sdk.Event, OracleResponsePacketData, bool) {
	request, err := keeper.GetRequest(ctx, requestID)
	if err != nil { // should never happen
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	// TODO: Refactor this code. For now we hardcode execute gas to 100k
	executeGas := uint64(100000)

	// Discard the request if execute gas is greater than EndBlockExecuteGasLimit.
	if executeGas > resolveContext.endBlockExecuteGasLimit {
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	estimatedGasConsumed, overflow := addUint64Overflow(resolveContext.gasConsumed, executeGas)
	if overflow || estimatedGasConsumed > resolveContext.endBlockExecuteGasLimit {
		return sdk.Event{},
			OracleResponsePacketData{},
			true
	}

	env, err := NewExecutionEnvironment(ctx, keeper, requestID)
	if err != nil { // should never happen
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	err = env.LoadRawDataReports(ctx, keeper)
	if err != nil { // should never happen
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
	if err != nil { // should never happen
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	result, gasUsed, errOwasm := owasm.Execute(
		&env, script.Code, "execute", request.Calldata, executeGas,
	)

	if gasUsed > executeGas {
		gasUsed = executeGas
	}

	resolveContext.gasConsumed, overflow = addUint64Overflow(resolveContext.gasConsumed, gasUsed)
	// Must never overflow because we already checked for overflow above with
	// gasConsumed + executeGas (which is >= gasUsed).
	if overflow {
		panic(sdk.ErrorGasOverflow{Descriptor: "oracle::handleEndBlock: Gas overflow"})
	}

	if errOwasm != nil {
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	errResult := keeper.AddResult(ctx, requestID, request.OracleScriptID, request.Calldata, result)
	if errResult != nil {
		keeper.SetResolve(ctx, requestID, types.Failure)
		return newRequestExecuteEvent(requestID, types.Failure),
			NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Failure, ""),
			false
	}

	keeper.SetResolve(ctx, requestID, types.Success)
	event := newRequestExecuteEvent(requestID, types.Success)
	event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyResult, string(result)))
	return event,
		NewOracleResponsePacketData(request.ClientID, requestID, int64(len(request.ReceivedValidators)), request.RequestTime, ctx.BlockTime().Unix(), types.Success, hex.EncodeToString(result)),
		false

}

func handleEndBlock(ctx sdk.Context, keeper Keeper) {
	pendingList := keeper.GetPendingResolveList(ctx)
	firstUnresolvedRequestIndex := len(pendingList)
	resolveContext := ResolveContext{
		endBlockExecuteGasLimit: keeper.GetParam(ctx, types.KeyEndBlockExecuteGasLimit),
		gasConsumed:             0,
	}
	events := []sdk.Event{}
	for i, requestID := range pendingList {
		event, packet, stopped := handleResolveRequest(ctx, keeper, requestID, &resolveContext)

		if stopped {
			firstUnresolvedRequestIndex = i
			break
		}

		request, err := keeper.GetRequest(ctx, requestID)
		events = append(events, event)
		sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, request.SourcePort, request.SourceChannel)
		if !found {
			fmt.Println("SOURCE NOT FOUND", request.SourcePort, request.SourceChannel)
			continue
		}

		destinationPort := sourceChannelEnd.Counterparty.PortID
		destinationChannel := sourceChannelEnd.Counterparty.ChannelID

		// get the next sequence
		sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(ctx, request.SourcePort, request.SourceChannel)
		if !found {
			fmt.Println("SEQUENCE NOT FOUND", request.SourcePort, request.SourceChannel)
			continue
		}

		err = keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
			sequence, request.SourcePort, request.SourceChannel, destinationPort, destinationChannel,
			1000000000, // Arbitrarily high timeout for now
		))

		if err != nil {
			fmt.Println("SEND PACKET ERROR", err)
		}
	}

	ctx.EventManager().EmitEvents(events)
	keeper.SetPendingResolveList(ctx, pendingList[firstUnresolvedRequestIndex:])
}
