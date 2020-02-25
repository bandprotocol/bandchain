package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/bandprotocol/d3n/chain/d3nlib"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	executeGas = 150000
)

// File to send new request to bandchain
func main() {
	// Get environment variable
	privS, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		// Default private key is eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877
		privS = "eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877"
	}
	nodeURI, ok := os.LookupEnv("NODE_URI")
	if !ok {
		// Default node uri is tcp://localhost:26657
		nodeURI = "tcp://localhost:26657"
	}

	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	tx, err := d3nlib.NewBandProvider(nodeURI, priv)

	valPrivB, _ := hex.DecodeString("06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b")
	var valPriv secp256k1.PrivKeySecp256k1
	copy(valPriv[:], valPrivB)

	valTx, err := d3nlib.NewBandProvider(nodeURI, valPriv)
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	switch args[0] {
	case "store":
		{
			coingecko, err := ioutil.ReadFile("../../datasources/coingecko_price.sh")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Coingecko script", sdk.Coins{}, coingecko, tx.Sender(),
				)},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			cryptoCompare, err := ioutil.ReadFile("../../datasources/crypto_compare_price.sh")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Crypto compare script", sdk.Coins{}, cryptoCompare, tx.Sender(),
				)},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			binance, err := ioutil.ReadFile("../../datasources/binance_price.sh")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Binance script", sdk.Coins{}, binance, tx.Sender(),
				)},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			oracleBytes, err := ioutil.ReadFile("../../owasm/res/crypto_price.wasm")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateOracleScript(tx.Sender(), "Crypto price script", oracleBytes, tx.Sender())},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))
		}
	case "send_token":
		{
			// Send token
			to, _ := sdk.AccAddressFromBech32("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
			fmt.Println(tx.SendTransaction([]sdk.Msg{bank.MsgSend{
				FromAddress: tx.Sender(),
				ToAddress:   to,
				Amount:      sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(10))),
			}}, 0, 10000000, "", "", "", flags.BroadcastBlock))
		}
	case "request":
		{
			switch args[1] {
			case "BTC":
				{
					fmt.Println(tx.SendTransaction(
						[]sdk.Msg{zoracle.NewMsgRequestData(1, []byte("BTC"), 4, 4, 100000, executeGas, tx.Sender())},
						0, 10000000, "", "", "",
						flags.BroadcastBlock,
					))
				}
			case "ETH":
				{
					fmt.Println(tx.SendTransaction(
						[]sdk.Msg{zoracle.NewMsgRequestData(1, []byte("ETH"), 4, 4, 100000, executeGas, tx.Sender())},
						0, 10000000, "", "", "",
						flags.BroadcastBlock,
					))
				}
			}
		}
	case "example-test":
		{
			bytes, err := ioutil.ReadFile("../../owasm/res/silly.wasm")
			if err != nil {
				panic(err)
			}

			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateOracleScript(tx.Sender(), "Silly script", bytes, tx.Sender())},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgCreateDataSource(tx.Sender(), "Mock Data source", sdk.Coins{}, []byte("exec"), tx.Sender())},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgRequestData(1, []byte("calldata"), 1, 1, 100, executeGas, tx.Sender())},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			fmt.Println(valTx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgReportData(1, []zoracle.RawDataReport{
					zoracle.NewRawDataReport(1, []byte("data1")),
				}, sdk.ValAddress(valTx.Sender()))},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))
		}
	case "deploy_oracle_scripts":
		{
			bytes, err := ioutil.ReadFile("../../owasm/res/silly/pkg/silly_bg.wasm")
			if err != nil {
				panic(err)
			}

			round, err := strconv.ParseUint(args[1], 10, 64)
			if round <= 0 || err != nil {
				panic("round should be more than 0")
			}

			for i := uint64(0); i < round; i++ {
				fmt.Println(tx.SendTransaction(
					[]sdk.Msg{zoracle.NewMsgCreateOracleScript(tx.Sender(), fmt.Sprintf("Silly script %d", i), bytes, tx.Sender())},
					0, 10000000, "", "", "",
					flags.BroadcastBlock,
				))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
