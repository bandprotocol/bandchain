import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Proposals } from '/imports/api/proposals/proposals.js';
import List from './List.jsx';

export default ProposalListContainer = withTracker((props) => {
    let proposalsHandle, proposals, proposalsExist;
    let loading = true;

    if (Meteor.isClient){
        proposalsHandle = Meteor.subscribe('proposals.list');
        loading = !proposalsHandle.ready();
    }

    if (Meteor.isServer || !loading){
        proposals = Proposals.find({}, {sort:{proposalId:-1}}).fetch();

        if (Meteor.isServer){
            loading = false;
            proposalsExist = !!proposals;
        }
        else{
            proposalsExist = !loading && !!proposals;
        }
    }

    return {
        loading,
        proposalsExist,
        proposals: proposalsExist ? proposals : {},
        history: props.history
    };
})(List);
