package oracle_test

import (
	"fmt"
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
	code := mustGetOwasmCode("beeb.wat")
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

func mustGetOwasmCode(filename string) []byte {
	absPath, _ := filepath.Abs(fmt.Sprintf("../../pkg/owasm/res/%s", filename))
	rawWAT, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	return wat2wasm(rawWAT)
}

func mustCompileOwasm(code []byte) []byte {
	compiled, err := api.Compile(code, types.MaxDataSize)
	if err != nil {
		panic(err)
	}
	return compiled
}
