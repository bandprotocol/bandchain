package db

import (
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

func (b *BandDB) handleMsgWithdrawDelegatorReward(msg distribution.MsgWithdrawDelegatorReward) error {
	info := b.DistrKeeper.GetDelegatorStartingInfo(b.ctx, msg.ValidatorAddress, msg.DelegatorAddress)
	latestReward := b.DistrKeeper.GetValidatorHistoricalRewards(b.ctx, msg.ValidatorAddress, info.PreviousPeriod)
	// CurrentReward must be reset after delegation.
	err := b.UpdateValidator(
		msg.ValidatorAddress,
		&Validator{
			CurrentReward: "0",
			CurrentRatio:  latestReward.CumulativeRewardRatio[0].Amount.String(),
		},
	)
	if err != nil {
		return err
	}
	var delegation Delegation
	err = b.tx.Where(Delegation{
		ValidatorAddress: msg.ValidatorAddress.String(),
		DelegatorAddress: msg.DelegatorAddress.String(),
	}).First(&delegation).Error
	if err != nil {
		return err
	}
	delegation.LastRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	return b.tx.Save(&delegation).Error
}
