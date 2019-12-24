import { Meteor } from 'meteor/meteor';
import { withTracker } from 'meteor/react-meteor-data';
import { Transactions } from '/imports/api/transactions/transactions.js';
import Transaction from './Transaction.jsx';

export default TransactionContainer = withTracker((props) => {
    let txId = props.match.params.txId.toUpperCase();

    let transactionsHandle, transaction, transactionExist;
    let loading = false;

    if (Meteor.isClient){
        transactionsHandle = Meteor.subscribe('transactions.findOne', txId);
        loading = !transactionsHandle.ready();
    }

    if (Meteor.isServer || !loading){
        transaction = Transactions.findOne({txhash: txId});

        if (Meteor.isServer){
            loading = false;
            transactionExist = !!transaction;
        }
        else{
            transactionExist = !loading && !!transaction;
        }
        if (props.location.search === '?new' && !transactionExist) {
            loading = true;
        }
    }

    return {
        loading,
        transactionExist,
        transaction: transactionExist ? transaction : {},
    };
})(Transaction);