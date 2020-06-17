package simapp

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
)

// Account is a data structure to store key of test account.
type Account struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

// nolint
var (
	Owner      Account
	Alice      Account
	Bob        Account
	Carol      Account
	Validator1 Account
	Validator2 Account
	Validator3 Account
)

// nolint
var (
	Coins1000000uband   = sdk.NewCoins(sdk.NewInt64Coin("uband", 1000000))
	Coins99999999uband  = sdk.NewCoins(sdk.NewInt64Coin("uband", 99999999))
	Coins100000000uband = sdk.NewCoins(sdk.NewInt64Coin("uband", 100000000))
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

func createArbitraryAccount(r *rand.Rand) Account {
	privkeySeed := make([]byte, 12)
	r.Read(privkeySeed)
	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	return Account{
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}

func createValidatorTx(chainID string, acc Account, moniker string, selfDelegation sdk.Coin) authtypes.StdTx {
	msg := staking.NewMsgCreateValidator(
		acc.ValAddress, acc.PubKey, selfDelegation,
		staking.NewDescription(moniker, "", "", "", ""),
		staking.NewCommissionRates(sdk.MustNewDecFromStr("0.125"), sdk.MustNewDecFromStr("0.3"), sdk.MustNewDecFromStr("0.01")),
		sdk.NewInt(1),
	)
	txMsg := authtypes.StdSignMsg{
		ChainID:       chainID,
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
		PubKey:    acc.PubKey,
		Signature: sigBytes,
	}}
	return authtypes.NewStdTx([]sdk.Msg{msg}, auth.NewStdFee(200000, sdk.Coins{}), sigs, "")
}

// NewSimApp creates instance of our app using in test.
func NewSimApp(chainID string, logger log.Logger) *bandapp.BandApp {
	db := dbm.NewMemDB()
	app := bandapp.NewBandApp(logger, db, nil, true, 0, map[int64]bool{}, "")
	genesis := bandapp.NewDefaultGenesisState()
	// Funds seed accounts and validators with 1000000uband and 100000000uband initially.
	authGenesis := auth.NewGenesisState(auth.DefaultParams(), []authexported.GenesisAccount{
		&auth.BaseAccount{Address: Owner.Address, Coins: Coins1000000uband},
		&auth.BaseAccount{Address: Alice.Address, Coins: Coins1000000uband},
		&auth.BaseAccount{Address: Bob.Address, Coins: Coins1000000uband},
		&auth.BaseAccount{Address: Carol.Address, Coins: Coins1000000uband},
		&auth.BaseAccount{Address: Validator1.Address, Coins: Coins100000000uband},
		&auth.BaseAccount{Address: Validator2.Address, Coins: Coins100000000uband},
		&auth.BaseAccount{Address: Validator3.Address, Coins: Coins100000000uband},
	})
	genesis[auth.ModuleName] = app.Codec().MustMarshalJSON(authGenesis)
	genutilGenesis := genutil.NewGenesisStateFromStdTx([]authtypes.StdTx{
		createValidatorTx(chainID, Validator1, "validator1", Coins100000000uband[0]),
		createValidatorTx(chainID, Validator2, "validator2", Coins1000000uband[0]),
		createValidatorTx(chainID, Validator3, "validator3", Coins99999999uband[0]),
	})
	genesis[genutil.ModuleName] = app.Codec().MustMarshalJSON(genutilGenesis)
	// Initialize the sim blockchain. We are ready for testing!
	app.InitChain(abci.RequestInitChain{
		ChainId:       chainID,
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), genesis),
	})
	return app
}
