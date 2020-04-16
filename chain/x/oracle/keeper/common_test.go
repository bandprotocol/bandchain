package keeper_test

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

type account struct {
	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address sdk.AccAddress
}

var OWNER account
var ALICE account
var BOB account
var CAROL account

var COINS_10_UBAND = sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
var COINS_20_UBAND = sdk.NewCoins(sdk.NewInt64Coin("uband", 20))

func init() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	OWNER = createArbitraryAccount(r)
	ALICE = createArbitraryAccount(r)
	BOB = createArbitraryAccount(r)
	CAROL = createArbitraryAccount(r)
}

func createArbitraryAccount(r *rand.Rand) account {
	privkeySeed := make([]byte, 12)
	r.Read(privkeySeed)
	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	return account{
		PrivKey: privKey,
		PubKey:  privKey.PubKey(),
		Address: sdk.AccAddress(privKey.PubKey().Address()),
	}
}

func createTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	db := dbm.NewMemDB()
	app := bandapp.NewBandApp(log.NewNopLogger(), db, nil, true, 0, map[int64]bool{}, "")
	app.InitChain(abci.RequestInitChain{
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), bandapp.NewDefaultGenesisState()),
	})
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx, app.OracleKeeper
}
