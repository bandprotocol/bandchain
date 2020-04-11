package db

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func (b *BandDB) HandleEndblockEvent(event abci.Event) {
	kvMap := make(map[string]string)
	for _, kv := range event.Attributes {
		kvMap[string(kv.Key)] = string(kv.Value)
	}

	switch event.Type {
	case oracle.EventTypeRequestExecute:
		{
			requestID, err := strconv.ParseInt(kvMap[oracle.AttributeKeyRequestID], 10, 64)
			if err != nil {
				panic(err)
			}

			numResolveStatus, err := strconv.ParseInt(kvMap[oracle.AttributeKeyResolveStatus], 10, 8)
			if err != nil {
				panic(err)
			}
			resolveStatus := oracle.ResolveStatus(numResolveStatus)

			// Get result from keeper
			var rawResult []byte
			rawResult = nil
			if resolveStatus == 1 {
				id := oracle.RequestID(requestID)
				request, sdkErr := b.OracleKeeper.GetRequest(b.ctx, id)
				if sdkErr != nil {
					panic(err)
				}
				result, sdkErr := b.OracleKeeper.GetResult(b.ctx, id, request.OracleScriptID, request.Calldata)
				if sdkErr != nil {
					panic(err)
				}
				rawResult = result.Data
			}
			err = b.ResolveRequest(requestID, resolveStatus, rawResult)
			if err != nil {
				panic(err)
			}
		}
	case staking.EventTypeCompleteUnbonding:
		{
			// Recalculate delegator account
			delegatorAddress, err := sdk.AccAddressFromBech32(kvMap[staking.AttributeKeyDelegator])
			if err != nil {
				panic(err)
			}
			err = b.SetAccountBalance(
				delegatorAddress,
				b.OracleKeeper.CoinKeeper.GetAllBalances(b.ctx, delegatorAddress),
				b.ctx.BlockHeight(),
			)
			if err != nil {
				panic(err)
			}
		}
	case channel.EventTypeSendPacket:
		packetType := ""
		data := []byte(kvMap[channel.AttributeKeyData])
		jsonMap := make(map[string]interface{})

		var responsePacket oracle.OracleResponsePacketData
		if err := oracle.ModuleCdc.UnmarshalJSON(data, &responsePacket); err == nil {
			packetType = "ORACLE RESPONSE"
			rawBytes, err := json.Marshal(responsePacket)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(rawBytes, &jsonMap)
			if err != nil {
				panic(err)
			}

			request, err := b.OracleKeeper.GetRequest(b.ctx, responsePacket.RequestID)
			if err != nil {
				panic(err)
			}
			oracleScript, err := b.OracleKeeper.GetOracleScript(b.ctx, request.OracleScriptID)
			if err != nil {
				panic(err)
			}
			jsonMap["oracleScriptID"] = request.OracleScriptID
			jsonMap["oracleScriptName"] = oracleScript.Name
			jsonMap["resolveStatus"] = parseResolveStatus(request.ResolveStatus)
		}

		if packetType == "" {
			panic("Unknown packet type")
		}

		sequence, err := strconv.ParseUint(kvMap[channel.AttributeKeySequence], 10, 64)
		if err != nil {
			panic(err)
		}

		chainID, err := b.getChainID(
			kvMap[channel.AttributeKeySrcChannel],
			kvMap[channel.AttributeKeySrcPort],
		)
		if err != nil {
			panic(err)
		}

		rawJson, err := json.Marshal(jsonMap)
		if err != nil {
			panic(err)
		}

		isIncoming := false
		err = b.tx.Create(&Packet{
			Type:        packetType,
			Sequence:    sequence,
			MyChannel:   kvMap[channel.AttributeKeySrcChannel],
			MyPort:      kvMap[channel.AttributeKeySrcPort],
			YourChainID: chainID,
			YourChannel: kvMap[channel.AttributeKeyDstChannel],
			YourPort:    kvMap[channel.AttributeKeyDstPort],
			BlockHeight: b.ctx.BlockHeight(),
			IsIncoming:  &isIncoming,
			Detail:      rawJson,
		}).Error

		if err != nil {
			panic(err)
		}
	}
}
