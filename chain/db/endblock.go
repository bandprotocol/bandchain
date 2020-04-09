package db

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
	}
}
