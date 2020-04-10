package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (b *BandDB) HandleBeginblockEvent(event abci.Event) {
	kvMap := make(map[string]string)
	for _, kv := range event.Attributes {
		kvMap[string(kv.Key)] = string(kv.Value)
	}

	switch event.Type {
	case slashing.EventTypeSlash:
		{
			if rawConsAddress, ok := kvMap[slashing.AttributeKeyJailed]; ok {
				consAddress, err := sdk.ConsAddressFromBech32(rawConsAddress)
				if err != nil {
					panic(err)
				}
				validator, found := b.StakingKeeper.GetValidatorByConsAddr(b.ctx, consAddress)
				if !found {
					panic("HandleBeginblockEvent: validator not found")
				}
				err = b.tx.Model(&Validator{}).
					Where(Validator{ConsensusAddress: tmbytes.HexBytes(consAddress.Bytes()).String()}).
					Update(&Validator{
						Tokens: validator.Tokens.Uint64(),
						Jailed: true,
					}).Error
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
