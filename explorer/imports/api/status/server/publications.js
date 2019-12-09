import { Meteor } from 'meteor/meteor';
import { Status } from '../status.js';
import { check } from 'meteor/check'

Meteor.publish('status.status', function () {
    return Status.find({chainId:Meteor.settings.public.chainId});
});

