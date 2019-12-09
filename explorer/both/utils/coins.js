import { Meteor } from 'meteor/meteor';
import numbro from 'numbro';

autoformat = (value) => {
	let formatter = '0,0.0000';
	value = Math.round(value * 1000) / 1000
	if (Math.round(value) === value)
		formatter = '0,0'
	else if (Math.round(value*10) === value*10)
		formatter = '0,0.0'
	else if (Math.round(value*100) === value*100)
		formatter = '0,0.00'
	else if (Math.round(value*1000) === value*1000)
		formatter = '0,0.000'
	return numbro(value).format(formatter)
}

export default class Coin {
	static StakingDenom = Meteor.settings.public.stakingDenom;
	static StakingDenomPlural = Meteor.settings.public.stakingDenomPlural || (Coin.StakingDenom + 's');
	static MintingDenom = Meteor.settings.public.mintingDenom;
	static StakingFraction = Number(Meteor.settings.public.stakingFraction);
	static MinStake = 1 / Number(Meteor.settings.public.stakingFraction);

	constructor(amount, denom=null) {
		if (typeof amount === 'object')
			({amount, denom} = amount)
		if (!denom || denom.toLowerCase() === Coin.MintingDenom.toLowerCase()) {
			this._amount = Number(amount);
		} else if (denom.toLowerCase() === Coin.StakingDenom.toLowerCase()) {
			this._amount = Number(amount) * Coin.StakingFraction;
		}
		else {
			throw Error(`unsupported denom ${denom}`);
		}
	}

	get amount () {
		return this._amount;
	}

	get stakingAmount () {
		return this._amount / Coin.StakingFraction;
	}

	toString (precision) {
		// default to display in mint denom if it has more than 4 decimal places
		let minStake = Coin.StakingFraction/(precision?Math.pow(10, precision):10000)
		if (this.amount < minStake) {
			return `${numbro(this.amount).format('0,0')} ${Coin.MintingDenom}`;
		} else {
			return `${precision?numbro(this.stakingAmount).format('0,0.' + '0'.repeat(precision)):autoformat(this.stakingAmount)} ${Coin.StakingDenom}`
		}
	}

	mintString (formatter) {
		let amount = this.amount
		if (formatter) {
			amount = numbro(amount).format(formatter)
		}
		return `${amount} ${Coin.MintingDenom}`;
	}

	stakeString (formatter) {
		let amount = this.stakingAmount
		if (formatter) {
			amount = numbro(amount).format(formatter)
		}
		return `${amount} ${Coin.StakingDenom}`;
	}
}