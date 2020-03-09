package keeper

import (
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
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

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper) {
	keyRequest := sdk.NewKVStoreKey(types.StoreKey)
	tkeyRequest := sdk.NewKVStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	config := sdk.GetConfig()
	SetBech32AddressPrefixesAndBip44CoinType(config)

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

	keeper := NewKeeper(cdc, keyRequest, bk, sk, supplyKeeper, pk.Subspace(types.DefaultParamspace))
	require.Equal(t, account.GetAddress(), addr)
	accountKeeper.SetAccount(ctx, account)

	require.Equal(t, account, accountKeeper.GetAccount(ctx, addr))

	// Set default parameter
	keeper.SetMaxDataSourceExecutableSize(ctx, types.DefaultMaxDataSourceExecutableSize)
	keeper.SetMaxOracleScriptCodeSize(ctx, types.DefaultMaxOracleScriptCodeSize)
	keeper.SetMaxCalldataSize(ctx, types.DefaultMaxCalldataSize)
	keeper.SetMaxDataSourceCountPerRequest(ctx, types.DefaultMaxDataSourceCountPerRequest)
	keeper.SetMaxRawDataReportSize(ctx, types.DefaultMaxRawDataReportSize)
	keeper.SetMaxResultSize(ctx, types.DefaultMaxResultSize)
	keeper.SetEndBlockExecuteGasLimit(ctx, types.DefaultEndBlockExecuteGasLimit)
	keeper.SetMaxNameLength(ctx, types.DefaultMaxNameLength)
	keeper.SetMaxDescriptionLength(ctx, types.DefaultDescriptionLength)

	return ctx, keeper
}

func SetupTestValidator(ctx sdk.Context, keeper Keeper, pk string, power int64) sdk.ValAddress {
	pubKey := NewPubKey(pk)
	validatorAddress := sdk.ValAddress(pubKey.Address())
	initTokens := sdk.TokensFromConsensusPower(power)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
	keeper.CoinKeeper.AddCoins(ctx, sdk.AccAddress(pubKey.Address()), initCoins)

	msgCreateValidator := staking.NewTestMsgCreateValidator(
		validatorAddress, pubKey, sdk.TokensFromConsensusPower(power),
	)
	stakingHandler := staking.NewHandler(keeper.StakingKeeper)
	stakingHandler(ctx, msgCreateValidator)

	keeper.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	return validatorAddress
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

func GetAddressFromPub(pub string) (addr sdk.AccAddress, err error) {
	return sdk.AccAddressFromHex(NewPubKey(pub).Address().String())
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
		20000,
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
