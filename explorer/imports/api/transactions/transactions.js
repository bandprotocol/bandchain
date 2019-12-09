import { Mongo } from 'meteor/mongo';
import { Blockscon } from '../blocks/blocks.js';
import { TxIcon } from '../../ui/components/Icons.jsx';

export const Transactions = new Mongo.Collection('transactions');

Transactions.helpers({
    block(){
        return Blockscon.findOne({height:this.height});
    }
})