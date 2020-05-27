package oracle_test

import (
	"crypto/sha256"
	"encoding/hex"
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
)

func createTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	app := simapp.NewSimApp()
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx, app.OracleKeeper
}

func deleteFile(data []byte) {
	hash := sha256.Sum256(data)
	filename := hex.EncodeToString(hash[:])
	path := filepath.Join(viper.GetString(cli.HomeFlag), "files", filename)
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func getTestDataSource() (ds types.DataSource, clear func()) {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile([]byte("executable"))
	return types.NewDataSource(Owner.Address, "Test data source", "For test only", filename),
		func() { deleteFile([]byte("executable")) }
}

func getTestOracleScript() (os types.OracleScript, clear func()) {
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile([]byte("code"))
	return types.NewOracleScript(Owner.Address, "Test oracle script",
		"For test only", filename, "", "test URL",
	), func() { deleteFile([]byte("code")) }
}

func loadDataSourceFromExecutable(path string) (os types.DataSource, clear func()) {
	executable, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	f := filecache.New(dir)
	filename := f.AddFile(executable)
	return types.NewDataSource(
		Owner.Address, "imported data source", "description",
		filename,
	), func() { deleteFile(executable) }
}

func loadOracleScriptFromWasm(path string) (os types.OracleScript, clear func()) {
	absPath, _ := filepath.Abs(path)
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
		"schema", "sourceCodeURL",
	), func() { deleteFile(code) }
}
