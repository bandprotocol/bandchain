import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Proposals } from '/imports/api/proposals/proposals.js';
import { Chain } from '/imports/api/chain/chain.js';
import Proposal from './Proposal.jsx';

export default ProposalContainer = withTracker((props) => {
    let proposalId = 0;
    if (props.match.params.id){
        proposalId = parseInt(props.match.params.id);
    }

    let chainHandle, proposalHandle, proposalListHandle, proposal, proposalCount, chain, proposalExist;
    let loading = true;

    if (Meteor.isClient){
        chainHandle = Meteor.subscribe('chain.status');
        proposalListHandle = Meteor.subscribe('proposals.list', proposalId);
        proposalHandle = Meteor.subscribe('proposals.one', proposalId);
        loading = !proposalHandle.ready() || !chainHandle.ready() || !proposalListHandle.ready();
    }

    if (Meteor.isServer || !loading){
        proposal = Proposals.findOne({proposalId:proposalId});
        proposalCount = Proposals.find({}).count();
        chain = Chain.findOne({chainId:Meteor.settings.public.chainId});

        if (Meteor.isServer){
            // loading = false;
            proposalExist = !!proposal;
        }
        else{
            proposalExist = !loading && !!proposal;
        }
    }

    return {
        loading,
        proposalExist,
        proposal: proposalExist ? proposal : {},
        chain: proposalExist ? chain : {},
        proposalCount: proposalExist? proposalCount: 0
    };
})(Proposal);
