package main

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/bandlib"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

var bandClient bandlib.BandStatefulClient

func init() {
	// TODO: Cleanup
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	priv, err := hex.DecodeString(
		"6b0b8909eadbbc220797dc0aada9558030d4a89972e51de8de525fc9de42bd40",
	)
	var privSecp256k1 secp256k1.PrivKeySecp256k1
	copy(privSecp256k1[:], priv)
	bandClient, err = bandlib.NewBandStatefulClient(
		"tcp://localhost:26657", privSecp256k1, 100, 10, "Bandoracled reports", "bandchain",
	)
	if err != nil {
		panic(err)
	}
}

func GetExecutable(dataSourceID int) ([]byte, error) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query(
		fmt.Sprintf("custom/oracle/%s/%d", oracle.QueryDataSourceByID, dataSourceID),
	)

	if err != nil {
		return nil, err
	}

	var dataSource oracle.DataSourceQuerierInfo
	err = cliCtx.Codec.UnmarshalJSON(res, &dataSource)
	if err != nil {
		return nil, err
	}
	return dataSource.Executable, nil
}
