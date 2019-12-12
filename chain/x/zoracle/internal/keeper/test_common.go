package keeper

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const Bech32MainPrefix = "band"

func createTestCodec() *codec.Codec {
	var cdc = codec.New()
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	return cdc
}

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32MainPrefix, Bech32MainPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
}

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	keyRequest := sdk.NewKVStoreKey(types.StoreKey)
	tkeyRequest := sdk.NewKVStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	config := sdk.GetConfig()
	SetBech32AddressPrefixes(config)

	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyRequest, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, isCheckTx, log.NewNopLogger())
	cdc := createTestCodec()

	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // account store key
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	addr, _ := sdk.AccAddressFromBech32("band1q8ysvjkslxdkhap2zqd2n5shhay606ru3cdjwr")

	account := accountKeeper.NewAccountWithAddress(
		ctx,
		addr,
	)
	// TODO: add feeCollectorAcc, notBondedPool, bondPool
	// REF: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/x/staking/keeper/test_common.go#L109
	blacklistedAddrs := map[string]bool{}
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	initTokens := sdk.TokensFromConsensusPower(10)                                       // 10^7 for staking
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(2))) // 2 = total validator address

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	sk := staking.NewKeeper(cdc, keyRequest, tkeyRequest, supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	sk.SetParams(ctx, staking.DefaultParams())

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	keeper := NewKeeper(cdc, keyRequest, bk, sk)
	require.Equal(t, account.GetAddress(), addr)
	accountKeeper.SetAccount(ctx, account)

	require.Equal(t, account, accountKeeper.GetAccount(ctx, addr))

	return ctx, keeper
}

// NewPubKey is a function to generate public key
func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}
