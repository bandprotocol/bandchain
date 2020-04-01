package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	priv secp256k1.PrivKeySecp256k1
	cdc  *codec.Codec
)

func main() {
	privS := "27313aa3fd8286b54d5dbe16a4fbbc55c7908e844e37a737997fc2ba74403812"
	privB, _ := hex.DecodeString(privS)

	copy(priv[:], privB)

	cdc = app.MakeCodec()
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cc := ctx.Config
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
	appGenTxs, accounts, err := CollectGenTxsAndAccounts(cdc, txsDir, genDoc)

	if err != nil {
		return appState, err
	}

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
func CollectGenTxsAndAccounts(cdc *codec.Codec, genTxsDir string,
	genDoc tmtypes.GenesisDoc,
) (appGenTxs []authtypes.StdTx, accounts genaccounts.GenesisAccounts, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, accounts, err
	}

	for idx, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, accounts, err
		}
		var genStdTx authtypes.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, accounts, err
		}

		// the memo flag is used to store
		// the ip and node-id, for example this may be:
		// "528fd3df22b31f4969b05652bfe8f0fe921321d5@192.168.2.37:26656"
		nodeAddrIP := genStdTx.GetMemo()
		if len(nodeAddrIP) == 0 {
			return appGenTxs, accounts, fmt.Errorf(
				"couldn't find node's address and IP in %s", fo.Name())
		}

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {
			return appGenTxs, accounts, errors.New(
				"each genesis transaction must provide a single genesis message")
		}

		createMsg := msgs[0].(stakingtypes.MsgCreateValidator)

		// delegate from band
		delegateMsg := stakingtypes.NewMsgDelegate(
			sdk.AccAddress(priv.PubKey().Address()),
			createMsg.ValidatorAddress,
			sdk.NewCoin("uband", sdk.NewInt(1000000000)),
		)

		delegateStdTx, err := signStdTx(uint64(idx), []sdk.Msg{delegateMsg})
		if err != nil {
			return appGenTxs, accounts, err
		}

		genAcc := genaccounts.NewGenesisAccountRaw(
			createMsg.DelegatorAddress,
			sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(1))),
			sdk.Coins{}, 0, 0, "", "",
		)
		if err := genAcc.Validate(); err != nil {
			return appGenTxs, accounts, err
		}

		appGenTxs = append(appGenTxs, genStdTx)
		appGenTxs = append(appGenTxs, delegateStdTx)
		accounts = append(accounts, genAcc)
	}
	return appGenTxs, accounts, nil
}

func signStdTx(
	seq uint64,
	msgs []sdk.Msg,
) (authtypes.StdTx, error) {
	txBldr := authtypes.NewTxBuilder(
		utils.GetTxEncoder(cdc), 0, seq,
		200000, 1, false, "bandchain", "", sdk.Coins{}, sdk.DecCoins{},
	)
	// build and sign the transaction
	signMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return authtypes.StdTx{}, err
	}

	sigBytes, err := priv.Sign(signMsg.Bytes())
	if err != nil {
		return authtypes.StdTx{}, err
	}
	sig := authtypes.StdSignature{
		PubKey:    priv.PubKey(),
		Signature: sigBytes,
	}
	return authtypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authtypes.StdSignature{sig}, signMsg.Memo), nil
}
