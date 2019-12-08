package keeper

import (
	"testing"

	"github.com/bandprotocol/bandx/oracle/x/oracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryRequestById(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before set new request
	acsBytes, err := querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// It must be requestID = 0
	request := types.DataPoint{RequestID: 0, CodeHash: []byte(nil), ReportEndAt: 0, Result: []byte(nil)}
	acs, errJSON := codec.MarshalJSONIndent(keeper.cdc, request)
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)

	// set code
	code := []byte("code")
	codeHash := keeper.SetCode(ctx, code)

	// set request
	datapoint := types.NewDataPoint(1, codeHash, 3)
	keeper.SetRequest(ctx, 1, datapoint)

	// create query
	acsBytes, err = querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// Use bytes format for comparison
	request = types.DataPoint{RequestID: 1, CodeHash: codeHash, ReportEndAt: 3, Result: []byte(nil)}
	acs, errJSON = codec.MarshalJSONIndent(keeper.cdc, request)
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)
}

func TestQueryPendingRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// Read the state. The state should contain 0 pending request.
	acsBytes, err := querier(
		ctx,
		[]string{"pending_request"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)
	// Use bytes format for comparison
	acs, errJSON := codec.MarshalJSONIndent(keeper.cdc, []uint64{})
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)

	// set request
	code := []byte("code")
	codeHash := keeper.SetCode(ctx, code)
	datapoint := types.NewDataPoint(2, codeHash, 3)
	keeper.SetRequest(ctx, 2, datapoint)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 2)
	keeper.SetPending(ctx, pendingRequests)

	// Read the state agian. The state should contain 1 pending request. That is reqID = 2.
	acsBytes, err = querier(
		ctx,
		[]string{"pending_request"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)
	// Use bytes format for comparison
	acs, errJSON = codec.MarshalJSONIndent(keeper.cdc, []uint64{2})
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)
}
