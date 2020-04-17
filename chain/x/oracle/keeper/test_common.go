package keeper

import (
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

const Bech32MainPrefix = "band"
const Bip44CoinType = 494

func createTestCodec() *codec.Codec {
	var cdc = codec.New()
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	return cdc
}

func SetBech32AddressPrefixesAndBip44CoinType(config *sdk.Config) {
	config.SetBech32PrefixForAccount(
		Bech32MainPrefix,
		Bech32MainPrefix+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForValidator(
		Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator,
		Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForConsensusNode(
		Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)
	config.SetCoinType(Bip44CoinType)
}

// TODO: Create a test context that encapsulates this.
var accountKeeper auth.AccountKeeper

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	keyRequest := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyBank := sdk.NewKVStoreKey(bank.StoreKey)
	keyIBC := sdk.NewKVStoreKey(ibc.StoreKey)
	keyCap := sdk.NewKVStoreKey(capability.StoreKey)

	config := sdk.GetConfig()
	SetBech32AddressPrefixesAndBip44CoinType(config)

	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyRequest, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIBC, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCap, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, isCheckTx, log.NewNopLogger())
	cdc := createTestCodec()
	appCodec := codecstd.NewAppCodec(cdc)

	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	pk := params.NewKeeper(appCodec, keyParams, tkeyParams)

	accountKeeper = auth.NewAccountKeeper(
		appCodec, // amino codec
		keyAcc,   // account store key
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
		appCodec,
		keyBank,
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(appCodec, keySupply, accountKeeper, bk, maccPerms)

	initTokens := sdk.TokensFromConsensusPower(10)                                       // 10^7 for staking
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(2))) // 2 = total validator address

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	sk := staking.NewKeeper(appCodec, keyRequest, bk, supplyKeeper, pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultParams())

	// set module accounts
	err = bk.SetBalances(ctx, notBondedPool.GetAddress(), totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	capabilityKeeper := capability.NewKeeper(appCodec, keyCap)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedOracleKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)

	ibcKeeper := ibc.NewKeeper(cdc, keyIBC, sk, scopedIBCKeeper)

	keeper := NewKeeper(cdc, keyRequest, bk, sk, ibcKeeper.ChannelKeeper, pk.Subspace(types.DefaultParamspace), scopedOracleKeeper)
	require.Equal(t, account.GetAddress(), addr)
	accountKeeper.SetAccount(ctx, account)

	require.Equal(t, account, accountKeeper.GetAccount(ctx, addr))

	// Set default parameter
	keeper.SetParam(ctx, types.KeyMaxExecutableSize, types.DefaultMaxDataSourceExecutableSize)
	keeper.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, types.DefaultMaxOracleScriptCodeSize)
	keeper.SetParam(ctx, types.KeyMaxCalldataSize, types.DefaultMaxCalldataSize)
	keeper.SetParam(ctx, types.KeyMaxDataSourceCountPerRequest, types.DefaultMaxDataSourceCountPerRequest)
	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, types.DefaultMaxRawDataReportSize)
	keeper.SetParam(ctx, types.KeyMaxResultSize, types.DefaultMaxResultSize)
	keeper.SetParam(ctx, types.KeyEndBlockExecuteGasLimit, types.DefaultEndBlockExecuteGasLimit)
	keeper.SetParam(ctx, types.KeyMaxNameLength, types.DefaultMaxNameLength)
	keeper.SetParam(ctx, types.KeyMaxDescriptionLength, types.DefaultDescriptionLength)
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, types.DefaultGasPerRawDataRequestPerValidator)
	keeper.SetParam(ctx, types.KeyExpirationBlockCount, types.DefaultExpirationBlockCount)
	keeper.SetParam(ctx, types.KeyExecuteGas, types.DefaultExecuteGas)
	keeper.SetParam(ctx, types.KeyPrepareGas, types.DefaultPrepareGas)

	return ctx, keeper
}

func SetupTestValidator(ctx sdk.Context, keeper Keeper, pk string, power int64) sdk.ValAddress {
	pubKey := NewPubKey(pk)
	validatorAddress := sdk.ValAddress(pubKey.Address())
	initTokens := sdk.TokensFromConsensusPower(power)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))

	addr := sdk.AccAddress(pubKey.Address())
	keeper.CoinKeeper.SetBalances(ctx, addr, initCoins)
	accountKeeper.SetAccount(ctx, accountKeeper.NewAccountWithAddress(ctx, addr))

	msgCreateValidator := staking.NewMsgCreateValidator(
		validatorAddress, pubKey,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(power)),
		staking.Description{},
		staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()), sdk.OneInt(),
	)

	stakingHandler := staking.NewHandler(keeper.StakingKeeper)
	_, err := stakingHandler(ctx, msgCreateValidator)
	if err != nil {
		panic(err)
	}

	keeper.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	return validatorAddress
}

func NewPubKey(pk string) crypto.PubKey {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

func GetAddressFromPub(pub string) sdk.AccAddress {
	return sdk.AccAddress(NewPubKey(pub).Address())
}

func NewUBandCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(amount)))
}

func newDefaultRequest() types.Request {
	return types.NewRequest(
		1,
		[]byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("validator1")), sdk.ValAddress([]byte("validator2"))},
		2,
		0,
		1581503227,
		100,
		"clientID",
	)
}

func GetTestOracleScript(path string) types.OracleScript {
	absPath, _ := filepath.Abs(path)
	code, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	return types.NewOracleScript(
		sdk.AccAddress([]byte("owner")),
		"silly script",
		"description",
		code,
	)
}

func GetTestDataSource() types.DataSource {
	return types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		"description",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
}
