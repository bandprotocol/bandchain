import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Validators } from '/imports/api/validators/validators.js';
import { ValidatorRecords } from '/imports/api/records/records.js';
import { Chain } from '/imports/api/chain/chain.js';
import Validator from './Validator.jsx';

export default ValidatorDetailsContainer = withTracker((props) => {
    let chainHandle;
    let validatorHandle;
    let validatorsHandle;
    let loading = true;

    if (Meteor.isClient){
        chainHandle = Meteor.subscribe('chain.status');
        validatorsHandle = Meteor.subscribe('validators.all', props.address);
        validatorHandle = Meteor.subscribe('validator.details', props.address);
        loading = !validatorHandle.ready() && !validatorsHandle.ready() && !chainHandle.ready();
    }

    let options = {address:props.address};

    let chainStatus;
    let validatorExist;
    let validator;
    let validatorRecords;

    if (Meteor.isServer || !loading){
        if (props.address.indexOf(Meteor.settings.public.bech32PrefixValAddr) != -1){
            options = {operator_address:props.address}
        }
        validator = Validators.findOne(options);

        if (validator){
            validatorRecords = ValidatorRecords.find({address:validator.address}, {sort:{height:-1}}).fetch();
        }

        chainStatus = Chain.findOne({chainId:Meteor.settings.public.chainId});

        if (Meteor.isServer){
            loading = false;
            validatorExist = !!validator && !!validatorRecords && !!chainStatus;
        }
        else{
            validatorExist = !loading && !!validator && !!validatorRecords && !!chainStatus;
        }

        // loading = false;
    }
    // console.log(props.state.limit);
    return {
        loading,
        validatorExist,
        validator: validatorExist ? validator : {},
        records: validatorExist ? validatorRecords : {},
        chainStatus: validatorExist ? chainStatus : {}
    };
})(Validator);
