package db

import (
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

func (b *BandDB) handleMsgWithdrawDelegatorReward(msg distribution.MsgWithdrawDelegatorReward) error {
	info := b.DistrKeeper.GetDelegatorStartingInfo(b.ctx, msg.ValidatorAddress, msg.DelegatorAddress)
	latestReward := b.DistrKeeper.GetValidatorHistoricalRewards(b.ctx, msg.ValidatorAddress, info.PreviousPeriod)
	// CurrentReward must be reset after delegation.
	cumulativeRewardRatio := "0"
	if !latestReward.CumulativeRewardRatio.IsZero() {
		cumulativeRewardRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	}
	err := b.UpdateValidator(
		msg.ValidatorAddress,
		&Validator{
			CurrentReward: "0",
			CurrentRatio:  cumulativeRewardRatio,
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
	delegation.LastRatio = cumulativeRewardRatio
	return b.tx.Save(&delegation).Error
}
