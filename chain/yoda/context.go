package yoda

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
)

type Context struct {
	client    rpcclient.Client
	validator sdk.ValAddress
	gasPrices sdk.DecCoins
	keys      chan keys.Info
	executor  executor
	fileCache filecache.Cache
}
