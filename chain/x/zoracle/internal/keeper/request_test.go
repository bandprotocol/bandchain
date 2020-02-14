package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

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

func TestRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	calldata := []byte("calldata")
	_, err := keeper.Request(ctx, 1, calldata, 2, 2, 100)
	require.NotNil(t, err)

	script := getTestOracleScript()
	keeper.SetOracleScript(ctx, 1, script)
	_, err = keeper.Request(ctx, 1, calldata, 2, 2, 100)
	require.NotNil(t, err)

	// Setup validator1
	_, err = keeper.Request(ctx, 1, calldata, 2, 2, 100)
	require.NotNil(t, err)

	// Setup validator2
	_, err = keeper.Request(ctx, 1, calldata, 2, 2, 100)
	require.NotNil(t, err)

	dataSource := getTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)
	requestID, err := keeper.Request(ctx, 1, calldata, 2, 2, 100)
	require.Nil(t, err)
	require.Equal(t, int64(1), requestID)
}

// TestAddNewReceiveValidator tests keeper can add valid validator to request
func TestAddNewReceiveValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Nil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.ReceivedValidators = []sdk.ValAddress{sdk.ValAddress([]byte("validator1"))}
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestAddNewReceiveValidatorOnInvalidRequest tests keeper must return if add on invalid request
func TestAddNewReceiveValidatorOnInvalidRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	err := keeper.AddNewReceiveValidator(ctx, 2, sdk.ValAddress([]byte("validator1")))
	require.Equal(t, types.CodeRequestNotFound, err.Code())
}

// TestAddInvalidValidator tests keeper return error if try to add new validator that doesn't contain in list.
func TestAddInvalidValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("validator3")))
	require.Equal(t, types.CodeInvalidValidator, err.Code())

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestAddDuplicateValidator tests keeper return error if try to add new validator that already in list.
func TestAddDuplicateValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	// First add must return nil
	err := keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Nil(t, err)

	// Second add must return duplicate error
	err = keeper.AddNewReceiveValidator(ctx, 1, sdk.ValAddress([]byte("validator1")))
	require.Equal(t, types.CodeDuplicateValidator, err.Code())

	// Check final output
	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.ReceivedValidators = []sdk.ValAddress{sdk.ValAddress([]byte("validator1"))}
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
func TestGetSetPendingRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)

	keeper.SetPendingRequests(ctx, []int64{1, 2, 3})

	reqIDs = keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{1, 2, 3}, reqIDs)

	keeper.SetPendingRequests(ctx, []int64{})
	reqIDs = keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)
}

// Can add new pending request if request doesn't exist in list,
// and return error if request has already existed in list.
func TestAddPendingRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{}, reqIDs)

	keeper.SetPendingRequests(ctx, []int64{1, 2})
	err := keeper.AddPendingRequest(ctx, 3)
	require.Nil(t, err)
	reqIDs = keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{1, 2, 3}, reqIDs)

	err = keeper.AddPendingRequest(ctx, 3)
	require.Equal(t, types.CodeDuplicateRequest, err.Code())
	reqIDs = keeper.GetPendingRequests(ctx)
	require.Equal(t, []int64{1, 2, 3}, reqIDs)
}

func TestHasToPutInPendingList(t *testing.T) {
	// TODO: Write test
}
