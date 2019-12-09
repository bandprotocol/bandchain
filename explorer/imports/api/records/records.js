import { Mongo } from 'meteor/mongo';
import { Validators } from '../validators/validators';

export const ValidatorRecords = new Mongo.Collection('validator_records');
export const Analytics = new Mongo.Collection('analytics');
export const MissedBlocksStats = new Mongo.Collection('missed_blocks_stats');
export const MissedBlocks = new  Mongo.Collection('missed_blocks');
export const VPDistributions = new Mongo.Collection('voting_power_distributions');
export const AverageData = new Mongo.Collection('average_data');
export const AverageValidatorData = new Mongo.Collection('average_validator_data');

MissedBlocksStats.helpers({
    proposerMoniker(){
        let validator = Validators.findOne({address:this.proposer});
        return (validator.description)?validator.description.moniker:this.proposer;
    },
    voterMoniker(){
        let validator = Validators.findOne({address:this.voter});
        return (validator.description)?validator.description.moniker:this.voter;
    }
})
