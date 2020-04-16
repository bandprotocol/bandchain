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

var Owner account
var Alice account
var Bob account
var Carol account

var Coins10uband = sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
var Coins20uband = sdk.NewCoins(sdk.NewInt64Coin("uband", 20))

var BasicName = "BASIC_NAME"
var BasicDesc = "BASIC_DESCRIPTION"
var BasicCode = []byte("BASIC_WASM_CODE")
var BasicExec = []byte("BASIC_EXECUTABLE")

func init() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	Owner = createArbitraryAccount(r)
	Alice = createArbitraryAccount(r)
	Bob = createArbitraryAccount(r)
	Carol = createArbitraryAccount(r)
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
