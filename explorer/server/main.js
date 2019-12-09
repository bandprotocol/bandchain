// Server entry point, imports all server code

import '/imports/startup/server';
import '/imports/startup/both';
// import moment from 'moment';
// import '/imports/api/blocks/blocks.js';

SYNCING = false;
COUNTMISSEDBLOCKS = false;
COUNTMISSEDBLOCKSSTATS = false;
RPC = Meteor.settings.remote.rpc;
LCD = Meteor.settings.remote.lcd;
timerBlocks = 0;
timerChain = 0;
timerConsensus = 0;
timerProposal = 0;
timerProposalsResults = 0;
timerMissedBlock = 0;
timerDelegation = 0;
timerAggregate = 0;

const DEFAULTSETTINGS = '/default_settings.json';

updateChainStatus = () => {
    Meteor.call('chain.updateStatus', (error, result) => {
        if (error){
            console.log("updateStatus: "+error);
        }
        else{
            console.log("updateStatus: "+result);
        }
    })
}

updateBlock = () => {
    Meteor.call('blocks.blocksUpdate', (error, result) => {
        if (error){
            console.log("updateBlocks: "+error);
        }
        else{
            console.log("updateBlocks: "+result);
        }
    })
}

getConsensusState = () => {
    Meteor.call('chain.getConsensusState', (error, result) => {
        if (error){
            console.log("get consensus: "+error)
        }
    })
}

getProposals = () => {
    Meteor.call('proposals.getProposals', (error, result) => {
        if (error){
            console.log("get proposal: "+ error);
        }
        if (result){
            console.log("get proposal: "+result);
        }
    });
}

getProposalsResults = () => {
    Meteor.call('proposals.getProposalResults', (error, result) => {
        if (error){
            console.log("get proposals result: "+error);
        }
        if (result){
            console.log("get proposals result: "+result);
        }
    });
}

updateMissedBlocks = () => {
    Meteor.call('ValidatorRecords.calculateMissedBlocks', (error, result) =>{
        if (error){
            console.log("missed blocks error: "+ error)
        }
        if (result){
            console.log("missed blocks ok:" + result);
        }
    });
/*
    Meteor.call('ValidatorRecords.calculateMissedBlocksStats', (error, result) =>{
        if (error){
            console.log("missed blocks stats error: "+ error)
        }
        if (result){
            console.log("missed blocks stats ok:" + result);
        }
    });
*/
}

getDelegations = () => {
    Meteor.call('delegations.getDelegations', (error, result) => {
        if (error){
            console.log("get delegations error: "+ error)
        }
        else{
            console.log("get delegations ok: "+ result)
        }
    });
}

aggregateMinutely = () =>{
    // doing something every min
    Meteor.call('Analytics.aggregateBlockTimeAndVotingPower', "m", (error, result) => {
        if (error){
            console.log("aggregate minutely block time error: "+error)
        }
        else{
            console.log("aggregate minutely block time ok: "+result)
        }
    });

    Meteor.call('coinStats.getCoinStats', (error, result) => {
        if (error){
            console.log("get coin stats error: "+error);
        }
        else{
            console.log("get coin stats ok: "+result)
        }
    });
}

aggregateHourly = () =>{
    // doing something every hour
    Meteor.call('Analytics.aggregateBlockTimeAndVotingPower', "h", (error, result) => {
        if (error){
            console.log("aggregate hourly block time error: "+error)
        }
        else{
            console.log("aggregate hourly block time ok: "+result)
        }
    });
}

aggregateDaily = () =>{
    // doing somthing every day
    Meteor.call('Analytics.aggregateBlockTimeAndVotingPower', "d", (error, result) => {
        if (error){
            console.log("aggregate daily block time error: "+error)
        }
        else{
            console.log("aggregate daily block time ok: "+result)
        }
    });

    Meteor.call('Analytics.aggregateValidatorDailyBlockTime', (error, result) => {
        if (error){
            console.log("aggregate validators block time error:"+ error)
        }
        else {
            console.log("aggregate validators block time ok:"+ result);
        }
    })
}



Meteor.startup(function(){
    if (Meteor.isDevelopment){
        process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0;
        import DEFAULTSETTINGSJSON from '../default_settings.json'
        Object.keys(DEFAULTSETTINGSJSON).forEach((key) => {
            if (Meteor.settings[key] == undefined) {
                console.warn(`CHECK SETTINGS JSON: ${key} is missing from settings`)
                Meteor.settings[key] = {};
            }
            Object.keys(DEFAULTSETTINGSJSON[key]).forEach((param) => {
                if (Meteor.settings[key][param] == undefined){
                    console.warn(`CHECK SETTINGS JSON: ${key}.${param} is missing from settings`)
                    Meteor.settings[key][param] = DEFAULTSETTINGSJSON[key][param]
                }
            })
        })
    }

    Meteor.call('chain.genesis', (err, result) => {
        if (err){
            console.log(err);
        }
        if (result){
            if (Meteor.settings.debug.startTimer){
                timerConsensus = Meteor.setInterval(function(){
                    getConsensusState();
                }, Meteor.settings.params.consensusInterval);

                timerBlocks = Meteor.setInterval(function(){
                    updateBlock();
                }, Meteor.settings.params.blockInterval);

                timerChain = Meteor.setInterval(function(){
                    updateChainStatus();
                }, Meteor.settings.params.statusInterval);

                timerProposal = Meteor.setInterval(function(){
                    getProposals();
                }, Meteor.settings.params.proposalInterval);

                timerProposalsResults = Meteor.setInterval(function(){
                    getProposalsResults();
                }, Meteor.settings.params.proposalInterval);

                timerMissedBlock = Meteor.setInterval(function(){
                    updateMissedBlocks();
                }, Meteor.settings.params.missedBlocksInterval);

                timerDelegation = Meteor.setInterval(function(){
                    getDelegations();
                }, Meteor.settings.params.delegationInterval);

                timerAggregate = Meteor.setInterval(function(){
                    let now = new Date();
                    if ((now.getUTCSeconds() == 0)){
                        aggregateMinutely();
                    }

                    if ((now.getUTCMinutes() == 0) && (now.getUTCSeconds() == 0)){
                        aggregateHourly();
                    }

                    if ((now.getUTCHours() == 0) && (now.getUTCMinutes() == 0) && (now.getUTCSeconds() == 0)){
                        aggregateDaily();
                    }
                }, 1000)
            }
        }
    })

});
