import { Mongo } from 'meteor/mongo';
import { Validators } from '../validators/validators.js';

export const Blockscon = new Mongo.Collection('blocks');

Blockscon.helpers({
    proposer(){
        return Validators.findOne({address:this.proposerAddress});
    }
});

// Blockscon.helpers({
//     sorted(limit) {
//         return Blockscon.find({}, {sort: {height:-1}, limit: limit});
//     }
// });


// Meteor.setInterval(function() {
//     Meteor.call('blocksUpdate', (error, result) => {
//         console.log(result);
//     })
// }, 30000000);