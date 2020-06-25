package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	bandsupply "github.com/bandprotocol/bandchain/chain/x/supply"
)

const (
	AppName          = "BandApp"
	Bech32MainPrefix = "band"
	Bip44CoinType    = 494
)

var (
	// DefaultCLIHome default home directories for bandcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.bandcli")

	// DefaultNodeHome default home directories for bandd
	DefaultNodeHome = os.ExpandEnv("$HOME/.bandd")

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		oracle.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

// BandApp extended ABCI application
type BandApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper  auth.AccountKeeper
	BankKeeper     bank.Keeper
	SupplyKeeper   supply.Keeper
	StakingKeeper  staking.Keeper
	SlashingKeeper slashing.Keeper
	MintKeeper     mint.Keeper
	DistrKeeper    distr.Keeper
	GovKeeper      gov.Keeper
	CrisisKeeper   crisis.Keeper
	ParamsKeeper   params.Keeper
	UpgradeKeeper  upgrade.Keeper
	EvidenceKeeper evidence.Keeper
	OracleKeeper   oracle.Keeper

	// Decoder for unmarshaling []byte into sdk.Tx
	TxDecoder sdk.TxDecoder
	// Deliver Context that is set during BeginBlock and unset during EndBlock; primarily for gas refund
	DeliverContext sdk.Context

	// the module manager
	mm *module.Manager
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc.Seal()
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

// NewBandApp returns a reference to an initialized BandApp.
func NewBandApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	baseAppOptions ...func(*bam.BaseApp),
) *BandApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(AppName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, supply.StoreKey, staking.StoreKey, mint.StoreKey,
		distr.StoreKey, slashing.StoreKey, gov.StoreKey, params.StoreKey, upgrade.StoreKey,
		evidence.StoreKey, oracle.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	app := &BandApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
		TxDecoder:      auth.DefaultTxDecoder(cdc),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.ParamsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[oracle.ModuleName] = app.ParamsKeeper.Subspace(oracle.DefaultParamspace)

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount)
	app.BankKeeper = bank.NewBaseKeeper(app.AccountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs())
	app.SupplyKeeper = supply.NewKeeper(cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms)
	// Wrapped supply keeper allows burned tokens to be transferred to community pool
	wrappedSupplyKeeper := bandsupply.WrapSupplyKeeperBurnToCommunityPool(app.SupplyKeeper)
	stakingKeeper := staking.NewKeeper(cdc, keys[staking.StoreKey], &wrappedSupplyKeeper, app.subspaces[staking.ModuleName])
	app.MintKeeper = mint.NewKeeper(cdc, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper, app.SupplyKeeper, auth.FeeCollectorName)
	app.CrisisKeeper = crisis.NewKeeper(app.subspaces[crisis.ModuleName], invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName)
	app.DistrKeeper = distr.NewKeeper(cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper, app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	// DistrKeeper must be set afterward due to the circular reference of supply-staking-distr
	wrappedSupplyKeeper.SetDistrKeeper(&app.DistrKeeper)
	app.SlashingKeeper = slashing.NewKeeper(cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName])

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.GovKeeper = gov.NewKeeper(cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.SupplyKeeper, &stakingKeeper, govRouter)

	// create evidence keeper with evidence router
	evidenceKeeper := evidence.NewKeeper(cdc, keys[evidence.StoreKey], app.subspaces[evidence.ModuleName], &stakingKeeper, app.SlashingKeeper)
	evidenceRouter := evidence.NewRouter()
	evidenceKeeper.SetRouter(evidenceRouter)
	app.EvidenceKeeper = *evidenceKeeper

	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], cdc)
	app.OracleKeeper = oracle.NewKeeper(cdc, keys[oracle.StoreKey], filepath.Join(viper.GetString(cli.HomeFlag), "files"), app.subspaces[oracle.ModuleName], stakingKeeper)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.SupplyKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		oracle.NewAppModule(app.OracleKeeper),
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.

	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName, staking.ModuleName)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, oracle.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName, supply.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, crisis.ModuleName,
		genutil.ModuleName, evidence.ModuleName, oracle.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, ante.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}
	return app
}

// Name returns the name of the App
func (app *BandApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *BandApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.DeliverContext = ctx
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *BandApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *BandApp) Commit() (res abci.ResponseCommit) {
	app.DeliverContext = sdk.Context{}
	return app.BaseApp.Commit()
}

// InitChainer application update at chain initialization
func (app *BandApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	app.DeliverContext = ctx // NOTE: This will be reset at the beginning of the first block.
	return app.mm.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height
func (app *BandApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *BandApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *BandApp) Codec() *codec.Codec {
	return app.cdc
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
