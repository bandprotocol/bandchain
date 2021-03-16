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
	bandante "github.com/bandprotocol/bandchain/chain/x/oracle/ante"
	bandsupply "github.com/bandprotocol/bandchain/chain/x/supply"
)

const (
	AppName          = "BandApp"
	Bech32MainPrefix = "band"
	Bip44CoinType    = 494
)

var (
	// DefaultCLIHome is the default home directories for bandcli.
	DefaultCLIHome = os.ExpandEnv("$HOME/.bandcli")
	// DefaultNodeHome is the default home directories for bandd.
	DefaultNodeHome = os.ExpandEnv("$HOME/.bandd")
	// ModuleBasics is in charge of setting up basic, non-dependant module elements.
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

// BandApp is the application of BandChain, extended base ABCI application.
type BandApp struct {
	*bam.BaseApp
	cdc            *codec.Codec
	invCheckPeriod uint
	// Keys to access the substores.
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey
	// Module keepers, publicly accessible to facilate testing and extending (see emitter).
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
	// Deliver context, set during InitGenesis/BeginBlock and cleared during Commit. It allows
	// anyone with access to BandApp to read/mutate consensus state anytime. USE WITH CARE!
	DeliverContext sdk.Context
	// Module manager.
	mm *module.Manager
	// List of hooks
	hooks []Hook
}

// MakeCodec returns BandChain codec.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	return cdc.Seal()
}

// SetBech32AddressPrefixesAndBip44CoinType sets the global Bech32 prefixes and HD wallet coin type.
func SetBech32AddressPrefixesAndBip44CoinType(config *sdk.Config) {
	accountPrefix := Bech32MainPrefix
	validatorPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)
	config.SetCoinType(Bip44CoinType)
}

// NewBandApp returns a reference to an initialized BandApp.
func NewBandApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	disableFeelessReports bool, requesters []string, baseAppOptions ...func(*bam.BaseApp),
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
	}
	// Initialize params keeper and module subspaces.
	app.ParamsKeeper = params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	evidenceSubspace := app.ParamsKeeper.Subspace(evidence.DefaultParamspace)
	govSubspace := app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	oracleSubspace := app.ParamsKeeper.Subspace(oracle.DefaultParamspace)
	// Add module keepers.
	app.AccountKeeper = auth.NewAccountKeeper(cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.BankKeeper = bank.NewBaseKeeper(app.AccountKeeper, bankSubspace, app.ModuleAccountAddrs())
	app.SupplyKeeper = supply.NewKeeper(cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms)
	// wrappedSupplyKeeper overrides burn token behavior to instead transfer to community pool.
	wrappedSupplyKeeper := bandsupply.WrapSupplyKeeperBurnToCommunityPool(app.SupplyKeeper)
	stakingKeeper := staking.NewKeeper(cdc, keys[staking.StoreKey], &wrappedSupplyKeeper, stakingSubspace)
	app.MintKeeper = mint.NewKeeper(cdc, keys[mint.StoreKey], mintSubspace, &stakingKeeper, app.SupplyKeeper, auth.FeeCollectorName)
	app.DistrKeeper = distr.NewKeeper(cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper, app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	// DistrKeeper must be set afterward due to the circular reference between supply-staking-distr.
	wrappedSupplyKeeper.SetDistrKeeper(&app.DistrKeeper)
	app.CrisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName)
	app.SlashingKeeper = slashing.NewKeeper(cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace)
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], cdc)
	app.OracleKeeper = oracle.NewKeeper(cdc, keys[oracle.StoreKey], filepath.Join(viper.GetString(cli.HomeFlag), "files"), auth.FeeCollectorName, oracleSubspace, app.SupplyKeeper, &stakingKeeper, app.DistrKeeper)
	// Register the proposal types.
	govRouter := gov.NewRouter()
	govRouter.
		AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.GovKeeper = gov.NewKeeper(cdc, keys[gov.StoreKey], govSubspace, app.SupplyKeeper, &stakingKeeper, govRouter)
	// Create evidence keeper with evidence router.
	evidenceKeeper := evidence.NewKeeper(cdc, keys[evidence.StoreKey], evidenceSubspace, &stakingKeeper, app.SlashingKeeper)
	evidenceRouter := evidence.NewRouter()
	evidenceKeeper.SetRouter(evidenceRouter)
	app.EvidenceKeeper = *evidenceKeeper
	// Register the staking hooks. NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks.
	app.StakingKeeper = *stakingKeeper.SetHooks(staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()))
	// Create the module manager. NOTE: Any module instantiated in the module manager that is later modified must be passed by reference here.
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
	// NOTE: Oracle module must occur before distr as it takes some fee to distribute to active oracle validators.
	// NOTE: During begin block slashing happens after distr.BeginBlocker so that there is nothing left
	// over in the validator fee pool, so as to keep the CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, mint.ModuleName, oracle.ModuleName, distr.ModuleName, slashing.ModuleName,
		evidence.ModuleName, staking.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisis.ModuleName, gov.ModuleName, staking.ModuleName, oracle.ModuleName,
	)
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName, supply.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, oracle.ModuleName, crisis.ModuleName,
		genutil.ModuleName, evidence.ModuleName,
	)
	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
	// Initialize stores.
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	// initialize BaseApp.
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler := ante.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer)
	if !disableFeelessReports {
		anteHandler = bandante.NewFeelessReportsAnteHandler(anteHandler, app.OracleKeeper)
	}
	if len(requesters) > 0 {
		anteHandler = bandante.NewWhiteListAnteHandler(anteHandler, app.OracleKeeper, requesters)
	}
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)
	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}
	return app
}

// Name returns the name of the App.
func (app *BandApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block.
func (app *BandApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.DeliverContext = ctx
	res := app.mm.BeginBlock(ctx, req)
	for _, hook := range app.hooks {
		hook.AfterBeginBlock(ctx, req, res)
	}
	return res
}

// EndBlocker application updates every end block.
func (app *BandApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.mm.EndBlock(ctx, req)
	for _, hook := range app.hooks {
		hook.AfterEndBlock(ctx, req, res)
	}
	return res
}

// Commit overrides the default BaseApp's ABCI commit by adding DeliverContext clearing.
func (app *BandApp) Commit() (res abci.ResponseCommit) {
	for _, hook := range app.hooks {
		hook.BeforeCommit()
	}
	app.DeliverContext = sdk.Context{}
	return app.BaseApp.Commit()
}

// InitChainer application update at chain initialization
func (app *BandApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	app.DeliverContext = ctx // NOTE: This will be reset at the beginning of the first block.
	res := app.mm.InitGenesis(ctx, genesisState)
	for _, hook := range app.hooks {
		hook.AfterInitChain(ctx, req, res)
	}
	return res
}

// DeliverTx overwrite DeliverTx to apply afterDeliverTx hook
func (app *BandApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	res := app.BaseApp.DeliverTx(req)
	for _, hook := range app.hooks {
		hook.AfterDeliverTx(app.DeliverContext, req, res)
	}
	return res
}

func (app *BandApp) Query(req abci.RequestQuery) abci.ResponseQuery {
	hookReq := req

	// when a client did not provide a query height, manually inject the latest
	if hookReq.Height == 0 {
		hookReq.Height = app.LastBlockHeight()
	}

	for _, hook := range app.hooks {
		res, stop := hook.ApplyQuery(hookReq)
		if stop {
			return res
		}
	}
	return app.BaseApp.Query(req)
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

// AddHook appends hook that will be call after process abci request
func (app *BandApp) AddHook(hook Hook) {
	app.hooks = append(app.hooks, hook)
}
