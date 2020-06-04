package oracle_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// func TestEndBlock(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)

// 	ctx = ctx.WithBlockHeight(2)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// 	calldata := []byte("calldata")
// 	sender := sdk.AccAddress([]byte("sender"))

// 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// 	keeper.SetOracleScript(ctx, 1, script)

// 	pubStr := []string{
// 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// 	}

// 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// 	dataSource := keep.GetTestDataSource()
// 	keeper.SetDataSource(ctx, 1, dataSource)

// 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// 	handleMsgRequestData(ctx, keeper, msg)

// 	keeper.SetReport(ctx, 1, 1, validatorAddress1, types.NewReport(0, []byte("answer1")))
// 	keeper.SetReport(ctx, 1, 1, validatorAddress2, types.NewReport(0, []byte("answer2")))

// 	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))
// 	handleEndBlock(ctx, keeper)

// 	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

// 	result, err := keeper.GetResult(ctx, 1, 1, calldata)
// 	require.Nil(t, err)
// 	require.Equal(t,
// 		types.Result{
// 			RequestTime:              1581589790,
// 			AggregationTime:          1581589999,
// 			RequestedValidatorsCount: 2,
// 			MinCount: 2,
// 			ReportedValidatorsCount:  0,
// 			Data:                     []byte("answer2"),
// 		},
// 		result,
// 	)

// 	actualRequest, err := keeper.GetRequest(ctx, 1)
// 	require.Nil(t, err)
// 	require.Equal(t, types.ResolveStatus_Success, actualRequest.ResolveStatus)
// }

func TestRequestOracleData(t *testing.T) {
	app, ctx, k := createTestInput()

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	for i := 1; i <= 3; i++ {
		// dataSource, clear := d
		// k.SetDataSource(ctx, i,)
	}
	absPath, _ := filepath.Abs("../../pkg/owasm/res/beeb.wasm")
	code, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	createOracleScript := types.NewMsgCreateOracleScript(
		Owner.Address, "Oracle script", "description", code, "schema", "sourceCodeURL", Owner.Address,
	)
	_, _, err = simapp.SignCheckDeliver(
		t, app.Codec(), app.BaseApp, header, []sdk.Msg{createOracleScript}, []uint64{0},
		[]uint64{uint64(4)}, true, true, Owner.PrivKey,
	)
	require.NoError(t, err)
	requestMsg := types.NewMsgRequestData(types.OracleScriptID(1), []byte("calldata"), 3, 2, "app_test", Alice.Address)
	_, _, err = simapp.SignCheckDeliver(
		t, app.Codec(), app.BaseApp, header, []sdk.Msg{requestMsg}, []uint64{0},
		[]uint64{uint64(5)}, true, true, Alice.PrivKey,
	)
	require.NoError(t, err)
	// simapp.SignCheckDeliver(t, app.Codec(), app, header)
}
