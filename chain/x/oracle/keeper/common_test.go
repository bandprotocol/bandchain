package keeper_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/simapp"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
	app := simapp.NewSimApp()
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx, app.OracleKeeper
}

func newDefaultRequest() types.Request {
	return types.NewRequest(
		1,
		[]byte("calldata"),
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress},
		2,
		0,
		1581503227,
		"clientID",
		nil,
		[]types.ExternalID{42},
	)
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func getTestDataSource(executable string) (ds types.DataSource, clear func()) {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile([]byte(executable))
	return types.NewDataSource(Owner.Address, "Test data source", "For test only", filename),
		func() { deleteFile(filepath.Join(dir, filename)) }
}

func getTestOracleScript() (os types.OracleScript, clear func()) {
	absPath, _ := filepath.Abs("../../../pkg/owasm/res/beeb.wasm")
	code, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile(code)
	return types.NewOracleScript(
		Owner.Address, "imported script", "description",
		filename,
		"schema",
		"sourceCodeURL",
	), func() { deleteFile(filepath.Join(dir, filename)) }
}

func getBadOracleScript() (os types.OracleScript, clear func()) {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile([]byte("bad_code"))
	return types.NewOracleScript(
		Owner.Address, "imported script", "description",
		filename,
		"beeb", "sourceCodeURL",
	), func() { deleteFile(filepath.Join(dir, filename)) }
}
