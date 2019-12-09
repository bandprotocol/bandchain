import React, { Component } from 'react';
import { Meteor } from 'meteor/meteor'
import moment from 'moment-timezone';
import { UncontrolledTooltip } from 'reactstrap';

export default class TimeStamp extends Component{
    constructor(props){
        super(props);
        this.ref = React.createRef();
    }

    render(){
        let format = "D MMM YYYY, h:mm:ssa z";
        let timezone = moment.tz.guess()
        let time = moment.utc(this.props.time);
        return <span>
            <span ref={this.ref}>{time.format(format)}</span>
            <UncontrolledTooltip placement='bottom' target={this.ref}>
                {time.tz(timezone).format(format)}
            </UncontrolledTooltip>
        </span>

    }
}