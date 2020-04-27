package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestBCDExport(t *testing.T) {
	db := db.NewMemDB()
	bcapp := NewBandConsumerApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")
	err := setGenesis(bcapp)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newBCapp := NewBandConsumerApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")
	_, _, err = newBCapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := db.NewMemDB()
	bcapp := NewBandConsumerApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")

	for acc := range maccPerms {
		require.True(t, bcapp.bankKeeper.BlacklistedAddr(bcapp.supplyKeeper.GetModuleAddress(acc)))
	}
}

func setGenesis(bcapp *BandConsumerApp) error {
	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(bcapp.Codec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	bcapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	bcapp.Commit()
	return nil
}
