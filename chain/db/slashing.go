package db

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

func (b *BandDB) handleMsgUnjail(msg slashing.MsgUnjail) error {
	jailed := false
	return b.tx.Model(&Validator{}).
		Where(Validator{OperatorAddress: msg.ValidatorAddr.String()}).
		Update(Validator{Jailed: &jailed}).Error
}
