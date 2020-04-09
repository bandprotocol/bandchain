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

func handleEndBlock(ctx sdk.Context, keeper Keeper) {
	pendingList := keeper.GetPendingResolveList(ctx)
	endBlockExecuteGasLimit := keeper.GetParam(ctx, types.KeyEndBlockExecuteGasLimit)
	gasConsumed := uint64(0)
	firstUnresolvedRequestIndex := len(pendingList)
	events := []sdk.Event{}
	for i, requestID := range pendingList {
		request, err := keeper.GetRequest(ctx, requestID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		// Discard the request if execute gas is greater than EndBlockExecuteGasLimit.
		if request.ExecuteGas > endBlockExecuteGasLimit {
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		estimatedGasConsumed, overflow := addUint64Overflow(gasConsumed, request.ExecuteGas)
		if overflow || estimatedGasConsumed > endBlockExecuteGasLimit {
			firstUnresolvedRequestIndex = i
			break
		}

		env, err := NewExecutionEnvironment(ctx, keeper, requestID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		err = env.LoadRawDataReports(ctx, keeper)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		result, gasUsed, errOwasm := owasm.Execute(
			&env, script.Code, "execute", request.Calldata, request.ExecuteGas,
		)

		if gasUsed > request.ExecuteGas {
			gasUsed = request.ExecuteGas
		}

		gasConsumed, overflow = addUint64Overflow(gasConsumed, gasUsed)
		// Must never overflow because we already checked for overflow above with
		// gasConsumed + request.ExecuteGas (which is >= gasUsed).
		if overflow {
			panic(sdk.ErrorGasOverflow{Descriptor: "oracle::handleEndBlock: Gas overflow"})
		}

		if errOwasm != nil {
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		errResult := keeper.AddResult(ctx, requestID, request.OracleScriptID, request.Calldata, result)
		if errResult != nil {
			keeper.SetResolve(ctx, requestID, types.Failure)
			events = append(events, newRequestExecuteEvent(requestID, types.Failure))
			continue
		}

		keeper.SetResolve(ctx, requestID, types.Success)
		event := newRequestExecuteEvent(requestID, types.Success)
		event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyResult, string(result)))
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

		packet := NewOracleResponsePacketData(requestID, request.ClientID, hex.EncodeToString(result))

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
