package oracle

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bandprotocol/bandx/oracle/wasm"
	keep "github.com/bandprotocol/bandx/oracle/x/oracle/internal/keeper"
	"github.com/bandprotocol/bandx/oracle/x/oracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/common"
)

func TestRequestSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	absPath, _ := filepath.Abs("../../wasm/res/test.wasm")
	code, err := wasm.ReadBytes(absPath)
	if err != nil {
		fmt.Println(err)
	}
	codeHash := crypto.Keccak256(code)
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgRequest(code, 5, sender)
	got := handleMsgRequest(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set request(datapoint) to be ok, got %v", got)

	// Check global request count
	require.Equal(t, uint64(1), keeper.GetRequestCount(ctx))
	request, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)

	// Check request id must be 1
	require.Equal(t, uint64(1), request.RequestID)

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
	preparePair := common.KVPair{
		Key:   []byte(types.AttributeKeyPrepare),
		Value: []byte("020000000000000004000000000000006375726c01000000000000004b0000000000000068747470733a2f2f6170692e636f696e6765636b6f2e636f6d2f6170692f76332f73696d706c652f70726963653f6964733d626974636f696e2676735f63757272656e636965733d75736404000000000000006375726c01000000000000003f0000000000000068747470733a2f2f6d696e2d6170692e63727970746f636f6d706172652e636f6d2f646174612f70726963653f6673796d3d425443267473796d733d555344"),
	}
	require.Equal(t, codeHashPair, ctx.EventManager().Events()[0].Attributes[1])
	require.Equal(t, preparePair, ctx.EventManager().Events()[0].Attributes[2])
}

func TestReportSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
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

	updates := keeper.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.Equal(t, 1, len(updates))

	msg := types.NewMsgReport(1, []byte("data"), validatorAddress)
	got := handleMsgReport(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected set report to be ok, got %v", got)

	//check event
	require.Equal(t, types.EventTypeReport, ctx.EventManager().Events()[2].Type)
}

func TestReportInvalidValidator(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	pubKey := keep.NewPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	validatorAddress := sdk.ValAddress(pubKey.Address())

	msg := types.NewMsgReport(1, []byte("data"), validatorAddress)
	got := handleMsgReport(ctx, keeper, msg)
	require.Equal(t, types.CodeInvalidValidator, got.Code)
}

func TestEndBlock(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	ctx = ctx.WithBlockHeight(0)
	absPath, _ := filepath.Abs("../../wasm/res/test.wasm")
	code, err := wasm.ReadBytes(absPath)
	if err != nil {
		fmt.Println(err)
	}
	codeHash := keeper.SetCode(ctx, code)

	// set request
	datapoint := types.NewDataPoint(1, codeHash, 3)
	keeper.SetRequest(ctx, 1, datapoint)

	// set pending
	pendingRequests := keeper.GetPending(ctx)
	pendingRequests = append(pendingRequests, 1)
	keeper.SetPending(ctx, pendingRequests)

	// set report
	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

	data1, _ := hex.DecodeString("02000000000000001b000000000000007b22626974636f696e223a7b22757364223a373436392e34397d7d0f000000000000007b22555344223a373531302e32317d")
	data2, _ := hex.DecodeString("02000000000000001b000000000000007b22626974636f696e223a7b22757364223a373436392e34397d7d0f000000000000007b22555344223a373532312e38317d")

	keeper.SetReport(ctx, 1, validatorAddress1, data1)
	keeper.SetReport(ctx, 1, validatorAddress2, data2)

	// blockheight update to 4
	ctx = ctx.WithBlockHeight(4)
	resultAfter, _ := hex.DecodeString("00000000000b6edb")

	// handle end block
	gotEndBlock := handleEndBlock(ctx, keeper)
	require.True(t, gotEndBlock.IsOK(), "expected end block to be ok, got %v", gotEndBlock)

	request, _ := keeper.GetRequest(ctx, 1)
	require.Equal(t, uint64(1), request.RequestID)
	require.Equal(t, resultAfter, request.Result)

	pendingRequests = keeper.GetPending(ctx)
	require.Equal(t, []uint64{}, pendingRequests)
}
