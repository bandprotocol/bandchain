package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
	viper.Set("nodeURI", nodeURI)

	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	tx := cmtx.NewTxSender(priv)

	args := os.Args[1:]
	switch args[0] {
	case "store":
		{
			file, err := os.Open("../../wasm/res/result.wasm")
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
			fmt.Println(tx.SendTransaction(zoracle.NewMsgStoreCode(bytes, "Crypto price (enum version)", tx.Sender()), flags.BroadcastBlock))
		}
	case "send_token":
		{
			// Send token
			to, _ := sdk.AccAddressFromBech32("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
			fmt.Println(tx.SendTransaction(bank.MsgSend{
				FromAddress: tx.Sender(),
				ToAddress:   to,
				Amount:      sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(10))),
			}, flags.BroadcastBlock))
		}
	case "request":
		{
			codeHash, _ := hex.DecodeString("c5d9b37939c8a4f2bb55a14ebb7bc9138e8f29485bf8fcbacff49210ad66d3dc")
			switch args[1] {
			case "BTC":
				{
					// BTC parameter
					params, _ := hex.DecodeString("00000000")
					// Send request by code hash
					fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))
				}
			case "ETH":
				{
					// ETH parameter
					params, _ := hex.DecodeString("00000001")
					// Send request by same code hash with new parameter
					fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))
				}
			}
		}
	}
}
