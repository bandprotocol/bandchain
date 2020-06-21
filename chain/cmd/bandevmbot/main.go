package main

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	crypto "github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/spf13/cobra"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// 	"github.com/tendermint/tendermint/libs/log"
// 	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
// )

// var (
// 	logger log.Logger
// )

// type ValidatorWithPower struct {
// 	Addr  common.Address
// 	Power *big.Int
// }

// func getValidators(nodeURI string) []ValidatorWithPower {
// 	node, err := rpchttp.New(nodeURI, "/websocket")
// 	if err != nil {
// 		panic(err)
// 	}
// 	validators, err := node.Validators(nil, 0, 0)
// 	if err != nil {
// 		panic(err)
// 	}

// 	vals := make([]ValidatorWithPower, len(validators.Validators))

// 	for idx, validator := range validators.Validators {
// 		pubKeyBytes, ok := validator.PubKey.(secp256k1.PubKeySecp256k1)
// 		if !ok {
// 			panic("fail to cast pubkey")
// 		}

// 		pubkey, err := crypto.DecompressPubkey(pubKeyBytes[:])
// 		if err != nil {
// 			panic(err)
// 		}
// 		vals[idx] = ValidatorWithPower{
// 			Addr:  crypto.PubkeyToAddress(*pubkey),
// 			Power: big.NewInt(validator.VotingPower),
// 		}

// 	}
// 	return vals
// }

// func updateValidators(rpcURI string, address string, node string, privateKey string, gasPrice uint64) {
// 	vals := getValidators(node)
// 	backgroundCtx := context.Background()
// 	contractAddress := common.HexToAddress(address)
// 	priv, err := crypto.HexToECDSA(privateKey)
// 	botAddress := crypto.PubkeyToAddress(priv.PublicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	evmClient, err := ethclient.Dial(rpcURI)
// 	if err != nil {
// 		panic(err)
// 	}
// 	abiJson, err := abi.JSON(bytes.NewReader(rawABI))
// 	if err != nil {
// 		panic(err)
// 	}
// 	data, err := abiJson.Pack("updateValidatorPowers", vals)
// 	if err != nil {
// 		panic(err)
// 	}
// 	chainID, err := evmClient.ChainID(backgroundCtx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	nonce, err := evmClient.NonceAt(backgroundCtx, botAddress, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tx := ethtypes.NewTransaction(
// 		nonce,
// 		contractAddress,
// 		big.NewInt(0),
// 		gasPrice,
// 		big.NewInt(2000000000),
// 		data,
// 	)

// 	signer := ethtypes.NewEIP155Signer(chainID)
// 	signTx, err := ethtypes.SignTx(tx, signer, priv)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = evmClient.SendTransaction(context.Background(), signTx)
// 	if err != nil {
// 		panic(err)
// 	}

// }

// const (
// 	flagRPCUri          = "rpc-uri"
// 	flagContractAddress = "contract-address"
// 	flagNodeUri         = "node-uri"
// 	flagPrivKey         = "priv-key"
// 	flagGasPrice        = "gas-price"
// 	flagPollInterval    = "poll-interval"
// )

func main() {

	// 	cmd := &cobra.Command{
	// 		Use:   "(--rpc-uri [rpc-uri]) (--contract-address [contract-address]) (--node-uri [node-uri]) (--priv-key [priv-key]) (--gas-price [gas-price]) (--poll-interval [poll-interval])",
	// 		Short: "Periodically update validator set to the destination EVM blockchain",
	// 		Args:  cobra.ExactArgs(0),
	// 		Long: strings.TrimSpace(
	// 			fmt.Sprintf(`Periodically update validator set to the destination EVM blockchain
	// Example:
	// $ --rpc-uri https://kovan.infura.io/v3/d3301689638b40dabad8395bf00d3945 --contract-address 0x0d8152D22a05A3Cf2cE1c5bEfCc2F8658f75a59d --node-uri http://d3n-debug.bandprotocol.com:26657 --priv-key AA0C65C16D4B8511C58122966F94192F6963D0EB7896435430BCDFF56E9F13B9 --gas-price 1000000 --poll-interval 24
	// `),
	// 		),
	// 		RunE: func(cmd *cobra.Command, args []string) error {

	// 			rpcURI, err := cmd.Flags().GetString(flagRPCUri)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			contractAddress, err := cmd.Flags().GetString(flagContractAddress)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			nodeURI, err := cmd.Flags().GetString(flagNodeUri)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			privateKey, err := cmd.Flags().GetString(flagPrivKey)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			rawGasPrice, err := cmd.Flags().GetString(flagGasPrice)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			gasPrice, err := strconv.ParseInt(rawGasPrice, 10, 64)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			rawInterval, err := cmd.Flags().GetString(flagPollInterval)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			interval, err := strconv.ParseInt(rawInterval, 10, 64)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			for {
	// 				updateValidators(rpcURI, contractAddress, nodeURI, privateKey, uint64(gasPrice))
	// 				fmt.Println("finish round")
	// 				time.Sleep(time.Duration(interval) * time.Hour)
	// 			}
	// 		},
	// 	}
	// 	cmd.Flags().String(flagRPCUri, "", "RPC URI")
	// 	cmd.Flags().String(flagContractAddress, "", "Address of contract")
	// 	cmd.Flags().String(flagNodeUri, "", "Node URI")
	// 	cmd.Flags().String(flagPrivKey, "", "Private key")
	// 	cmd.Flags().String(flagGasPrice, "", "Gas Price")
	// 	cmd.Flags().String(flagPollInterval, "", "Interval of update validatos (Hours)")

	// 	err := cmd.Execute()
	// 	if err != nil {
	// 		logger.Error(fmt.Sprintf("Failed executing: %s, exiting...", err))
	// 		os.Exit(1)

	// 	}
}
