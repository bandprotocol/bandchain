import React, { Component } from 'react';
import { Link, Switch, Route } from 'react-router-dom';
import numbro from 'numbro';
import moment from 'moment';
import { Markdown } from 'react-showdown';
import Block from '../components/Block.jsx';
import Avatar from '../components/Avatar.jsx';
import PowerHistory from '../components/PowerHistory.jsx';
import { Badge, Row, Col, Card,
    CardBody, Spinner, Nav, NavItem, NavLink } from 'reactstrap';
import KeybaseCheck from '../components/KeybaseCheck.jsx';
import ValidatorDelegations from './Delegations.jsx';
import ValidatorTransactions from '../components/TransactionsContainer.js';
import { DelegationButtons } from '../ledger/LedgerActions.jsx';
import { Helmet } from 'react-helmet';
import LinkIcon from '../components/LinkIcon.jsx';
import i18n from 'meteor/universe:i18n';
import TimeStamp from '../components/TimeStamp.jsx';
import SentryBoundary from '../components/SentryBoundary.jsx';

const T = i18n.createComponent();

addhttp = (url) => {
    if (!/^(f|ht)tps?:\/\//i.test(url)) {
        url = "http://" + url;
    }
    return url;
}

const StatusBadge = (props) =>{
    const statusColor = ['secondary', 'warning', 'success'];
    const statusText = ['Unbonded', 'Unbonding', 'Active'];
    return <h3>
        {props.jailed?<Badge color='danger'><T>validators.jailed</T></Badge>:''}
        <Badge color={statusColor[props.bondingStatus]}>{statusText[props.bondingStatus]}</Badge>
    </h3>;
}

export default class Validator extends Component{
    constructor(props){
        let showdown  = require('showdown');
        showdown.setFlavor('github');
        super(props);
        this.state = {
            identity: "",
            records: "",
            history: "",
            updateTime: "",
            user: localStorage.getItem(CURRENTUSERADDR)
        }
        this.getUserDelegations();
    }

    getUserDelegations() {
        if (this.state.user && this.props.validator && this.props.validator.address) {
            Meteor.call('accounts.getDelegation', this.state.user, this.props.validator.operator_address, (err, res) => {
                if (res && res.shares > 0) {
                    res.tokenPerShare = this.props.validator.tokens/this.props.validator.delegator_shares
                    this.setState({
                        currentUserDelegation: res
                    })
                } else {
                    this.setState({
                        currentUserDelegation: null
                    })
                }

            })
        } else if (this.state.currentUserDelegation != null) {
            this.setState({currentUserDelegation: null})
        }
    }

    static getDerivedStateFromProps(props, state) {
        if (state.user !== localStorage.getItem(CURRENTUSERADDR)) {
            return {user: localStorage.getItem(CURRENTUSERADDR)};
        }
        return null;
    }

    isSameValidator(prevProps) {
        if (this.props.validator == prevProps.validator)
            return true
        if (this.props.validator == null || prevProps.validator == null)
            return false
        return this.props.validator.address === prevProps.validator.address;
    }

    componentDidUpdate(prevProps, prevState){
        if (!this.isSameValidator(prevProps) || this.state.user !== prevState.user)
            this.getUserDelegations();
        if (this.props.validator != prevProps.validator){
            // if (this.props.validator.description.identity != prevProps.validator.description.identity){
            if ((this.props.validator.description) && (this.props.validator.description != prevProps.validator.description)){
                // console.log(prevProps.validator.description);
                if (this.state.identity != this.props.validator.description.identity){
                    this.setState({identity:this.props.validator.description.identity});
                }
            }

            if (this.props.validator.commission){
                if (this.props.validator.commission.update_time == Meteor.settings.public.genesisTime){
                    this.setState({
                        updateTime: "Never changed"
                    });
                }
                else{
                    Meteor.call('Validators.findCreateValidatorTime', this.props.validator.delegator_address, (error, result) => {
                        if (error){
                            console.warn(error);
                        }
                        else{
                            if (result){
                                if (result == this.props.validator.commission.update_time){
                                    this.setState({
                                        updateTime: "Never changed"
                                    });
                                }
                                else{
                                    this.setState({
                                        updateTime: "Updated "+moment(this.props.validator.commission.update_time).fromNow()
                                    });
                                }
                            }
                            else{
                                this.setState({
                                    updateTime: "Updated "+moment(this.props.validator.commission.update_time).fromNow()
                                });
                            }
                        }
                    });
                }
            }
        }

        if (this.props.validatorExist && this.props.validator != prevProps.validator){
            let powerHistory = this.props.validator.history()
            if (powerHistory.length > 0){
                this.setState({
                    history: powerHistory.map((history, i) => {
                        return <PowerHistory
                            key={i}
                            type={history.type}
                            prevVotingPower={history.prev_voting_power}
                            votingPower={history.voting_power}
                            time={history.block_time}
                            height={history.height}
                            address={this.props.validator.operator_address}
                        />
                    })
                })
            }
        }

        if (this.props.records != prevProps.records){
            if (this.props.records.length > 0){
                this.setState({
                    records: this.props.records.map((record, i) => {
                        return <Block key={i} exists={record.exists} height={record.height} />
                    })
                })
            }
        }
    }

    renderShareLink() {
        let validator = this.props.validator;
        let primaryLink = `/validator/${validator.operator_address}`
        let otherLinks = [
            {label: 'Delegate', url: `${primaryLink}/delegate`},
            {label: 'Transfer', url: `/account/${validator.delegator_address}/send`}
        ]

        return <LinkIcon link={primaryLink} otherLinks={otherLinks} />
    }

    render() {
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            if (this.props.validatorExist){
                let moniker = (this.props.validator.description&&this.props.validator.description.moniker)?this.props.validator.description.moniker:this.props.validator.address;
                let identity = (this.props.validator.description&&this.props.validator.description.identity)?this.props.validator.description.identity:"";
                let website = (this.props.validator.description&&this.props.validator.description.website)?this.props.validator.description.website:undefined;
                let details = (this.props.validator.description&&this.props.validator.description.details)?this.props.validator.description.details:"";

                return <Row className="validator-details">
                    <Helmet>
                        <title>{ moniker } - Cosmos Validator | The Big Dipper</title>
                    <meta name="description" content={details} />
                  </Helmet>
                    <Col xs={12}>
                        <Link to="/validators" className="btn btn-link"><i className="fas fa-caret-left"></i> <T>common.backToList</T></Link>
                    </Col>
                    <Col md={4}>
                    <Card body className="text-center">
                        <div className="shareLink d-flex align-self-end">{this.renderShareLink()}</div>
                        <div className="validator-avatar"><Avatar moniker={moniker} profileUrl={this.props.validator.profile_url} address={this.props.validator.address} list={false}/></div>
                        <div className="moniker text-primary">{website?<a href={addhttp(this.props.validator.description.website)} target="_blank">{moniker} <i className="fas fa-link"></i></a>:moniker}</div>
                        <div className="identity"><KeybaseCheck identity={identity} showKey /></div>
                        <div className="details"><Markdown markup={ details } /></div>
                        <div className="website"></div>
                      </Card>
                    <Card>
                          <div className="card-header"><T>validators.uptime</T> <Link className="float-right" to={"/validator/"+this.props.validator.address+"/missed/blocks"}><T>common.more</T>...</Link></div>
                          <SentryBoundary>
                          <CardBody>
                                <Row>
                                    <Col xs={8} className="label"><T numBlocks={Meteor.settings.public.uptimeWindow}>validators.lastNumBlocks</T></Col>
                            <Col xs={4} className="value text-right">{this.props.validator.uptime}%</Col>
                                    <Col md={12} className="blocks-list">{this.state.records}</Col>
                                </Row>
                            </CardBody>
                            </SentryBoundary>
                        </Card>
                  </Col>
                    <Col md={8}>
                        <Card>
                        <div className="card-header"><T>validators.validatorInfo</T></div>
                        <CardBody>
                            <Row>
                                <Col xs={12}><StatusBadge bondingStatus={this.props.validator.status} jailed={this.props.validator.jailed} /></Col>
                                <Col sm={4} className="label"><T>validators.operatorAddress</T></Col>
                                <Col sm={8} className="value address" data-operator-address={this.props.validator.operator_address}>{this.props.validator.operator_address}</Col>
                                <Col sm={4} className="label"><T>validators.selfDelegationAddress</T></Col>
                                <Col sm={8} className="value address" data-delegator-address={this.props.validator.delegator_address}><Link to={"/account/"+this.props.validator.delegator_address}>{this.props.validator.delegator_address}</Link></Col>
                                <Col sm={4} className="label"><T>validators.commissionRate</T></Col>
                                <Col sm={8} className="value">{this.props.validator.commission&&this.props.validator.commission.commission_rates?numbro(this.props.validator.commission.commission_rates.rate*100).format('0.00')+"%":''} <small className="text-secondary">({this.state.updateTime})</small></Col>
                                <Col sm={4} className="label"><T>validators.maxRate</T></Col>
                                <Col sm={8} className="value">{this.props.validator.commission&&this.props.validator.commission.commission_rates?numbro(this.props.validator.commission.commission_rates.max_rate*100).format('0.00')+"%":''}</Col>
                                <Col sm={4} className="label"><T>validators.maxChangeRate</T></Col>
                                <Col sm={8} className="value">{this.props.validator.commission&&this.props.validator.commission.commission_rates?numbro(this.props.validator.commission.commission_rates.max_change_rate*100).format('0.00')+"%":''}</Col>
                            </Row>
                          </CardBody>
                      </Card>
                        <Card>
                        <div className="card-header"><T>common.votingPower</T></div>
                            <CardBody className="voting-power-card">
                            {this.state.user?<DelegationButtons validator={this.props.validator}
                                currentDelegation={this.state.currentUserDelegation}
                                history={this.props.history} stakingParams={this.props.chainStatus.staking.params}/>:''}
                            <Row>
                                {this.props.validator.voting_power?<Col xs={12}><h1 className="display-4 voting-power"><Badge color="primary" >{numbro(this.props.validator.voting_power).format('0,0')}</Badge></h1><span>(~{numbro(this.props.validator.voting_power/this.props.chainStatus.activeVotingPower).format('0.00%')})</span></Col>:''}
                                <Col sm={4} className="label"><T>validators.selfDelegationRatio</T></Col>
                                    <Col sm={8} className="value">{this.props.validator.self_delegation?<span>{numbro(this.props.validator.self_delegation).format("0,0.00%")} <small className="text-secondary">(~{numbro(this.props.validator.voting_power*this.props.validator.self_delegation).format({thousandSeparated: true,mantissa:0})} {Meteor.settings.public.stakingDenom})</small></span>:'N/A'}</Col>
                                    <Col sm={4} className="label"><T>validators.proposerPriority</T></Col>
                                <Col sm={8} className="value">{this.props.validator.proposer_priority?numbro(this.props.validator.proposer_priority).format('0,0'):'N/A'}</Col>
                                <Col sm={4} className="label"><T>validators.delegatorShares</T></Col>
                                    <Col sm={8} className="value">{numbro(this.props.validator.delegator_shares).format('0,0.00')}</Col>
                                    {(this.state.currentUserDelegation)?<Col sm={4} className="label"><T>validators.userDelegateShares</T></Col>:''}
                                    {(this.state.currentUserDelegation)?<Col sm={8} className="value">{numbro(this.state.currentUserDelegation.shares).format('0,0.00')}</Col>:''}
                                <Col sm={4} className="label"><T>validators.tokens</T></Col>
                                    <Col sm={8} className="value">{numbro(this.props.validator.tokens).format('0,0.00')}</Col>
                                    {(this.props.validator.jailed)?<Col xs={12} >
                                <Row><Col md={4} className="label"><T>validators.unbondingHeight</T></Col>
                                            <Col md={8} className="value">{numbro(this.props.validator.unbonding_height).format('0,0')}</Col>
                                            <Col md={4} className="label"><T>validators.unbondingTime</T></Col>
                                            <Col md={8} className="value"><TimeStamp time={this.props.validator.unbonding_time}/></Col>
                                        </Row></Col>:''}
                                </Row>
                          </CardBody>
                      </Card>
                    <Nav pills>
                            <NavItem>
                        <NavLink tag={Link} to={"/validator/"+this.props.validator.operator_address} active={!(this.props.location.pathname.match(/(delegations|transactions)/gm))}><T>validators.powerChange</T></NavLink>
                      </NavItem>
                            <NavItem>
                                <NavLink tag={Link} to={"/validator/"+this.props.validator.operator_address+"/delegations"} active={(this.props.location.pathname.match(/delegations/gm) && this.props.location.pathname.match(/delegations/gm).length > 0)}><T>validators.delegations</T></NavLink>
                      </NavItem>
                            <NavItem>
                                <NavLink tag={Link} to={"/validator/"+this.props.validator.operator_address+"/transactions"} active={(this.props.location.pathname.match(/transactions/gm) && this.props.location.pathname.match(/transactions/gm).length > 0)}><T>validators.transactions</T></NavLink>
                      </NavItem>
                        </Nav>
                        <Switch>
                            <Route exact path="/(validator|validators)/:address" render={() => <div className="power-history">{this.state.history}</div> } />
                            <Route path="/(validator|validators)/:address/delegations" render={() => <ValidatorDelegations address={this.props.validator.operator_address} tokens={this.props.validator.tokens} shares={this.props.validator.delegator_shares} />} />
                            <Route path="/(validator|validators)/:address/transactions" render={() => <ValidatorTransactions validator={this.props.validator.operator_address} delegator={this.props.validator.delegator_address} limit={100}/>} />
                      </Switch>

                        <Link to="/validators" className="btn btn-link"><i className="fas fa-caret-left"></i> <T>common.backToList</T></Link>
                  </Col>
                </Row>
            }
            else{
                return <div><T>validators.validatorNotExists</T></div>
            }
        }
    }

}
