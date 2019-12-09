import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Chain } from '/imports/api/chain/chain.js';
import Consensus from './Consensus.jsx';

export default ConsensusContainer = withTracker((curr) => {
    let consensusHandle;
    let loading = true;
    let consensus;

    if (Meteor.isClient){
        consensusHandle = Meteor.subscribe('chain.status');
        loading = !consensusHandle.ready();    
    }
    
    let consensusExist;

    if (Meteor.isServer || !loading){
        consensus = Chain.findOne({chainId:Meteor.settings.public.chainId});

        if (Meteor.isServer){
            // loading = false;
            consensusExist = !!consensus;
        }   
        else{
            consensusExist = !loading && !!consensus;
        }
    }

    
    // console.log(props.state.limit);
    return {
        loading,
        consensusExist,
        consensus: consensusExist ? consensus : {}
    };
})(Consensus);

