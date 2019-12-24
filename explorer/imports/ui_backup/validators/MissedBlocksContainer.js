import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Validators } from '/imports/api/validators/validators.js';
import { Status } from '/imports/api/status/status.js';
import { MissedBlocksStats, MissedBlocks } from '../../api/records/records.js';
import MissedBlocksComponent from './MissedBlocks.jsx';

export default MissedBlocksContainer = withTracker((props) => {
    let statusHandle;
    let validatorsHandle;
    // let missedBlockHandle;
    let loading = true;
    let address = props.match.params.address;

    if (Meteor.isClient){
        statusHandle = Meteor.subscribe('status.status');
        validatorsHandle = Meteor.subscribe('validator.details', address);

        if (props.type == 'voter'){
            // missedBlockHandle = Meteor.subscribe('missedblocks.validator', address, 'voter');
            missedRecordHandle = Meteor.subscribe('missedrecords.validator', address, 'voter');
        }
        else{
            // missedBlockHandle = Meteor.subscribe('missedblocks.validator', address, 'proposer');
            missedRecordHandle = Meteor.subscribe('missedrecords.validator', address, 'proposer');
        }

        // loading = !validatorsHandle.ready() && !statusHandle.ready() && !missedBlockHandle.ready();
        loading = !validatorsHandle.ready() && !statusHandle.ready() &&!missedRecordHandle.ready();
    }


    let validator;
    let status;
    let missedBlocks;
    let validatorExist;
    let statusExist;
    let missedBlocksExist;
    let missedRecords;

    if (Meteor.isServer || !loading){
        validator = Validators.findOne({address:address});
        status = Status.findOne({chainId:Meteor.settings.public.chainId});
        if (props.type == 'voter'){
            // missedBlocks = MissedBlocksStats.find({voter:address}, {sort:{count:-1}}).fetch();
            missedRecordsStats = MissedBlocks.find({voter:address, blockHeight: -1}).fetch();
            missedRecords = MissedBlocks.find({voter:address, blockHeight: {'$gt': 0}}, {sort:{blockHeight:-1}}).fetch();
        }
        else {
            // missedBlocks = MissedBlocksStats.find({proposer:address}, {sort:{count:-1}}).fetch();
            missedRecordsStats = MissedBlocks.find({proposer:address, blockHeight: -1}).fetch();
            missedRecords = MissedBlocks.find({proposer:address, blockHeight: {'$gt': 0}}, {sort:{blockHeight:-1}}).fetch();
        }

        if (Meteor.isServer){
            loading = false;
            validatorExist = !!validator;
            statusExist = !!status;
            // missedBlocksExist = !!missedBlocks;
            missedBlocksExist = !!missedRecords;
        }
        else{
            validatorExist = !loading && !!validator;
            statusExist = !loading && !!status;
            // missedBlocksExist = !loading && !!missedBlocks;
            missedBlocksExist = !loading && !!missedRecords;
        }
    }

    return {
        loading,
        validatorExist,
        statusExist,
        missedBlocksExist,
        validator: validatorExist ? validator : {},
        status: statusExist ? status : {},
        // missedBlocks: missedBlocksExist ? missedBlocks : {},
        missedRecords: missedBlocksExist ? missedRecords : [],
        missedRecordsStats: missedBlocksExist ? missedRecordsStats : [],
    };
})(MissedBlocksComponent);
