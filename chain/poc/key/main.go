package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/keys"
)

// Script for get private key from local database
func main() {
	kb, _ := keys.NewKeyBaseFromDir(os.ExpandEnv("$HOME/.bandcli"))
	privKey, _ := kb.ExportPrivateKeyObject("your_key_name", "your_password")
	fmt.Printf("%x", privKey)
}
