package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Context struct {
	// Run context
	key    keyring.Info
	client *rpchttp.HTTP
}
