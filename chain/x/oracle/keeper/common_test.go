package keeper_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/simapp"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/bandprotocol/bandchain/go-owasm/api"
)

const (
	ChainID = "bandchain"
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

func getTestDataSource(executable string) (ds types.DataSource, clear func()) {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile([]byte(executable))
	return types.NewDataSource(Owner.Address, "Test data source", "For test only", filename),
		func() { deleteFile(filepath.Join(dir, filename)) }
}

func getTestOracleScript() (os types.OracleScript, clear func()) {
	absPath, _ := filepath.Abs("../../../pkg/owasm/res/beeb.wat")
	rawWAT, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	code := wat2wasm(rawWAT)
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	compiledCode, err := api.Compile(code, types.MaxDataSize)
	if err != nil {
		panic(err)
	}
	filename := f.AddFile(compiledCode)
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

// wat2wasm compiles the given Wat content to Wasm, relying on the host's wat2wasm program.
func wat2wasm(wat []byte) []byte {
	inputFile, err := ioutil.TempFile("", "input")
	if err != nil {
		panic(err)
	}
	defer os.Remove(inputFile.Name())
	outputFile, err := ioutil.TempFile("", "output")
	if err != nil {
		panic(err)
	}
	defer os.Remove(outputFile.Name())
	if _, err := inputFile.Write(wat); err != nil {
		panic(err)
	}
	if err := exec.Command("wat2wasm", inputFile.Name(), "-o", outputFile.Name()).Run(); err != nil {
		panic(err)
	}
	output, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		panic(err)
	}
	return output
}
