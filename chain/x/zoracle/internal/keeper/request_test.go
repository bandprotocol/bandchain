package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func newDefaultRequest() types.Request {
	return types.NewRequest(
		1,
		[]byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("validator1")), sdk.ValAddress([]byte("validator2"))},
		2,
		0,
		1581503227,
		100,
	)
}

func TestGetterSetterRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetRequest(ctx, 1)
	require.NotNil(t, err)

	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestAddSubmitValidator tests keeper can add valid validator to request
func TestAddSubmitValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.AddSubmitValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Nil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.SubmittedValidatorList = []sdk.ValAddress{sdk.ValAddress([]byte("validator1"))}
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestAddSubmitValidatorOnInvalidRequest tests keeper must return if add on invalid request
func TestAddSubmitValidatorOnInvalidRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	err := keeper.AddSubmitValidator(ctx, 2, sdk.ValAddress([]byte("validator1")))
	require.Equal(t, types.CodeRequestNotFound, err.Code())
}

// TestAddInvalidSubmitValidator tests keeper return error if try to add new validator that doesn't contain in list.
func TestAddInvalidSubmitValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.AddSubmitValidator(ctx, 1, sdk.ValAddress([]byte("validator3")))
	require.Equal(t, types.CodeInvalidValidator, err.Code())

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestAddDuplicateSubmitValidator tests keeper return error if try to add new validator that already in list.
func TestAddDuplicateSubmitValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	// First add must return nil
	err := keeper.AddSubmitValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Nil(t, err)

	// Second add must return duplicate error
	err = keeper.AddSubmitValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Equal(t, types.CodeDuplicateValidator, err.Code())

	// Check final output
	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.SubmittedValidatorList = []sdk.ValAddress{sdk.ValAddress([]byte("validator1"))}
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestSetResolved tests keeper can set resolved status to request
func TestSetResolved(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.SetResolve(ctx, 1, true)
	require.Nil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.IsResolved = true
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestSetResolvedOnInvalidRequest tests keeper must return if set on invalid request
func TestSetResolvedOnInvalidRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	err := keeper.SetResolve(ctx, 2, true)
	require.Equal(t, types.CodeRequestNotFound, err.Code())
}

// Can get/set unresolved request correctly and set empty case
func TestGetSetUnresolvedRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetUnresolvedRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)

	keeper.SetUnresolvedRequests(ctx, []int64{1, 2, 3})

	reqIDs = keeper.GetUnresolvedRequests(ctx)
	require.Equal(t, []int64{1, 2, 3}, reqIDs)

	keeper.SetUnresolvedRequests(ctx, []int64{})
	reqIDs = keeper.GetUnresolvedRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)
}

// Can set pending request will set only unique request IDs
func TestGetSetUnresolvedRequestsUnique(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetUnresolvedRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)

	keeper.SetUnresolvedRequests(ctx, []int64{3, 2, 3, 1, 2, 1, 3, 2, 1})
	reqIDs = keeper.GetUnresolvedRequests(ctx)
	// no guarantee of an order
	require.Equal(t, []int64{3, 2, 1}, reqIDs)
}
