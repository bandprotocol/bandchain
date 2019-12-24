import React, { Component } from 'react';
import { Meteor } from 'meteor/meteor';
import { Link } from 'react-router-dom';
import { Badge, Progress, Row, Col, Card, Spinner, UncontrolledTooltip } from 'reactstrap';
import numbro from 'numbro';
import Avatar from '../components/Avatar.jsx';
import Account from './Account.jsx';
import TimeStamp from '../components/TimeStamp.jsx';

export default class AccountTooltip extends Account{
    constructor(props){
        super(props);

        this.ref = React.createRef();
    }

    shouldComponentUpdate(nextProps, nextState) {
        return (
            this.props.address !== nextProps.address ||
            this.state.address !== nextState.address ||
            this.state.moniker !== nextState.moniker)
    }

    getFields() {
        return {
            status: 1,
            description: 1,
            delegator_shares: 1,
            operator_address: 1,
            tokens: 1,
            commission: 1,
            unbonding_time: 1,
            jailed: 1,
            delegator_address: 1,
            address: 1,
            operator_pubkey: 1,
            voting_power: 1,
            lastSeen: 1,
            uptime: 1,
            self_delegation: 1,
            profile_url: 1
        }
    }

    renderDetailTooltip() {
        if (!this.state.validator)
            return
        let validator = this.state.validator;
        let moniker = validator.description && validator.description.moniker || validator.address;
        let isActive = validator.status == 2 && !validator.jailed;

        return <UncontrolledTooltip key='tooltip' className='validator-tooltip' placement='right' flip={false} target={this.ref} autohide={false} fade={false}>
            <Card body className='validator-tooltip-card'>
                <Row className='d-flex justify-content-center'>
                    <h4 className="moniker text-primary">{moniker}</h4>
                </Row>
                <Row className='d-flex justify-content-center avatar'>
                    <Avatar moniker={moniker} profileUrl={validator.profile_url} address={validator.address}/>
                </Row>
                <Row className="voting-power data">
                    <i className="material-icons">power</i>
                    {validator.voting_power?numbro(validator.voting_power).format('0,0'):0}
                </Row>
                <Row className="self-delegation data">
                    <i className="material-icons">equalizer</i>
                    {validator.self_delegation?numbro(validator.self_delegation).format('0.00%'):'N/A'}
                </Row>
                {(isActive)?<Row className="commission data">
                    <i className="material-icons">call_split</i>
                    {numbro(validator.commission.rate).format('0.00%')}
                </Row>:null}
                {(!isActive)?<Row className="last-seen data">
                    {validator.lastSeen?<TimeStamp time={validator.lastSeen}/>:
                     (validator.unbonding_time?<TimeStamp time={validator.unbonding_time}/>:null)}
                </Row>:null}
                {(!isActive)?<Row className="bond-status data" xs={2}>
                    <Col xs={6}>{(validator.status == 0)?<Badge color="secondary">Unbonded</Badge>:<Badge color="warning">Unbonding</Badge>}</Col>
                    <Col xs={6}>{validator.jailed?<Badge color="danger">Jailed</Badge>:''}</Col>
                </Row>:null}
                {(isActive)?<Row className="uptime data">
                    <i className="material-icons">flash_on</i>
                    <Progress value={validator.uptime} style={{width:'80%'}}>
                       {validator.uptime?numbro(validator.uptime/100).format('0%'):0}
                    </Progress>
                </Row>:null}
            </Card>
        </UncontrolledTooltip>
    }

    render(){
        return [
            <span ref={this.ref} key='link'>
                <Link to={this.state.address}>{this.userIcon()}{this.state.moniker}</Link>
            </span>,
            this.renderDetailTooltip()
        ]
    }
}
