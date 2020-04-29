package keeper_test

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

type account struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

var (
	Owner account
	Alice account
	Bob   account
	Carol account
)

var (
	BasicName          = "BASIC_NAME"
	BasicDesc          = "BASIC_DESCRIPTION"
	BasicCode          = []byte("BASIC_WASM_CODE")
	BasicSchema        = "BASIC_SCHEMA"
	BasicSourceCodeURL = "BASIC_SOURCE_CODE_URL"
	BasicExec          = []byte("BASIC_EXECUTABLE")
	BasicCalldata      = []byte("BASIC_CALLDATA")
	CoinsZero          = sdk.NewCoins()
	Coins10uband       = sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	Coins20uband       = sdk.NewCoins(sdk.NewInt64Coin("uband", 20))
	Coins1000000uband  = sdk.NewCoins(sdk.NewInt64Coin("uband", 1000000))
)

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
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}

func createTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	db := dbm.NewMemDB()
	app := bandapp.NewBandApp(log.NewNopLogger(), db, nil, true, 0, map[int64]bool{}, "")
	genesis := bandapp.NewDefaultGenesisState()
	// Funds all the seed accounts with 1000000uband initially.
	authGenesis := auth.NewGenesisState(auth.DefaultParams(), []authexported.GenesisAccount{
		&auth.BaseAccount{Address: Owner.Address},
		&auth.BaseAccount{Address: Alice.Address},
		&auth.BaseAccount{Address: Bob.Address},
		&auth.BaseAccount{Address: Carol.Address},
	})
	genesis[auth.ModuleName] = app.Codec().MustMarshalJSON(authGenesis)
	bankGenesis := bank.NewGenesisState(bank.DefaultGenesisState().SendEnabled, []bank.Balance{
		{
			Address: Owner.Address,
			Coins:   Coins1000000uband,
		},
		{
			Address: Alice.Address,
			Coins:   Coins1000000uband,
		},
		{
			Address: Bob.Address,
			Coins:   Coins1000000uband,
		},
		{
			Address: Carol.Address,
			Coins:   Coins1000000uband,
		},
	}, sdk.NewCoins(sdk.NewInt64Coin("uband", 4000000)))
	genesis[bank.ModuleName] = app.Codec().MustMarshalJSON(bankGenesis)
	// Initialize the sim blockchain. We are ready for testing!
	app.InitChain(abci.RequestInitChain{
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), genesis),
	})
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	// Funds all the seed accounts with 1000000uband initially.
	return app, ctx, app.OracleKeeper
}
