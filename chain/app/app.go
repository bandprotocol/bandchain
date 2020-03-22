package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	bandsupply "github.com/bandprotocol/bandchain/chain/x/supply"
	"github.com/bandprotocol/bandchain/chain/x/zoracle"
)

const (
	appName          = "BandApp"
	Bech32MainPrefix = "band"
	Bip44CoinType    = 494
)

var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.bandcli")

	// default home directories for the application daemon
	DefaultNodeHome = os.ExpandEnv("$HOME/.bandd")

	// NewBasicManager is in charge of setting up basic module elements
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},
		// D3N-specific modules
		zoracle.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:           nil,
		distr.ModuleName:                nil,
		mint.ModuleName:                 {supply.Minter},
		staking.BondedPoolName:          {supply.Burner, supply.Staking},
		staking.NotBondedPoolName:       {supply.Burner, supply.Staking},
		gov.ModuleName:                  {supply.Burner},
		transfer.GetModuleAccountName(): {supply.Minter, supply.Burner},
	}
)

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

type bandApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// Keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// Keepers
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
	IBCKeeper      ibc.Keeper
	TransferKeeper transfer.Keeper
	ZoracleKeeper  zoracle.Keeper

	// Decoder for unmarshaling []byte into sdk.Tx
	TxDecoder sdk.TxDecoder
	// Deliver Context that is set during BeingBlock and unset during EndBlock; primarily for gas refund
	DeliverContext sdk.Context

	// Module Manager
	mm *module.Manager
}

// NewBandApp is a constructor function for bandApp
func NewBandApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *bandApp {

	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	cdc := codecstd.MakeCodec(ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, bank.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, ibc.StoreKey, transfer.StoreKey,
		evidence.StoreKey, zoracle.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &bandApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		TxDecoder:      auth.DefaultTxDecoder(cdc),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.ParamsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	// Set specific supspaces
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.ParamsKeeper.Subspace(crisis.DefaultParamspace)
	zoracleSubspace := app.ParamsKeeper.Subspace(zoracle.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.AccountKeeper = auth.NewAccountKeeper(
		appCodec,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.BankKeeper = bank.NewBaseKeeper(
		appCodec,
		keys[bank.StoreKey],
		app.AccountKeeper,
		bankSubspace,
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.SupplyKeeper = supply.NewKeeper(
		appCodec,
		keys[supply.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		maccPerms,
	)

	// Wrapped supply keeper allows burned tokens to be transfereed to community pool
	wrappedSupplyKeeper := bandsupply.WrapSupplyKeeperBurnToCommunityPool(app.SupplyKeeper)

	// The staking keeper
	StakingKeeper := staking.NewKeeper(
		appCodec,
		keys[staking.StoreKey],
		app.BankKeeper,
		&wrappedSupplyKeeper,
		stakingSubspace,
	)

	app.MintKeeper = mint.NewKeeper(
		appCodec,
		keys[mint.StoreKey],
		mintSubspace,
		&StakingKeeper,
		app.SupplyKeeper,
		auth.FeeCollectorName,
	)

	app.DistrKeeper = distr.NewKeeper(
		appCodec,
		keys[distr.StoreKey],
		distrSubspace,
		app.BankKeeper,
		&StakingKeeper,
		app.SupplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	// DistrKeeper must be set afterward due to the circular reference of supply-staking-distr
	wrappedSupplyKeeper.SetDistrKeeper(&app.DistrKeeper)

	app.SlashingKeeper = slashing.NewKeeper(
		appCodec,
		keys[slashing.StoreKey],
		&StakingKeeper,
		slashingSubspace,
	)

	app.CrisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper))
	app.GovKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], govSubspace, app.SupplyKeeper, &StakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: StakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *StakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks()),
	)

	app.IBCKeeper = ibc.NewKeeper(app.cdc, keys[ibc.StoreKey], app.StakingKeeper)
	transferCapKey := app.IBCKeeper.PortKeeper.BindPort(bank.ModuleName)
	app.TransferKeeper = transfer.NewKeeper(
		app.cdc, keys[transfer.StoreKey], transferCapKey,
		app.IBCKeeper.ChannelKeeper, app.BankKeeper, app.SupplyKeeper,
	)

	app.ZoracleKeeper = zoracle.NewKeeper(
		cdc,
		keys[zoracle.StoreKey],
		app.BankKeeper,
		app.StakingKeeper,
		zoracleSubspace,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper, app.SupplyKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		crisis.NewAppModule(&app.CrisisKeeper),
		zoracle.NewAppModule(app.ZoracleKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.BankKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.SupplyKeeper, app.StakingKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper, app.SupplyKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.SupplyKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(zoracle.ModuleName, crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, supply.ModuleName,
		crisis.ModuleName, genutil.ModuleName, evidence.ModuleName, zoracle.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(
		app.AccountKeeper,
		app.SupplyKeeper,
		app.IBCKeeper,
		auth.DefaultSigVerificationGasConsumer,
	))

	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

func (app *bandApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, app.cdc, genesisState)
}

func (app *bandApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.DeliverContext = ctx
	return app.mm.BeginBlock(ctx, req)
}

func (app *bandApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *bandApp) Commit() (res abci.ResponseCommit) {
	app.DeliverContext = sdk.Context{}
	return app.BaseApp.Commit()
}

func (app *bandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
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
				if _, ok := msg.(zoracle.MsgReportData); !ok {
					isAllReportTxs = false
					break
				}
			}
			if isAllReportTxs && !stdTx.Fee.Amount.IsZero() {
				err := app.SupplyKeeper.SendCoinsFromModuleToAccount(
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

func (app *bandApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *bandApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
