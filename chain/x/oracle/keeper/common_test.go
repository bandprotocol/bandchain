package keeper_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/simapp"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainID                       = "bandchain"
	ChainIDA                      = "chainA"
	ChainIDB                      = "chainB"
	TestClientIDA                 = "clientA"
	TestClientIDB                 = "clientB"
	TestPortA                     = "testporta"
	TestPortB                     = "testportb"
	TestChannelA                  = "testchannela"
	TestChannelB                  = "testchannelb"
	TestConnectionA               = "connectionAtoB"
	TestConnectionB               = "connectionBtoA"
	TrustingPeriod  time.Duration = time.Hour * 24 * 7 * 2
	UbdPeriod       time.Duration = time.Hour * 24 * 7 * 3
	MaxClockDrift   time.Duration = time.Second * 10
)

var (
	Owner      = simapp.Owner
	Alice      = simapp.Alice
	Bob        = simapp.Bob
	Carol      = simapp.Carol
	Validator1 = simapp.Validator1
	Validator2 = simapp.Validator2
	Validator3 = simapp.Validator3
)

var (
	BasicName          = "BASIC_NAME"
	BasicDesc          = "BASIC_DESCRIPTION"
	BasicSchema        = "BASIC_SCHEMA"
	BasicSourceCodeURL = "BASIC_SOURCE_CODE_URL"
	BasicFilename      = "BASIC_FILENAME"
	BasicCalldata      = []byte("BASIC_CALLDATA")
	CoinsZero          = sdk.NewCoins()
	Coins10uband       = sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	Coins20uband       = sdk.NewCoins(sdk.NewInt64Coin("uband", 20))
)

func createTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	app := simapp.NewSimApp(ChainID, log.NewNopLogger())
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx, app.OracleKeeper
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func getTestDataSource() (ds types.DataSource, clear func()) {
	dir, err := ioutil.TempDir("/tmp", "filecache")
	if err != nil {
		panic(err)
	}
	f := filecache.New(dir)
	filename := f.AddFile([]byte("executable"))
	return types.NewDataSource(Owner.Address, "Test data source", "For test only", filename),
		func() { deleteFile(filepath.Join(dir, filename)) }
}

func getTestOracleScript() (os types.OracleScript, clear func()) {
	dir, err := ioutil.TempDir("/tmp", "filecache")
	if err != nil {
		panic(err)
	}
	f := filecache.New(dir)
	filename := f.AddFile([]byte("code"))
	return types.NewOracleScript(Owner.Address, "Test oracle script",
		"For test only", filename, "", "test URL",
	), func() { deleteFile(filepath.Join(dir, filename)) }
}

func loadOracleScriptFromWasm(path string) (os types.OracleScript, clear func()) {
	absPath, _ := filepath.Abs(path)
	code, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	dir, err := ioutil.TempDir("/tmp", "filecache")
	if err != nil {
		panic(err)
	}
	f := filecache.New(dir)
	filename := f.AddFile(code)
	return types.NewOracleScript(
		Owner.Address, "imported script", "description",
		filename, "schema", "sourceCodeURL",
	), func() { deleteFile(filepath.Join(dir, filename)) }
}
