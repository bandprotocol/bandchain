package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bandprotocol/d3n/chain/d3nlib"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
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
	// case "store":
	// 	{
	// 		file, err := os.Open("../../wasm/res/result.wasm")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		defer file.Close()

	// 		stats, statsErr := file.Stat()
	// 		if statsErr != nil {
	// 			panic(statsErr)
	// 		}

	// 		size := stats.Size()
	// 		bytes := make([]byte, size)

	// 		bufr := bufio.NewReader(file)
	// 		_, err = bufr.Read(bytes)

	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println(tx.SendTransaction(
	// 			[]sdk.Msg{zoracle.NewMsgStoreCode(bytes, "Crypto price (enum version)", tx.Sender())},
	// 			0, 10000000, "", "", "",
	// 			flags.BroadcastBlock,
	// 		))
	// 	}
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
	// case "request":
	// 	{
	// 		codeHash, _ := hex.DecodeString("148b6ddfdd2e1a6791b992592160ccd3cef0cea0c5f88ffdbdae7ea4044d9841")
	// 		switch args[1] {
	// 		case "BTC":
	// 			{
	// 				// BTC parameter
	// 				params, _ := hex.DecodeString("00000000")
	// 				// Send request by code hash
	// 				fmt.Println(tx.SendTransaction(
	// 					[]sdk.Msg{zoracle.NewMsgRequest(codeHash, params, 6, tx.Sender())},
	// 					0, 10000000, "", "", "",
	// 					flags.BroadcastBlock,
	// 				))
	// 			}
	// 		case "ETH":
	// 			{
	// 				// ETH parameter
	// 				params, _ := hex.DecodeString("00000001")
	// 				// Send request by same code hash with new parameter
	// 				fmt.Println(tx.SendTransaction(
	// 					[]sdk.Msg{zoracle.NewMsgRequest(codeHash, params, 6, tx.Sender())},
	// 					0, 10000000, "", "", "",
	// 					flags.BroadcastBlock,
	// 				))
	// 			}
	// 		}
	// 	}
	case "example":
		{
			file, err := os.Open("../../owasm/res/silly.wasm")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			stats, statsErr := file.Stat()
			if statsErr != nil {
				panic(statsErr)
			}

			size := stats.Size()
			bytes := make([]byte, size)

			bufr := bufio.NewReader(file)
			_, err = bufr.Read(bytes)

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
				[]sdk.Msg{zoracle.NewMsgRequestData(1, []byte("calldata"), 1, 1, 100, tx.Sender())},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))

			// _ = valTx.Sender()

			fmt.Println(valTx.SendTransaction(
				[]sdk.Msg{zoracle.NewMsgReportData(1, []zoracle.RawDataReport{
					zoracle.NewRawDataReport(1, []byte("data1")),
				}, sdk.ValAddress(valTx.Sender()))},
				0, 10000000, "", "", "",
				flags.BroadcastBlock,
			))
		}
	}
}
