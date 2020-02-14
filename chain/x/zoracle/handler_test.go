package zoracle

import (
	"testing"
	"time"

	keep "github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRequestSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, sender)

	// Test here
	got := handleMsgRequest(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)

	// Check global request count
	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	expectRequest := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102,
	)
	require.Equal(t, expectRequest, actualRequest)

	require.Equal(t, int64(1), keeper.GetRawDataRequestCount(ctx, 1))

	rawRequests := []types.RawDataRequest{
		types.NewRawDataRequest(1, []byte("band-protocol")),
	}
	require.Equal(t, rawRequests, keeper.GetRawDataRequests(ctx, 1))
}

func TestRequestInvalidDataSource(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, sender)
	got := handleMsgRequest(ctx, keeper, msg)
	require.False(t, got.IsOK())

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	got = handleMsgRequest(ctx, keeper, msg)
	require.False(t, got.IsOK())
}

// func TestRequestInvalidCodeHash(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	sender := sdk.AccAddress([]byte("sender"))

// 	codeHash, _ := hex.DecodeString("c0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0de")
// 	params, _ := hex.DecodeString("00000000")
// 	msg := types.NewMsgRequest(codeHash, params, 5, sender)
// 	got := handleMsgRequest(ctx, keeper, msg)
// 	require.False(t, got.IsOK(), "expected request is an invalid tx")
// 	require.Equal(t, types.CodeInvalidInput, got.Code)
// }

// func TestRequestInvalidWasmCode(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, []byte("Fake code"), "Fake script", sender)
// 	params, _ := hex.DecodeString("00000000")
// 	msg := types.NewMsgRequest(codeHash, params, 5, sender)
// 	got := handleMsgRequest(ctx, keeper, msg)
// 	require.False(t, got.IsOK(), "expected request is an invalid tx")
// 	require.Equal(t, types.WasmError, got.Code)
// }

// func TestReportSuccess(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	validatorAddress := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50",
// 	)

// 	// set request = 2
// 	name := "Script1"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
// 	request := types.NewRequest(codeHash, []byte("params"), 3)
// 	keeper.SetRequest(ctx, 2, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 2)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	// set blockheight
// 	ctx = ctx.WithBlockHeight(3)

// 	// report data
// 	msg := types.NewMsgReport(2, []byte("data"), validatorAddress)
// 	got := handleMsgReport(ctx, keeper, msg)
// 	require.True(t, got.IsOK(), "expected set report to be ok, got %v", got)

// 	//check event
// 	require.Equal(t, types.EventTypeReport, ctx.EventManager().Events()[2].Type)
// }

// func TestReportInvalidValidator(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	pubKey := keep.NewPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
// 	validatorAddress := sdk.ValAddress(pubKey.Address())

// 	// set request = 2
// 	name := "Script1"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
// 	request := types.NewRequest(codeHash, []byte("params"), 3)
// 	keeper.SetRequest(ctx, 1, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 1)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	msg := types.NewMsgReport(1, []byte("data"), validatorAddress)
// 	got := handleMsgReport(ctx, keeper, msg)
// 	require.Equal(t, types.CodeInvalidValidator, got.Code)
// }

// func TestOutOfReportPeriod(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	validatorAddress := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50",
// 	)

// 	// set request = 2
// 	name := "Script1"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
// 	request := types.NewRequest(codeHash, []byte("params"), 3)
// 	keeper.SetRequest(ctx, 2, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 2)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	// set blockheight
// 	ctx = ctx.WithBlockHeight(10)

// 	// report data
// 	msg := types.NewMsgReport(2, []byte("data"), validatorAddress)
// 	got := handleMsgReport(ctx, keeper, msg)
// 	require.Equal(t, types.CodeOutOfReportPeriod, got.Code)
// }

// TODO: Left this code as reference code after implemented handle store oracle script please remove these.
// func TestStoreCodeSuccess(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
// 	code, err := wasm.ReadBytes(absPath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	name := "Crypto price"
// 	owner := sdk.AccAddress([]byte("owner"))
// 	codeHash := types.NewStoredCode(code, name, owner).GetCodeHash()

// 	msg := types.NewMsgStoreCode(code, name, owner)
// 	got := handleMsgStoreCode(ctx, keeper, msg)
// 	require.True(t, got.IsOK(), "expected store code to be ok, got %v", got)

// 	// Check codehash from event
// 	require.Equal(t, types.EventTypeStoreCode, got.Events.ToABCIEvents()[0].Type)
// 	require.Equal(t, hex.EncodeToString(codeHash), string(got.Events.ToABCIEvents()[0].Attributes[0].Value))
// 	require.Equal(t, name, string(got.Events.ToABCIEvents()[0].Attributes[1].Value))

// 	// Check value in store
// 	sc, err := keeper.GetCode(ctx, codeHash)
// 	require.Nil(t, err)

// 	require.Equal(t, owner, sc.Owner)
// 	require.Equal(t, code, []byte(sc.Code))
// }

// func TestStoreCodeFailed(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	code := []byte("Code")
// 	name := "Failed script"
// 	owner := sdk.AccAddress([]byte("owner"))

// 	keeper.SetCode(ctx, code, name, owner)

// 	msg := types.NewMsgStoreCode(code, name, owner)
// 	got := handleMsgStoreCode(ctx, keeper, msg)

// 	require.False(t, got.IsOK())

// 	require.Equal(t, types.CodeInvalidInput, got.Code)
// }

// func TestDeleteCodeSuccess(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	code := []byte("Code")
// 	name := "script"
// 	owner := sdk.AccAddress([]byte("owner"))
// 	codeHash := keeper.SetCode(ctx, code, name, owner)

// 	msg := types.NewMsgDeleteCode(codeHash, owner)
// 	got := handleMsgDeleteCode(ctx, keeper, msg)

// 	require.Equal(t, types.EventTypeDeleteCode, got.Events.ToABCIEvents()[0].Type)
// 	require.Equal(t, hex.EncodeToString(codeHash), string(got.Events.ToABCIEvents()[0].Attributes[0].Value))
// }

// func TestDeleteCodeInvalidHash(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	code := []byte("Code")
// 	name := "script"
// 	owner := sdk.AccAddress([]byte("owner"))
// 	codeHash := keeper.SetCode(ctx, code, name, owner)
// 	invalidCodeHash := append(codeHash[:31], byte('b'))

// 	msg := types.NewMsgDeleteCode(invalidCodeHash, owner)
// 	got := handleMsgDeleteCode(ctx, keeper, msg)

// 	require.False(t, got.IsOK())

// 	require.Equal(t, types.CodeInvalidInput, got.Code)
// }

// func TestDeleteCodeInvalidOwner(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	code := []byte("Code")
// 	name := "script"
// 	owner := sdk.AccAddress([]byte("owner"))
// 	codeHash := keeper.SetCode(ctx, code, name, owner)

// 	other := sdk.AccAddress([]byte("other"))
// 	msg := types.NewMsgDeleteCode(codeHash, other)
// 	got := handleMsgDeleteCode(ctx, keeper, msg)

// 	require.False(t, got.IsOK())

// 	require.Equal(t, types.CodeInvalidOwner, got.Code)
// }

// func TestEndBlock(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	ctx = ctx.WithBlockHeight(0)

// 	validatorAddress1 := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50",
// 	)
// 	validatorAddress2 := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51",
// 	)

// 	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
// 	code, err := wasm.ReadBytes(absPath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	name := "Crypto price"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, code, name, sender)

// 	params, _ := hex.DecodeString("00000000")
// 	// set request
// 	request := types.NewRequest(codeHash, params, 3)
// 	keeper.SetRequest(ctx, 1, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 1)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	data1, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")
// 	data2, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139312e32357d7d222c227b5c225553445c223a373230392e31357d225d")

// 	keeper.SetReport(ctx, 1, validatorAddress1, data1)

// 	// blockheight update to 2
// 	ctx = ctx.WithBlockHeight(2)

// 	gotEndBlock := handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	_, err = keeper.GetResult(ctx, 1, codeHash, params)
// 	// Result must not found
// 	require.NotNil(t, err)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{1}, pendingRequests)

// 	// blockheight update to 4
// 	keeper.SetReport(ctx, 1, validatorAddress2, data2)
// 	ctx = ctx.WithBlockHeight(4)
// 	resultAfter, _ := hex.DecodeString("00000000000afd5b")

// 	// handle end block
// 	gotEndBlock = handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	request, _ = keeper.GetRequest(ctx, 1)

// 	result, err := keeper.GetResult(ctx, 1, codeHash, params)
// 	require.Nil(t, err)
// 	require.Equal(t, resultAfter, result)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{}, pendingRequests)
// }

// func TestEndBlockQuickResolve(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	ctx = ctx.WithBlockHeight(0)

// 	validatorAddress1 := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50",
// 	)
// 	validatorAddress2 := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51",
// 	)

// 	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
// 	code, err := wasm.ReadBytes(absPath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	name := "Crypto price"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, code, name, sender)

// 	params, _ := hex.DecodeString("00000000")
// 	// set request
// 	request := types.NewRequest(codeHash, params, 10000)
// 	keeper.SetRequest(ctx, 1, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 1)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	data1, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")
// 	data2, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139312e32357d7d222c227b5c225553445c223a373230392e31357d225d")

// 	keeper.SetReport(ctx, 1, validatorAddress1, data1)

// 	// blockheight update to 100
// 	ctx = ctx.WithBlockHeight(100)

// 	gotEndBlock := handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	_, err = keeper.GetResult(ctx, 1, codeHash, params)
// 	// Result must not found
// 	require.NotNil(t, err)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{1}, pendingRequests)

// 	keeper.SetReport(ctx, 1, validatorAddress2, data2)
// 	// blockheight update to 300
// 	ctx = ctx.WithBlockHeight(300)
// 	resultAfter, _ := hex.DecodeString("00000000000afd5b")

// 	// handle end block should apply quick resolve
// 	gotEndBlock = handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	request, _ = keeper.GetRequest(ctx, 1)

// 	result, err := keeper.GetResult(ctx, 1, codeHash, params)
// 	require.Nil(t, err)
// 	require.Equal(t, resultAfter, result)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{}, pendingRequests)
// }

// func TestEndBlockReportEnd(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	ctx = ctx.WithBlockHeight(0)

// 	validatorAddress1 := setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50",
// 	)
// 	setupTestValidator(
// 		ctx,
// 		keeper,
// 		"0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51",
// 	)

// 	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
// 	code, err := wasm.ReadBytes(absPath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	name := "Crypto price"
// 	sender := sdk.AccAddress([]byte("sender"))
// 	codeHash := keeper.SetCode(ctx, code, name, sender)

// 	params, _ := hex.DecodeString("00000000")
// 	// set request
// 	request := types.NewRequest(codeHash, params, 300)
// 	keeper.SetRequest(ctx, 1, request)

// 	// set pending
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 1)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

// 	data1, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")

// 	keeper.SetReport(ctx, 1, validatorAddress1, data1)

// 	// blockheight update to 100
// 	ctx = ctx.WithBlockHeight(100)

// 	gotEndBlock := handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	_, err = keeper.GetResult(ctx, 1, codeHash, params)
// 	// Result must not found
// 	require.NotNil(t, err)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{1}, pendingRequests)

// 	// blockheight update to 300
// 	ctx = ctx.WithBlockHeight(300)
// 	resultAfter, _ := hex.DecodeString("00000000000afe22")

// 	// handle end block after report end
// 	gotEndBlock = handleEndBlock(ctx, keeper)
// 	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

// 	request, _ = keeper.GetRequest(ctx, 1)

// 	result, err := keeper.GetResult(ctx, 1, codeHash, params)
// 	require.Nil(t, err)
// 	require.Equal(t, resultAfter, result)

// 	pendingRequests = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []uint64{}, pendingRequests)
// }

func mockDataSource(ctx sdk.Context, keeper Keeper) sdk.Result {
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source_1"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")
	sender := sdk.AccAddress([]byte("sender"))
	msg := types.NewMsgCreateDataSource(owner, name, fee, executable, sender)
	return handleMsgCreateDataSource(ctx, keeper, msg)
}

func TestCreateDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	got := mockDataSource(ctx, keeper)
	require.True(t, got.IsOK(), "expected set data source to be ok, got %v", got)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), dataSource.Owner)
	require.Equal(t, "data_source_1", dataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), dataSource.Fee)
	require.Equal(t, []byte("executable"), dataSource.Executable)
}

func TestEditDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newFee, newExecutable, sender)
	got := handleMsgEditDataSource(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected edit data source to be ok, got %v", got)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	require.Equal(t, newFee, dataSource.Fee)
	require.Equal(t, newExecutable, dataSource.Executable)
}

func TestEditDataSourceByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newFee, newExecutable, sender)
	got := handleMsgEditDataSource(ctx, keeper, msg)
	require.False(t, got.IsOK())
	require.Equal(t, types.CodeInvalidOwner, got.Code)
}
