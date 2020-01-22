package zoracle

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bandprotocol/d3n/chain/wasm"
	keep "github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/common"
)

func setupTestValidator(ctx sdk.Context, keeper Keeper) sdk.ValAddress {
	pubKey := keep.NewPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	validatorAddress := sdk.ValAddress(pubKey.Address())
	initTokens := sdk.TokensFromConsensusPower(10)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
	keeper.CoinKeeper.AddCoins(ctx, sdk.AccAddress(pubKey.Address()), initCoins)

	msgCreateValidator := staking.NewTestMsgCreateValidator(
		validatorAddress, pubKey, sdk.TokensFromConsensusPower(10),
	)
	stakingHandler := staking.NewHandler(keeper.StakingKeeper)
	stakingHandler(ctx, msgCreateValidator)

	keeper.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	return validatorAddress
}

func TestRequestSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
	code, err := wasm.ReadBytes(absPath)
	if err != nil {
		fmt.Println(err)
	}
	name := "Crypto price"
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, code, name, sender)
	params, _ := hex.DecodeString("00000000")
	msg := types.NewMsgRequest(codeHash, params, 5, sender)
	got := handleMsgRequest(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)

	// Check global request count
	require.Equal(t, uint64(1), keeper.GetRequestCount(ctx))
	request, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)

	// Check codeHash must match
	require.Equal(t, codeHash, request.CodeHash)

	// Check reportEndAt
	require.Equal(t, uint64(ctx.BlockHeight()+5), request.ReportEndAt)

	// Check pending request list
	require.Equal(t, []uint64{1}, keeper.GetPending(ctx))

	// check event
	require.Equal(t, types.EventTypeRequest, ctx.EventManager().Events()[0].Type)

	// check codeHash, prepare attribute
	codeHashPair := common.KVPair{
		Key:   []byte(types.AttributeKeyCodeHash),
		Value: []byte(hex.EncodeToString(codeHash)),
	}
	namePair := common.KVPair{
		Key:   []byte(types.AttributeKeyCodeName),
		Value: []byte("Crypto price"),
	}
	preparePair := common.KVPair{
		Key:   []byte(types.AttributeKeyPrepare),
		Value: []byte("5b7b22636d64223a226375726c222c2261726773223a5b2268747470733a2f2f6170692e636f696e6765636b6f2e636f6d2f6170692f76332f73696d706c652f70726963653f6964733d626974636f696e2676735f63757272656e636965733d757364225d7d2c7b22636d64223a226375726c222c2261726773223a5b2268747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d425443267473796d733d555344225d7d5d"),
	}
	require.Equal(t, codeHashPair, ctx.EventManager().Events()[0].Attributes[1])
	require.Equal(t, namePair, ctx.EventManager().Events()[0].Attributes[2])
	require.Equal(t, preparePair, ctx.EventManager().Events()[0].Attributes[3])
}

func TestRequestInvalidCodeHash(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	sender := sdk.AccAddress([]byte("sender"))

	codeHash, _ := hex.DecodeString("c0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0dec0de")
	params, _ := hex.DecodeString("00000000")
	msg := types.NewMsgRequest(codeHash, params, 5, sender)
	got := handleMsgRequest(ctx, keeper, msg)
	require.False(t, got.IsOK(), "expected request is an invalid tx")
	require.Equal(t, types.CodeInvalidInput, got.Code)
}

func TestRequestInvalidWasmCode(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, []byte("Fake code"), "Fake script", sender)
	params, _ := hex.DecodeString("00000000")
	msg := types.NewMsgRequest(codeHash, params, 5, sender)
	got := handleMsgRequest(ctx, keeper, msg)
	require.False(t, got.IsOK(), "expected request is an invalid tx")
	require.Equal(t, types.WasmError, got.Code)
}

func TestReportSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	validatorAddress := setupTestValidator(ctx, keeper)

	// set request = 2
	name := "Script1"
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
	request := types.NewRequest(codeHash, []byte("params"), 3)
	keeper.SetRequest(ctx, 2, request)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 2)
	keeper.SetPending(ctx, pendingRequests)

	// set blockheight
	ctx = ctx.WithBlockHeight(3)

	// report data
	msg := types.NewMsgReport(2, []byte("data"), validatorAddress)
	got := handleMsgReport(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set report to be ok, got %v", got)

	//check event
	require.Equal(t, types.EventTypeReport, ctx.EventManager().Events()[2].Type)
}

func TestReportInvalidValidator(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	pubKey := keep.NewPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	validatorAddress := sdk.ValAddress(pubKey.Address())

	// set request = 2
	name := "Script1"
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
	request := types.NewRequest(codeHash, []byte("params"), 3)
	keeper.SetRequest(ctx, 1, request)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 1)
	keeper.SetPending(ctx, pendingRequests)

	msg := types.NewMsgReport(1, []byte("data"), validatorAddress)
	got := handleMsgReport(ctx, keeper, msg)
	require.Equal(t, types.CodeInvalidValidator, got.Code)
}

func TestOutOfReportPeriod(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	validatorAddress := setupTestValidator(ctx, keeper)

	// set request = 2
	name := "Script1"
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, []byte("Code"), name, sender)
	request := types.NewRequest(codeHash, []byte("params"), 3)
	keeper.SetRequest(ctx, 2, request)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 2)
	keeper.SetPending(ctx, pendingRequests)

	// set blockheight
	ctx = ctx.WithBlockHeight(10)

	// report data
	msg := types.NewMsgReport(2, []byte("data"), validatorAddress)
	got := handleMsgReport(ctx, keeper, msg)
	require.Equal(t, types.CodeOutOfReportPeriod, got.Code)
}

func TestStoreCodeSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
	code, err := wasm.ReadBytes(absPath)
	if err != nil {
		fmt.Println(err)
	}
	name := "Crypto price"
	owner := sdk.AccAddress([]byte("owner"))
	codeHash := types.NewStoredCode(code, name, owner).GetCodeHash()

	msg := types.NewMsgStoreCode(code, name, owner)
	got := handleMsgStoreCode(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected store code to be ok, got %v", got)

	// Check codehash from event
	require.Equal(t, types.EventTypeStoreCode, got.Events.ToABCIEvents()[0].Type)
	require.Equal(t, hex.EncodeToString(codeHash), string(got.Events.ToABCIEvents()[0].Attributes[0].Value))
	require.Equal(t, name, string(got.Events.ToABCIEvents()[0].Attributes[1].Value))

	// Check value in store
	sc, err := keeper.GetCode(ctx, codeHash)
	require.Nil(t, err)

	require.Equal(t, owner, sc.Owner)
	require.Equal(t, code, sc.Code)
}

func TestStoreCodeFailed(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	code := []byte("Code")
	name := "Failed script"
	owner := sdk.AccAddress([]byte("owner"))

	keeper.SetCode(ctx, code, name, owner)

	msg := types.NewMsgStoreCode(code, name, owner)
	got := handleMsgStoreCode(ctx, keeper, msg)

	require.False(t, got.IsOK())

	require.Equal(t, types.CodeInvalidInput, got.Code)
}

func TestDeleteCodeSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	code := []byte("Code")
	name := "script"
	owner := sdk.AccAddress([]byte("owner"))
	codeHash := keeper.SetCode(ctx, code, name, owner)

	msg := types.NewMsgDeleteCode(codeHash, owner)
	got := handleMsgDeleteCode(ctx, keeper, msg)

	require.Equal(t, types.EventTypeDeleteCode, got.Events.ToABCIEvents()[0].Type)
	require.Equal(t, hex.EncodeToString(codeHash), string(got.Events.ToABCIEvents()[0].Attributes[0].Value))
}

func TestDeleteCodeInvalidHash(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	code := []byte("Code")
	name := "script"
	owner := sdk.AccAddress([]byte("owner"))
	codeHash := keeper.SetCode(ctx, code, name, owner)
	invalidCodeHash := append(codeHash[:31], byte('b'))

	msg := types.NewMsgDeleteCode(invalidCodeHash, owner)
	got := handleMsgDeleteCode(ctx, keeper, msg)

	require.False(t, got.IsOK())

	require.Equal(t, types.CodeInvalidInput, got.Code)
}

func TestDeleteCodeInvalidOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	code := []byte("Code")
	name := "script"
	owner := sdk.AccAddress([]byte("owner"))
	codeHash := keeper.SetCode(ctx, code, name, owner)

	other := sdk.AccAddress([]byte("other"))
	msg := types.NewMsgDeleteCode(codeHash, other)
	got := handleMsgDeleteCode(ctx, keeper, msg)

	require.False(t, got.IsOK())

	require.Equal(t, types.CodeInvalidOwner, got.Code)
}

func TestEndBlock(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	ctx = ctx.WithBlockHeight(0)
	absPath, _ := filepath.Abs("../../wasm/res/result.wasm")
	code, err := wasm.ReadBytes(absPath)
	if err != nil {
		fmt.Println(err)
	}
	name := "Crypto price"
	sender := sdk.AccAddress([]byte("sender"))
	codeHash := keeper.SetCode(ctx, code, name, sender)

	params, _ := hex.DecodeString("00000000")
	// set request
	request := types.NewRequest(codeHash, params, 3)
	keeper.SetRequest(ctx, 1, request)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 1)
	keeper.SetPending(ctx, pendingRequests)

	// set report
	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

	data1, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")
	data2, _ := hex.DecodeString("5b227b5c22626974636f696e5c223a7b5c227573645c223a373139342e32357d7d222c227b5c225553445c223a373231342e31327d225d")

	keeper.SetReport(ctx, 1, validatorAddress1, data1)
	keeper.SetReport(ctx, 1, validatorAddress2, data2)

	// blockheight update to 2
	ctx = ctx.WithBlockHeight(2)

	gotEndBlock := handleEndBlock(ctx, keeper)
	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

	_, err = keeper.GetResult(ctx, 1, codeHash, params)
	// Result must not found
	require.NotNil(t, err)

	pendingRequests = keeper.GetPending(ctx)
	require.Equal(t, []uint64{1}, pendingRequests)

	// blockheight update to 4
	ctx = ctx.WithBlockHeight(4)
	resultAfter, _ := hex.DecodeString("00000000000afe22")

	// handle end block
	gotEndBlock = handleEndBlock(ctx, keeper)
	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

	request, _ = keeper.GetRequest(ctx, 1)

	result, err := keeper.GetResult(ctx, 1, codeHash, params)
	require.Nil(t, err)
	require.Equal(t, resultAfter, result)

	pendingRequests = keeper.GetPending(ctx)
	require.Equal(t, []uint64{}, pendingRequests)
}
