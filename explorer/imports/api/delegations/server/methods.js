import { Meteor } from 'meteor/meteor';
import { Delegations } from '../delegations.js';
import { Validators } from '../../validators/validators.js';

Meteor.methods({
    'delegations.getDelegations': function(){
        this.unblock();
        let validators = Validators.find({}).fetch();
        let delegations = [];
        console.log("=== Getting delegations ===");
        for (v in validators){
            if (validators[v].operator_address){
                let url = LCD + '/staking/validators/'+validators[v].operator_address+"/delegations";
                try{
                    let response = HTTP.get(url);
                    if (response.statusCode == 200){
                        let delegation = JSON.parse(response.content).result;
                        // console.log(delegation);
                        delegations = delegations.concat(delegation);
                    }
                    else{
                        console.log(response.statusCode);
                    }
                }
                catch (e){
                    console.log(e);
                }    
            }
        }

        for (i in delegations){
            if (delegations[i] && delegations[i].shares)
                delegations[i].shares = parseFloat(delegations[i].shares);
        }

        // console.log(delegations);
        let data = {
            delegations: delegations,
            createdAt: new Date(),
        }

        return Delegations.insert(data);
    }
    // 'blocks.averageBlockTime'(address){
    //     let blocks = Blockscon.find({proposerAddress:address}).fetch();
    //     let heights = blocks.map((block, i) => {
    //         return block.height;
    //     });
    //     let blocksStats = Analytics.find({height:{$in:heights}}).fetch();
    //     // console.log(blocksStats);

    //     let totalBlockDiff = 0;
    //     for (b in blocksStats){
    //         totalBlockDiff += blocksStats[b].timeDiff;
    //     }
    //     return totalBlockDiff/heights.length;
    // }
})