package db

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (b *BandDB) handleMsgCreateValidator(msg staking.MsgCreateValidator) error {
	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msg.Pubkey)
	if err != nil {
		return err
	}
	err = b.CreateValidator(
		msg.ValidatorAddress,
		pubkey,
		msg.Description.Moniker,
		msg.Description.Identity,
		msg.Description.Website,
		msg.Description.Details,
		msg.Commission.Rate,
		msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate,
		msg.MinSelfDelegation,
		msg.Value,
		b.ctx.BlockHeight(),
	)
	if err != nil {
		return err
	}

	return b.SetDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Value.Amount.String(),
		sdk.NewDecCoins(),
	)
}

func (b *BandDB) handleMsgEditValidator(msg staking.MsgEditValidator) error {
	validator, isFound := b.GetValidator(msg.ValidatorAddress)
	if !isFound {
		return fmt.Errorf(fmt.Sprintf("validator %s does not exist.", msg.ValidatorAddress.String()))
	}

	commissionRate := msg.CommissionRate.String()
	minSelfDelegation := msg.MinSelfDelegation.ToDec().String()
	if msg.Description.Moniker != staking.DoNotModifyDesc {
		validator.Moniker = &msg.Description.Moniker
	}
	if msg.Description.Identity != staking.DoNotModifyDesc {
		validator.Identity = &msg.Description.Identity
	}
	if msg.Description.Website != staking.DoNotModifyDesc {
		validator.Website = &msg.Description.Website
	}
	if msg.Description.Details != staking.DoNotModifyDesc {
		validator.Details = &msg.Description.Details
	}
	if msg.CommissionRate != nil {
		validator.CommissionRate = &commissionRate
	}
	if msg.MinSelfDelegation != nil {
		validator.MinSelfDelegation = &minSelfDelegation
	}

	return b.tx.Save(&validator).Error
}

func (b *BandDB) handleMsgDelegate(msg staking.MsgDelegate) error {
	return b.delegate(msg.DelegatorAddress, msg.ValidatorAddress)
}

func (b *BandDB) handleMsgUndelegate(msg staking.MsgUndelegate) error {
	err := b.undelegate(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return err
	}
	return b.updateUnbondingDelegations(msg.DelegatorAddress, msg.ValidatorAddress)
}

func (b *BandDB) handleMsgBeginRedelegate(msg staking.MsgBeginRedelegate) error {
	err := b.undelegate(msg.DelegatorAddress, msg.ValidatorSrcAddress)
	if err != nil {
		return err
	}
	return b.delegate(msg.DelegatorAddress, msg.ValidatorDstAddress)
}

func (b *BandDB) updateUnbondingDelegations(del sdk.AccAddress, val sdk.ValAddress) error {
	// Delete all records
	err := b.db.Delete(
		&UnbondingDelegation{},
		"delegator_address = ? AND validator_address = ?",
		del.String(), val.String(),
	).Error
	if err != nil {
		return err
	}

	unbondingList, found := b.StakingKeeper.GetUnbondingDelegation(b.ctx, del, val)
	if !found {
		// All unbonding delegation have been unbonded
		return nil
	}
	for _, ud := range unbondingList.Entries {
		balance := ud.Balance.Uint64()
		initBalance := ud.InitialBalance.Uint64()
		err := b.db.Create(&UnbondingDelegation{
			DelegatorAddress: del.String(),
			ValidatorAddress: val.String(),
			Balance:          &balance,
			InitialBalance:   &initBalance,
			CompletionTime:   ud.CompletionTime.UnixNano() / int64(time.Millisecond),
			CreationHeight:   ud.CreationHeight,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BandDB) updateUnbondingDelegationsOfValidator(val sdk.ValAddress) error {
	dels := b.StakingKeeper.GetUnbondingDelegationsFromValidator(b.ctx, val)
	for _, del := range dels {
		err := b.updateUnbondingDelegations(del.DelegatorAddress, val)
		if err != nil {
			return err
		}
	}
	return nil
}
