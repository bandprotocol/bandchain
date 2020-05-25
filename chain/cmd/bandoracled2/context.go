package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Context struct {
	client    *rpchttp.HTTP
	chainRestServerURI string
	validator sdk.ValAddress
	gasPrices sdk.DecCoins
	keys      chan keyring.Info
	executor  executor
}
