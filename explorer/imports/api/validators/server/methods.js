import { Meteor } from 'meteor/meteor';
import { Transactions } from '../../transactions/transactions.js';
import { Blockscon } from '../../blocks/blocks.js';
import { Delegations } from '../../delegations/delegations.js';

Meteor.methods({
    'Validators.findCreateValidatorTime': function(address){
        // look up the create validator time to consider if the validator has never updated the commission
        let tx = Transactions.findOne({$and:[
            {"tx.value.msg.value.delegator_address":address},
            {"tx.value.msg.type":"cosmos-sdk/MsgCreateValidator"},
            {code:{$exists:false}}
        ]});

        if (tx){
            let block = Blockscon.findOne({height:tx.height});
            if (block){
                return block.time;
            }
        }
        else{
            // no such create validator tx
            return false;
        }
    },
    // async 'Validators.getAllDelegations'(address){
    'Validators.getAllDelegations'(address){
        let url = LCD + '/staking/validators/'+address+'/delegations';

        try{
            let delegations = HTTP.get(url);
            if (delegations.statusCode == 200){
                delegations = JSON.parse(delegations.content).result;
                delegations.forEach((delegation, i) => {
                    if (delegations[i] && delegations[i].shares)
                        delegations[i].shares = parseFloat(delegations[i].shares);
                })
                
                return delegations;
            };
        }
        catch (e){
            console.log(e);
        }
    }
});