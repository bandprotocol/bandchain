import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { VPDistributions } from '/imports/api/records/records.js';
import TwentyEighty from './TwentyEighty.jsx';

export default TwentyEightyContainer = withTracker((props) => {
    let chartHandle, stats, statsExist;
    let loading = true;
    
    if (Meteor.isClient){
        chartHandle = Meteor.subscribe('vpDistribution.latest');
        loading = !chartHandle.ready();
    }

    if (Meteor.isServer || !loading){
        stats = VPDistributions.findOne({});

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
})(TwentyEighty);

