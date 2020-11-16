package keeper_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryPendingRequests(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)

	// Add 3 requests
	k.SetRequestLastExpired(ctx, 40)
	k.SetRequest(ctx, 41, defaultRequest())
	k.SetRequest(ctx, 42, defaultRequest())
	k.SetRequest(ctx, 43, defaultRequest())
	k.SetRequestCount(ctx, 43)

	// Fulfill some requests
	k.SetReport(ctx, 41, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator2.ValAddress, true, nil))

	q := keeper.NewQuerier(k)

	test := func(args []string, expected []types.RequestID) {
		raw, err := q(ctx, append([]string{types.QueryPendingRequests}, args...), abci.RequestQuery{})
		require.NoError(t, err)

		var queryRequest types.QueryResult
		require.NoError(t, json.Unmarshal(raw, &queryRequest))

		var rawRequestIDs []string
		types.ModuleCdc.MustUnmarshalJSON(queryRequest.Result, &rawRequestIDs)

		var requestIDs []types.RequestID
		for _, r := range rawRequestIDs {
			id, err := strconv.ParseInt(r, 10, 64)
			require.NoError(t, err)

			requestIDs = append(requestIDs, types.RequestID(id))
		}

		require.Equal(t, expected, requestIDs)
	}

	test([]string{}, []types.RequestID{41, 42, 43})
	test([]string{testapp.Validator1.ValAddress.String()}, []types.RequestID{42, 43})
	test([]string{testapp.Validator2.ValAddress.String()}, []types.RequestID{41, 43})
	test([]string{testapp.Validator3.ValAddress.String()}, nil)
}
