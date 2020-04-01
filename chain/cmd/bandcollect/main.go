package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/bandlib"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const flagGenTxDir = "gentx-dir"

var provider bandlib.BandProvider

func main() {
	privS := "27313aa3fd8286b54d5dbe16a4fbbc55c7908e844e37a737997fc2ba74403812"
	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)
	provider, _ = bandlib.NewBandProvider("http://d3n-debug.bandprotocol.com:26657", priv)

	cdc := app.MakeCodec()
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cc := ctx.Config
	// config.SetRoot(viper.GetString(cli.HomeFlag))
	cc.SetRoot(os.ExpandEnv("$HOME/.bandd"))

	genDoc, err := tmtypes.GenesisDocFromFile(cc.GenesisFile())
	if err != nil {
		panic(err)
	}

	genTxsDir := filepath.Join(cc.RootDir, "config", "gentx")

	_, err = GenAppStateFromConfig(cdc, cc, genTxsDir, *genDoc)
	if err != nil {
		panic(err)
	}
}

// GenAppStateFromConfig gets the genesis app state from the config
func GenAppStateFromConfig(cdc *codec.Codec, config *cfg.Config,
	txsDir string, genDoc tmtypes.GenesisDoc,
) (appState json.RawMessage, err error) {

	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err := CollectStdTxs(
		cdc, config.Moniker, txsDir, genDoc)
	if err != nil {
		return appState, err
	}

	config.P2P.PersistentPeers = persistentPeers
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return appState, errors.New("there must be at least one genesis tx")
	}

	// create the app state
	appGenesisState, err := genutil.GenesisStateFromGenDoc(cdc, genDoc)
	if err != nil {
		return appState, err
	}

	appGenesisState, err = genutil.SetGenTxsInAppGenesisState(cdc, appGenesisState, appGenTxs)
	if err != nil {
		return appState, err
	}
	appState, err = codec.MarshalJSONIndent(cdc, appGenesisState)
	if err != nil {
		return appState, err
	}

	genDoc.AppState = appState
	err = genutil.ExportGenesisFile(&genDoc, config.GenesisFile())
	return appState, err
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns
// the list of appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(cdc *codec.Codec, moniker, genTxsDir string,
	genDoc tmtypes.GenesisDoc,
) (appGenTxs []authtypes.StdTx, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, persistentPeers, err
	}

	// addresses and IPs (and port) validator server info
	var addressesIPs []string

	for idx, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, persistentPeers, err
		}
		var genStdTx authtypes.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, persistentPeers, err
		}
		appGenTxs = append(appGenTxs, genStdTx)

		// the memo flag is used to store
		// the ip and node-id, for example this may be:
		// "528fd3df22b31f4969b05652bfe8f0fe921321d5@192.168.2.37:26656"
		nodeAddrIP := genStdTx.GetMemo()
		if len(nodeAddrIP) == 0 {
			return appGenTxs, persistentPeers, fmt.Errorf(
				"couldn't find node's address and IP in %s", fo.Name())
		}

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {
			return appGenTxs, persistentPeers, errors.New(
				"each genesis transaction must provide a single genesis message")
		}

		createMsg := msgs[0].(stakingtypes.MsgCreateValidator)

		// delegate from band
		delegateMsg := stakingtypes.NewMsgDelegate(
			provider.Sender(), createMsg.ValidatorAddress, sdk.NewCoin("uband", sdk.NewInt(1000000000)),
		)

		stdTx, err := provider.SignStdTx(uint64(idx), []sdk.Msg{delegateMsg})
		if err != nil {
			return appGenTxs, persistentPeers, err
		}
		appGenTxs = append(appGenTxs, stdTx)

		// exclude itself from persistent peers
		if createMsg.Description.Moniker != moniker {
			addressesIPs = append(addressesIPs, nodeAddrIP)
		}
	}

	sort.Strings(addressesIPs)
	persistentPeers = strings.Join(addressesIPs, ",")

	return appGenTxs, persistentPeers, nil
}
