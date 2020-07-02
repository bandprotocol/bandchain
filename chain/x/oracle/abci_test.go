package oracle_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
)

func fromHex(hexStr string) []byte {
	res, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return res
}

func TestRollingSeedCorrect(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	// Initially rolling seed should be all zeros.
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000000"), k.GetRollingSeed(ctx))
	// Every begin block, the rolling seed should get updated.
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("0100000000000000000000000000000000000000000000000000000000000000"),
	})
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000001"), k.GetRollingSeed(ctx))
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("0200000000000000000000000000000000000000000000000000000000000000"),
	})
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000102"), k.GetRollingSeed(ctx))
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	require.Equal(t, fromHex("00000000000000000000000000000000000000000000000000000000000102ff"), k.GetRollingSeed(ctx))
}
