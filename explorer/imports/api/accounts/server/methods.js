import { Meteor } from 'meteor/meteor';
import { HTTP } from 'meteor/http';
import { Validators } from '/imports/api/validators/validators.js';
const fetchFromUrl = (url) => {
    try{
        let res = HTTP.get(LCD + url);
        if (res.statusCode == 200){
            return res
        };
    }
    catch (e){
        console.log(e);
    }
}

Meteor.methods({
    'accounts.getAccountDetail': function(address){
        this.unblock();
        let url = LCD + '/auth/accounts/'+ address;
        try{
            let available = HTTP.get(url);
            if (available.statusCode == 200){
                let response = JSON.parse(available.content).result;
                let account;
                if (response.type === 'cosmos-sdk/Account')
                    account = response.value;
                else if (response.type === 'cosmos-sdk/DelayedVestingAccount' || response.type === 'cosmos-sdk/ContinuousVestingAccount')
                    account = response.value.BaseVestingAccount.BaseAccount
                if (account && account.account_number != null)
                    return account
                return null
            }
        }
        catch (e){
            console.log(e)
        }
    },
    'accounts.getBalance': function(address){
        this.unblock();
        let balance = {}

        // get available atoms
        let url = LCD + '/bank/balances/'+ address;
        try{
            let available = HTTP.get(url);
            if (available.statusCode == 200){
                // console.log(JSON.parse(available.content))
                balance.available = JSON.parse(available.content).result;
                if (balance.available && balance.available.length > 0)
                    balance.available = balance.available[0];
            }
        }
        catch (e){
            console.log(e)
        }

        // get delegated amnounts
        url = LCD + '/staking/delegators/'+address+'/delegations';
        try{
            let delegations = HTTP.get(url);
            if (delegations.statusCode == 200){
                balance.delegations = JSON.parse(delegations.content).result;
            }
        }
        catch (e){
            console.log(e);
        }
        // get unbonding
        url = LCD + '/staking/delegators/'+address+'/unbonding_delegations';
        try{
            let unbonding = HTTP.get(url);
            if (unbonding.statusCode == 200){
                balance.unbonding = JSON.parse(unbonding.content).result;
            }
        }
        catch (e){
            console.log(e);
        }

        // get rewards
        url = LCD + '/distribution/delegators/'+address+'/rewards';
        try{
            let rewards = HTTP.get(url);
            if (rewards.statusCode == 200){
                balance.rewards = JSON.parse(rewards.content).result.total;
            }
        }
        catch (e){
            console.log(e);
        }

        // get commission
        let validator = Validators.findOne(
            {$or: [{operator_address:address}, {delegator_address:address}, {address:address}]})
        if (validator) {
            let url = LCD + '/distribution/validators/' + validator.operator_address;
            balance.operator_address = validator.operator_address;
            try {
                let rewards = HTTP.get(url);
                if (rewards.statusCode == 200){
                    let content = JSON.parse(rewards.content).result;
                    if (content.val_commission && content.val_commission.length > 0)
                        balance.commission = content.val_commission[0];
                }

            }
            catch (e){
                console.log(e)
            }
        }

        return balance;
    },
    'accounts.getDelegation'(address, validator){
        let url = `/staking/delegators/${address}/delegations/${validator}`;
        let delegations = fetchFromUrl(url);
        delegations = delegations && delegations.data.result;
        if (delegations && delegations.shares)
            delegations.shares = parseFloat(delegations.shares);

        url = `/staking/redelegations?delegator=${address}&validator_to=${validator}`;
        let relegations = fetchFromUrl(url);
        relegations = relegations && relegations.data.result;
        let completionTime;
        if (relegations) {
            relegations.forEach((relegation) => {
                let entries = relegation.entries
                let time = new Date(entries[entries.length-1].completion_time)
                if (!completionTime || time > completionTime)
                    completionTime = time
            })
            delegations.redelegationCompletionTime = completionTime;
        }

        url = `/staking/delegators/${address}/unbonding_delegations/${validator}`;
        let undelegations = fetchFromUrl(url);
        undelegations = undelegations && undelegations.data.result;
        if (undelegations) {
            delegations.unbonding = undelegations.entries.length;
            delegations.unbondingCompletionTime = undelegations.entries[0].completion_time;
        }
        return delegations;
    },
    'accounts.getAllDelegations'(address){
        let url = LCD + '/staking/delegators/'+address+'/delegations';

        try{
            let delegations = HTTP.get(url);
            if (delegations.statusCode == 200){
                delegations = JSON.parse(delegations.content).result;
                if (delegations && delegations.length > 0){
                    delegations.forEach((delegation, i) => {
                        if (delegations[i] && delegations[i].shares)
                            delegations[i].shares = parseFloat(delegations[i].shares);
                    })
                }

                return delegations;
            };
        }
        catch (e){
            console.log(e);
        }
    },
    'accounts.getAllUnbondings'(address){
        let url = LCD + '/staking/delegators/'+address+'/unbonding_delegations';

        try{
            let unbondings = HTTP.get(url);
            if (unbondings.statusCode == 200){
                unbondings = JSON.parse(unbondings.content).result;
                return unbondings;
            };
        }
        catch (e){
            console.log(e);
        }
    },
    'accounts.getAllRedelegations'(address, validator){
        let url = `/staking/redelegations?delegator=${address}&validator_from=${validator}`;
        let result = fetchFromUrl(url);
        if (result && result.data) {
            let redelegations = {}
            result.data.forEach((redelegation) => {
                let entries = redelegation.entries;
                redelegations[redelegation.validator_dst_address] = {
                    count: entries.length,
                    completionTime: entries[0].completion_time
                }
            })
            return redelegations
        }
    }
})