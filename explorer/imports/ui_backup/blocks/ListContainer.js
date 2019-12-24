import { Meteor } from 'meteor/meteor';
// import React, { Component } from 'react';
import { withTracker } from 'meteor/react-meteor-data';
import { Blockscon } from '/imports/api/blocks/blocks.js';

import Blocks from './List.jsx';

export default BlocksContainer = withTracker((props) => {
    let heightHandle;
    let loading = true;

    if (Meteor.isClient){
        heightHandle = Meteor.subscribe('blocks.height', props.limit);
        loading = (!heightHandle.ready() && props.limit == Meteor.settings.public.initialPageSize);
    }

    let blocks;
    let blocksExist;

    if (Meteor.isServer || !loading){
        blocks = Blockscon.find({}, {sort: {height:-1}}).fetch();
        
        if (Meteor.isServer){
            // loading = false;
            blocksExist = !!blocks;
        }
        else{
            blocksExist = !loading && !!blocks;
        }
    }

    return {
        loading,
        blocksExist,
        blocks: blocksExist ? blocks : {}
    };
})(Blocks);
