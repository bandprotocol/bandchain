package keeper_test

import (
	"math/rand"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
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
	Owner      account
	Alice      account
	Bob        account
	Carol      account
	Validator1 account
	Validator2 account
	Validator3 account
)

var (
	BasicName           = "BASIC_NAME"
	BasicDesc           = "BASIC_DESCRIPTION"
	BasicSchema         = "BASIC_SCHEMA"
	BasicSourceCodeURL  = "BASIC_SOURCE_CODE_URL"
	BasicFilename       = "BASIC_FILENAME"
	BasicCalldata       = []byte("BASIC_CALLDATA")
	CoinsZero           = sdk.NewCoins()
	Coins10uband        = sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	Coins20uband        = sdk.NewCoins(sdk.NewInt64Coin("uband", 20))
	Coins1000000uband   = sdk.NewCoins(sdk.NewInt64Coin("uband", 1000000))
	Coins100000000uband = sdk.NewCoins(sdk.NewInt64Coin("uband", 100000000))
	Coins99999999uband  = sdk.NewCoins(sdk.NewInt64Coin("uband", 99999999))
)

func init() {
	bandapp.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	Owner = createArbitraryAccount(r)
	Alice = createArbitraryAccount(r)
	Bob = createArbitraryAccount(r)
	Carol = createArbitraryAccount(r)
	Validator1 = createArbitraryAccount(r)
	Validator2 = createArbitraryAccount(r)
	Validator3 = createArbitraryAccount(r)
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

func createValidatorTx(acc account, moniker string, selfDelegation sdk.Coin) authtypes.StdTx {
	msg := staking.NewMsgCreateValidator(
		acc.ValAddress, acc.PubKey, selfDelegation,
		staking.NewDescription(moniker, "", "", "", ""),
		staking.NewCommissionRates(sdk.MustNewDecFromStr("0.125"), sdk.MustNewDecFromStr("0.3"), sdk.MustNewDecFromStr("0.01")),
		sdk.NewInt(1),
	)
	txMsg := authtypes.StdSignMsg{
		ChainID:       "bandchain",
		AccountNumber: 0,
		Sequence:      0,
		Fee:           auth.NewStdFee(200000, sdk.Coins{}),
		Msgs:          []sdk.Msg{msg},
		Memo:          "",
	}
	sigBytes, err := acc.PrivKey.Sign(txMsg.Bytes())
	if err != nil {
		panic(err)
	}

	sigs := []authtypes.StdSignature{{
		PubKey:    acc.PubKey.Bytes(),
		Signature: sigBytes,
	}}
	return authtypes.NewStdTx([]sdk.Msg{msg}, auth.NewStdFee(200000, sdk.Coins{}), sigs, "")
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
		&auth.BaseAccount{Address: Validator1.Address},
		&auth.BaseAccount{Address: Validator2.Address},
		&auth.BaseAccount{Address: Validator3.Address},
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
		{
			Address: Validator1.Address,
			Coins:   Coins100000000uband,
		},
		{
			Address: Validator2.Address,
			Coins:   Coins100000000uband,
		},
		{
			Address: Validator3.Address,
			Coins:   Coins100000000uband,
		},
	}, sdk.NewCoins(sdk.NewInt64Coin("uband", 304000000)))
	genesis[bank.ModuleName] = app.Codec().MustMarshalJSON(bankGenesis)
	genutilGenesis := genutil.NewGenesisStateFromStdTx([]authtypes.StdTx{
		createValidatorTx(Validator1, "validator1", Coins100000000uband[0]),
		createValidatorTx(Validator2, "validator2", Coins1000000uband[0]),
		createValidatorTx(Validator3, "validator3", Coins99999999uband[0]),
	})
	genesis[genutil.ModuleName] = app.Codec().MustMarshalJSON(genutilGenesis)
	// Initialize the sim blockchain. We are ready for testing!
	app.InitChain(abci.RequestInitChain{
		ChainId:       "bandchain",
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), genesis),
	})
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	// Funds all the seed accounts with 1000000uband initially.
	return app, ctx, app.OracleKeeper
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}
