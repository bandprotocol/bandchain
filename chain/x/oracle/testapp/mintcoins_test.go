package testapp

import (
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func TestMintCoins(t *testing.T) {
	app := NewSimApp("test-chain", log.NewNopLogger())
	ctx := app.NewContext(false, abci.Header{})
	require.NoError(t, app.DistrKeeper.FundCommunityPool(ctx, Coins1000000uband, Owner.Address))
	require.NoError(t, app.MintKeeper.MintCoins(ctx, Coins1000000uband))
}
