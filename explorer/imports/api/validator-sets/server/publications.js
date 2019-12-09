import { Meteor } from 'meteor/meteor';
import { ValidatorSets } from '../validator-set.js';

Meteor.publish('validatorSets.all', function () {
    return ValidatorSets.find();
});
