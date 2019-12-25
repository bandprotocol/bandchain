package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	request := types.NewRequest([]byte(nil), []byte(nil), 0)
	acs, errJSON := codec.MarshalJSONIndent(keeper.cdc, types.NewRequestWithReport(request, []byte{}, []types.ValidatorReport{}))
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)

	// set code
	code := []byte("code")
	owner := sdk.AccAddress([]byte("owner"))
	codeHash := keeper.SetCode(ctx, code, owner)

	// set request
	request = types.NewRequest(codeHash, []byte("params"), 3)
	keeper.SetRequest(ctx, 1, request)
	result := []byte("result")
	keeper.SetResult(ctx, 1, codeHash, []byte("params"), result)

	// create query
	acsBytes, err = querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// Use bytes format for comparison
	request = types.NewRequest(codeHash, []byte("params"), 3)
	acs, errJSON = codec.MarshalJSONIndent(
		keeper.cdc,
		types.NewRequestWithReport(request, result, []types.ValidatorReport{}),
	)
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
	owner := sdk.AccAddress([]byte("owner"))
	code := []byte("code")
	codeHash := keeper.SetCode(ctx, code, owner)
	request := types.NewRequest(codeHash, []byte("params"), 3)
	keeper.SetRequest(ctx, 2, request)

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
