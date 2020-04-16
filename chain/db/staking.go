package db

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (b *BandDB) handleMsgCreateValidator(msg staking.MsgCreateValidator) error {
	// pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msg.PubKey)
	// if err != nil {
	// 	return err
	// }
	err := b.CreateValidator(
		msg.ValidatorAddress,
		msg.PubKey,
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
		sdk.NewDecCoins(sdk.NewCoins()),
	)
}

func (b *BandDB) handleMsgEditValidator(msg staking.MsgEditValidator) error {
	validator, isFound := b.GetValidator(msg.ValidatorAddress)
	if !isFound {
		return fmt.Errorf(fmt.Sprintf("validator %s does not exist.", msg.ValidatorAddress.String()))
	}

	if msg.Description.Moniker != staking.DoNotModifyDesc {
		validator.Moniker = msg.Description.Moniker
	}
	if msg.Description.Identity != staking.DoNotModifyDesc {
		validator.Identity = msg.Description.Identity
	}
	if msg.Description.Website != staking.DoNotModifyDesc {
		validator.Website = msg.Description.Website
	}
	if msg.Description.Details != staking.DoNotModifyDesc {
		validator.Details = msg.Description.Details
	}
	if msg.CommissionRate != nil {
		validator.CommissionRate = msg.CommissionRate.String()
	}
	if msg.MinSelfDelegation != nil {
		validator.MinSelfDelegation = msg.MinSelfDelegation.ToDec().String()
	}

	return b.tx.Save(&validator).Error
}

func (b *BandDB) handleMsgDelegate(msg staking.MsgDelegate) error {
	return b.delegate(msg.DelegatorAddress, msg.ValidatorAddress)
}

func (b *BandDB) handleMsgUndelegate(msg staking.MsgUndelegate) error {
	return b.undelegate(msg.DelegatorAddress, msg.ValidatorAddress)
}

func (b *BandDB) handleMsgBeginRedelegate(msg staking.MsgBeginRedelegate) error {
	err := b.undelegate(msg.DelegatorAddress, msg.ValidatorSrcAddress)
	if err != nil {
		return err
	}
	return b.delegate(msg.DelegatorAddress, msg.ValidatorDstAddress)
}
