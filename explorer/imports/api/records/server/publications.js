import { Meteor } from 'meteor/meteor';
import { ValidatorRecords, Analytics, MissedBlocks, MissedBlocksStats, VPDistributions } from '../records.js';
import { Validators } from '../../validators/validators.js';

Meteor.publish('validator_records.all', function () {
    return ValidatorRecords.find();
});

Meteor.publish('validator_records.uptime', function(address, num){
    return ValidatorRecords.find({address:address},{limit:num, sort:{height:-1}});
});

Meteor.publish('analytics.history', function(){
    return Analytics.find({},{sort:{height:-1},limit:50});
});

Meteor.publish('vpDistribution.latest', function(){
    return VPDistributions.find({},{sort:{height:-1}, limit:1});
});

publishComposite('missedblocks.validator', function(address, type){
    let conditions = {};
    if (type == 'voter'){
        conditions = {
            voter: address
        }
    }
    else{
        conditions = {
            proposer: address
        }
    }
    return {
        find(){
            return MissedBlocksStats.find(conditions)
        },
        children: [
            {
                find(stats){
                    return Validators.find(
                        {},
                        {fields:{address:1, description:1, profile_url:1}}
                    )
                }
            }
        ]
    }
});

publishComposite('missedrecords.validator', function(address, type){
    return {
        find(){
            return MissedBlocks.find(
                {[type]: address},
                {sort: {updatedAt: -1}}
            )
        },
        children: [
            {
                find(){
                    return Validators.find(
                        {},
                        {fields:{address:1, description:1, operator_address:1}}
                    )
                }
            }
        ]
    }
});
