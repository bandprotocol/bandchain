package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const Bech32MainPrefix = "band"

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32MainPrefix, Bech32MainPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
}

// Script for get private key from local database
func main() {
	kb, _ := keys.NewKeyBaseFromDir(os.ExpandEnv("$HOME/.bandcli"))
	privKey, _ := kb.ExportPrivateKeyObject("your_key_name", "your_password")
	fmt.Printf("%x", privKey)

	pkBytes, _ := hex.DecodeString("03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
	var pk secp256k1.PubKeySecp256k1
	copy(pk[:], pkBytes)
	config := sdk.GetConfig()
	SetBech32AddressPrefixes(config)
	pkStr := sdk.MustBech32ifyConsPub(pk)
	fmt.Println(pkStr)
}
