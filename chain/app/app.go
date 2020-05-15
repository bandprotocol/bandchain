package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	ibcclient "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	"github.com/bandprotocol/bandchain/chain/pkg/owasm"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
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
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		oracle.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {auth.Minter},
		staking.BondedPoolName:    {auth.Burner, auth.Staking},
		staking.NotBondedPoolName: {auth.Burner, auth.Staking},
		gov.ModuleName:            {auth.Burner},
	}
)

// BandApp extended ABCI application
type BandApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tKeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper    auth.AccountKeeper
	BankKeeper       bank.Keeper
	StakingKeeper    staking.Keeper
	CapabilityKeeper *capability.Keeper
	SlashingKeeper   slashing.Keeper
	MintKeeper       mint.Keeper
	DistrKeeper      distr.Keeper
	GovKeeper        gov.Keeper
	CrisisKeeper     crisis.Keeper
	ParamsKeeper     params.Keeper
	UpgradeKeeper    upgrade.Keeper
	EvidenceKeeper   evidence.Keeper
	IBCKeeper        *ibc.Keeper
	OracleKeeper     oracle.Keeper

	// Decoder for unmarshaling []byte into sdk.Tx
	TxDecoder sdk.TxDecoder
	// Deliver Context that is set during BeginBlock and unset during EndBlock; primarily for gas refund
	DeliverContext sdk.Context

	// the module manager
	mm *module.Manager
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

	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	cdc := codecstd.MakeCodec(ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	bApp := bam.NewBaseApp(AppName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	keys := sdk.NewKVStoreKeys(
		auth.StoreKey, bank.StoreKey, staking.StoreKey,
		mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, ibc.StoreKey, upgrade.StoreKey,
		evidence.StoreKey, capability.StoreKey, oracle.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capability.MemStoreKey)

	app := &BandApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		memKeys:        memKeys,
		subspaces:      make(map[string]params.Subspace),
		TxDecoder:      auth.DefaultTxDecoder(cdc),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[oracle.ModuleName] = app.ParamsKeeper.Subspace(oracle.DefaultParamspace)

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capability.NewKeeper(appCodec, keys[capability.StoreKey], memKeys[capability.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedOracleKeeper := app.CapabilityKeeper.ScopeToModule(oracle.ModuleName)

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		appCodec, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bank.NewBaseKeeper(
		appCodec, keys[bank.StoreKey], app.AccountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs(),
	)

	// TODO: Revisit wrap supply of new version and test
	// Wrapped supply keeper allows burned tokens to be transferred to community pool
	// wrappedSupplyKeeper := bandsupply.WrapSupplyKeeperBurnToCommunityPool(app.SupplyKeeper)

	stakingKeeper := staking.NewKeeper(
		appCodec, keys[staking.StoreKey], app.AccountKeeper, app.BankKeeper, app.subspaces[staking.ModuleName],
	)
	app.MintKeeper = mint.NewKeeper(
		appCodec, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, auth.FeeCollectorName,
	)
	app.CrisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.BankKeeper, auth.FeeCollectorName,
	)
	app.DistrKeeper = distr.NewKeeper(
		appCodec, keys[distr.StoreKey], app.subspaces[distr.ModuleName], app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)

	// DistrKeeper must be set afterward due to the circular reference of supply-staking-distr
	// wrappedSupplyKeeper.SetDistrKeeper(&app.DistrKeeper)

	app.SlashingKeeper = slashing.NewKeeper(
		appCodec, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], appCodec, home)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.GovKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], app.subspaces[gov.ModuleName],
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.IBCKeeper = ibc.NewKeeper(app.cdc, keys[ibc.StoreKey], stakingKeeper, scopedIBCKeeper)

	// create evidence keeper with evidence router
	evidenceKeeper := evidence.NewKeeper(
		appCodec, keys[evidence.StoreKey], &stakingKeeper, app.SlashingKeeper,
	)
	evidenceRouter := evidence.NewRouter().
		AddRoute(ibcclient.RouterKey, ibcclient.HandlerClientMisbehaviour(app.IBCKeeper.ClientKeeper))

	evidenceKeeper.SetRouter(evidenceRouter)
	app.EvidenceKeeper = *evidenceKeeper

	app.OracleKeeper = oracle.NewKeeper(
		cdc,
		keys[oracle.StoreKey],
		filepath.Join(viper.GetString(cli.HomeFlag), "files"),
		owasm.Execute,
		app.subspaces[oracle.ModuleName],
		app.BankKeeper,
		app.StakingKeeper,
		app.IBCKeeper.ChannelKeeper,
		scopedOracleKeeper,
		&app.IBCKeeper.PortKeeper,
	)

	oracleModule := oracle.NewAppModule(app.OracleKeeper)

	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(oracle.ModuleName, oracleModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.DeliverTx),
		auth.NewAppModule(appCodec, app.AccountKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(*app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(appCodec, app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		oracleModule,
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.

	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, mint.ModuleName, distr.ModuleName,
		slashing.ModuleName, staking.ModuleName, ibc.ModuleName,
	)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, oracle.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		capability.ModuleName, auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, crisis.ModuleName,
		ibc.ModuleName, genutil.ModuleName, evidence.ModuleName, oracle.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.AccountKeeper, app.BankKeeper, *app.IBCKeeper,
			ante.DefaultSigVerificationGasConsumer,
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion()
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	// Initialize and seal the capability keeper so all persistent capabilities
	// are loaded in-memory and prevent any further modules from creating scoped
	// sub-keepers.
	// This must be done during creation of baseapp rather than in InitChain so
	// that in-memory capabilities get regenerated on app restart
	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.CapabilityKeeper.InitializeAndSeal(ctx)

	return app
}

// Name returns the name of the App
func (app *BandApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *BandApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.DeliverContext = ctx
	return app.mm.BeginBlock(ctx, req)
}

// DeliverTx application updates every transction
func (app *BandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	response := app.BaseApp.DeliverTx(req)

	if response.IsOK() {
		// Refund 100% of gas fee for any successful transaction that only contains MsgReportData
		tx, err := app.TxDecoder(req.Tx)
		if err != nil { // Should never happen because BaseApp.DeliverTx succeeds
			panic(err)
		}
		isAllReportTxs := true
		if stdTx, ok := tx.(auth.StdTx); ok {
			for _, msg := range tx.GetMsgs() {
				if _, ok := msg.(oracle.MsgReportData); !ok {
					isAllReportTxs = false
					break
				}
			}
			if isAllReportTxs && !stdTx.Fee.Amount.IsZero() {
				err := app.BankKeeper.SendCoinsFromModuleToAccount(
					app.DeliverContext,
					auth.FeeCollectorName,
					stdTx.GetSigners()[0],
					stdTx.Fee.Amount,
				)
				if err != nil { // Should never happen because we just return the collected fee
					panic(err)
				}

			}
		}
	}

	return response
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
	return app.mm.InitGenesis(ctx, app.cdc, genesisState)
}

// LoadHeight loads a particular height
func (app *BandApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *BandApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
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
