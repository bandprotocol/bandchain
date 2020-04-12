package keeper

import (
	"testing"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	_, err := keeper.AddRequest(ctx, 1, calldata, 2, 2, "clientID")
	require.NotNil(t, err)

	script := GetTestOracleScript("../../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)
	_, err = keeper.AddRequest(ctx, 1, calldata, 2, 2, "clientID")
	require.NotNil(t, err)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := SetupTestValidator(
		ctx,
		keeper,
		pubStr[0],
		10,
	)
	_, err = keeper.AddRequest(ctx, 1, calldata, 2, 2, "clientID")
	require.NotNil(t, err)

	validatorAddress2 := SetupTestValidator(
		ctx,
		keeper,
		pubStr[1],
		100,
	)
	requestID, err := keeper.AddRequest(ctx, 1, calldata, 2, 2, "clientID")
	require.Nil(t, err)
	require.Equal(t, types.RequestID(1), requestID)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	expectRequest := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 22, "clientID",
	)
	require.Equal(t, expectRequest, actualRequest)
}

func TestRequestCallDataSizeTooBig(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	script := GetTestOracleScript("../../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	SetupTestValidator(
		ctx,
		keeper,
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		10,
	)
	SetupTestValidator(
		ctx,
		keeper,
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
		100,
	)

	// Set MaxCalldataSize to 0
	keeper.SetParam(ctx, types.KeyMaxCalldataSize, 0)
	// Should fail because size of "calldata" is > 0
	_, err := keeper.AddRequest(ctx, 1, []byte("calldata"), 2, 2, "clientID")
	require.NotNil(t, err)

	// Set MaxCalldataSize to 20
	keeper.SetParam(ctx, types.KeyMaxCalldataSize, 20)
	// Should pass because size of "calldata" is < 20
	_, err = keeper.AddRequest(ctx, 1, []byte("calldata"), 2, 2, "clientID")
	require.Nil(t, err)
}

func TestRequestExceedEndBlockExecuteGasLimit(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	script := GetTestOracleScript("../../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	SetupTestValidator(
		ctx,
		keeper,
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		10,
	)
	SetupTestValidator(
		ctx,
		keeper,
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
		100,
	)

	// Set EndBlockExecuteGasLimit to 10000
	keeper.SetParam(ctx, types.KeyEndBlockExecuteGasLimit, 10000)
	// Should fail because required execute gas is > 10000
	_, err := keeper.AddRequest(ctx, 1, []byte("calldata"), 2, 2, "clientID")
	require.NotNil(t, err)

	// Set EndBlockExecuteGasLimit to 120000
	keeper.SetParam(ctx, types.KeyEndBlockExecuteGasLimit, 120000)
	// Should fail because required execute gas is < 120000
	_, err = keeper.AddRequest(ctx, 1, []byte("calldata"), 2, 2, "clientID")
	require.Nil(t, err)
}

// TestSetResolved tests keeper can set resolved status to request
func TestSetResolved(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)

	err := keeper.SetResolve(ctx, 1, types.Success)
	require.Nil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	request.ResolveStatus = types.Success
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// TestSetResolvedOnInvalidRequest tests keeper must return if set on invalid request
func TestSetResolvedOnInvalidRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	keeper.SetRequest(ctx, 1, request)
	err := keeper.SetResolve(ctx, 2, types.Success)
	require.NotNil(t, err)
}

// Can get/set pending request correctly and set empty case
func TestGetSetPendingRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPendingResolveList(ctx)

	require.Equal(t, []types.RequestID{}, reqIDs)

	keeper.SetPendingResolveList(ctx, []types.RequestID{1, 2, 3})

	reqIDs = keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1, 2, 3}, reqIDs)

	keeper.SetPendingResolveList(ctx, []types.RequestID{})
	reqIDs = keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, reqIDs)
}

// Can add new pending request if request doesn't exist in list,
// and return error if request has already existed in list.
func TestAddPendingRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPendingResolveList(ctx)

	require.Equal(t, []types.RequestID{}, reqIDs)

	keeper.SetPendingResolveList(ctx, []types.RequestID{1, 2})
	err := keeper.AddPendingRequest(ctx, 3)
	require.Nil(t, err)
	reqIDs = keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1, 2, 3}, reqIDs)

	err = keeper.AddPendingRequest(ctx, 3)
	require.NotNil(t, err)
	reqIDs = keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1, 2, 3}, reqIDs)
}

func TestHasToPutInPendingList(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	require.False(t, keeper.ShouldBecomePendingResolve(ctx, 1))
	request := newDefaultRequest()
	request.SufficientValidatorCount = 1
	keeper.SetRequest(ctx, 1, request)
	require.False(t, keeper.ShouldBecomePendingResolve(ctx, 1))

	err := keeper.AddReport(ctx, 1, []types.RawDataReportWithID{}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
	require.Nil(t, err)
	require.True(t, keeper.ShouldBecomePendingResolve(ctx, 1))

	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{}, sdk.ValAddress([]byte("validator2")), sdk.AccAddress([]byte("validator2")))
	require.Nil(t, err)
	require.False(t, keeper.ShouldBecomePendingResolve(ctx, 1))
}

func TestValidateDataSourceCount(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	// Set MaxDataSourceCountPerRequest to 3
	keeper.SetParam(ctx, types.KeyMaxDataSourceCountPerRequest, 3)

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	keeper.SetRawDataRequest(ctx, 1, 101, types.NewRawDataRequest(0, []byte("calldata1")))
	err := keeper.ValidateDataSourceCount(ctx, 1)
	require.Nil(t, err)

	keeper.SetRawDataRequest(ctx, 1, 102, types.NewRawDataRequest(0, []byte("calldata2")))
	err = keeper.ValidateDataSourceCount(ctx, 1)
	require.Nil(t, err)

	keeper.SetRawDataRequest(ctx, 1, 103, types.NewRawDataRequest(0, []byte("calldata3")))
	err = keeper.ValidateDataSourceCount(ctx, 1)
	require.Nil(t, err)

	// Validation of "104" will return an error because MaxDataSourceCountPerRequest was set to 3.
	keeper.SetRawDataRequest(ctx, 1, 104, types.NewRawDataRequest(0, []byte("calldata4")))
	err = keeper.ValidateDataSourceCount(ctx, 1)
	require.NotNil(t, err)
}

func TestPayDataSourceFees(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	sender := sdk.AccAddress([]byte("sender"))
	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, NewUBandCoins(100))
	require.Nil(t, err)

	owner1 := sdk.AccAddress([]byte("owner1"))
	owner2 := sdk.AccAddress([]byte("owner2"))

	dataSource1 := types.NewDataSource(
		sdk.AccAddress([]byte("owner1")),
		"data_source1",
		"description1",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 40)),
		[]byte("executable1"),
	)
	keeper.SetDataSource(ctx, 1, dataSource1)

	dataSource2 := types.NewDataSource(
		sdk.AccAddress([]byte("owner2")),
		"data_source2",
		"description2",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 50)),
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(1, []byte("calldata")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(2, []byte("calldata")))

	err = keeper.PayDataSourceFees(ctx, 1, sender)
	require.Nil(t, err)

	balance := keeper.CoinKeeper.GetAllBalances(ctx, sender)
	require.Equal(t, NewUBandCoins(10), balance)

	owner1Balance := keeper.CoinKeeper.GetAllBalances(ctx, owner1)
	require.Equal(t, NewUBandCoins(40), owner1Balance)

	owner2Balance := keeper.CoinKeeper.GetAllBalances(ctx, owner2)
	require.Equal(t, NewUBandCoins(50), owner2Balance)
}
