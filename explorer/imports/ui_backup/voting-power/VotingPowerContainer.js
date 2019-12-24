import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Validators } from '/imports/api/validators/validators.js';
import VotingPower from './VotingPower.jsx';

export default VotingPowerContainer = withTracker((props) => {
    let chartHandle, stats, statsExist
    let loading = true;

    if (Meteor.isClient){
        chartHandle = Meteor.subscribe('validators.voting_power');
        loading = !chartHandle.ready();
    }

    if (Meteor.isServer || !loading){
        stats = Validators.find({},{sort:{voting_power:-1}}).fetch();

        if (Meteor.isServer){
            // loading = false;
            statsExist = !!stats;
        }
        else{
            statsExist = !loading && !!stats;
        }
    }

    return {
        loading,
        statsExist,
        stats: statsExist ? stats : {}
    };
})(VotingPower);

