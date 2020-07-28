package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Context struct {
	client    rpcclient.Client
	gasPrices sdk.DecCoins
	keys      chan keys.Info
	amount    sdk.Coins
}
