import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Transactions } from '/imports/api/transactions/transactions.js';
import List from './List.jsx';

export default ValidatorDetailsContainer = withTracker((props) => {
    let transactionsHandle, transactions, transactionsExist;
    let loading = true;

    if (Meteor.isClient){
        transactionsHandle = Meteor.subscribe('transactions.list', props.limit);
        loading = (!transactionsHandle.ready() && props.limit == Meteor.settings.public.initialPageSize);
    }

    if (Meteor.isServer || !loading){
        transactions = Transactions.find({}, {sort:{height:-1}}).fetch();

        if (Meteor.isServer){
            // loading = false;
            transactionsExist = !!transactions;
        }
        else{
            transactionsExist = !loading && !!transactions;
        }
    }
    
    return {
        loading,
        transactionsExist,
        transactions: transactionsExist ? transactions : {},
    };
})(List);