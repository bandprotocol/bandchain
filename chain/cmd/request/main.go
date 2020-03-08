package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/bandprotocol/d3n/chain/bandlib"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	prepareGas = 10000
	executeGas = 150000
)

// File to send new request to bandchain
func main() {
	// Get environment variable
	privS, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		// Default private key is eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877
		privS = "27313aa3fd8286b54d5dbe16a4fbbc55c7908e844e37a737997fc2ba74403812"
	}
	nodeURI, ok := os.LookupEnv("NODE_URI")
	if !ok {
		// Default node uri is tcp://localhost:26657
		nodeURI = "tcp://localhost:26657"
	}

	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	tx, err := bandlib.NewBandStatefulClient(nodeURI, priv, 10, 5, "Request script txs")
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
				zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Coingecko script",
					"The Script that queries crypto price from https://cryptocompare.com",
					sdk.Coins{}, coingecko, tx.Sender(),
				),
				1000000, "",
			))

			cryptoCompare, err := ioutil.ReadFile("../../datasources/crypto_compare_price.sh")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Crypto compare script",
					"The Script that queries crypto price from https://cryptocompare.com",
					sdk.Coins{}, cryptoCompare, tx.Sender(),
				),
				1000000, "",
			))

			binance, err := ioutil.ReadFile("../../datasources/binance_price.sh")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Binance script",
					"The Script that queries crypto price from https://www.binance.com/en",
					sdk.Coins{}, binance, tx.Sender(),
				),
				1000000, "",
			))

			oracleBytes, err := ioutil.ReadFile("../../owasm/res/crypto_price.wasm")
			if err != nil {
				panic(err)
			}
			fmt.Println(tx.SendTransaction(
				zoracle.NewMsgCreateOracleScript(
					tx.Sender(), "Crypto price script",
					"Oracle script for getting an average crypto price from many sources.",
					oracleBytes, tx.Sender(),
				),
				3000000, "",
			))
		}
	case "send_token":
		{
			// Send token
			to, _ := sdk.AccAddressFromBech32("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
			fmt.Println(tx.SendTransaction(bank.MsgSend{
				FromAddress: tx.Sender(),
				ToAddress:   to,
				Amount:      sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(10))),
			}, 1000000, ""))
		}
	case "request":
		{
			switch args[1] {
			case "BTC":
				{
					fmt.Println(tx.SendTransaction(
						zoracle.NewMsgRequestData(
							1, []byte("BTC"), 4, 4, 100000, prepareGas, executeGas, tx.Sender(),
						), 1000000, "",
					))
				}
			case "ETH":
				{
					fmt.Println(tx.SendTransaction(
						zoracle.NewMsgRequestData(
							1, []byte("ETH"), 4, 4, 100000, prepareGas, executeGas, tx.Sender(),
						), 1000000, "",
					))
				}
			}
		}
	case "requests":
		{
			txResponses := make(chan sdk.TxResponse, 2)
			errResponses := make(chan error, 2)
			go func() {
				txRes, err := tx.SendTransaction(
					zoracle.NewMsgRequestData(
						1, []byte("BTC"), 4, 4, 100000, prepareGas, executeGas, tx.Sender(),
					), 1000000, "",
				)

				if err != nil {
					errResponses <- err
				}
				txResponses <- txRes
			}()
			go func() {
				txRes, err := tx.SendTransaction(
					zoracle.NewMsgRequestData(
						1, []byte("ETH"), 4, 4, 100000, prepareGas, executeGas, tx.Sender(),
					), 1000000, "",
				)

				if err != nil {
					errResponses <- err
				}
				txResponses <- txRes
			}()
			for i := 0; i < 2; i++ {
				select {
				case txRes := <-txResponses:
					{
						fmt.Println(txRes)
					}
				case err := <-errResponses:
					{
						fmt.Println(err)
					}
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
				zoracle.NewMsgCreateOracleScript(
					tx.Sender(), "Silly script", "Test oracle script", bytes, tx.Sender()),
				3000000, "",
			))

			fmt.Println(tx.SendTransaction(
				zoracle.NewMsgCreateDataSource(
					tx.Sender(), "Mock Data source", "Mock Script",
					sdk.Coins{}, []byte("exec"), tx.Sender(),
				), 1000000, "",
			))

			fmt.Println(tx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgRequestData(1, []byte("calldata"), 1, 1, 100, prepareGas, executeGas, tx.Sender())},
				0, 1000000, "", "", "",
				flags.BroadcastBlock,
			))

			fmt.Println(valTx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgReportData(
					1,
					sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(1))),
					[]zoracle.RawDataReport{
						zoracle.NewRawDataReport(1, []byte("data1")),
					},
					sdk.ValAddress(valTx.Sender()),
				)},
				0, 1000000, "", "", "",
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
					zoracle.NewMsgCreateOracleScript(
						tx.Sender(), fmt.Sprintf("Silly script %d", i), "Test oracle script",
						bytes, tx.Sender(),
					), 1000000, "",
				))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
