import React, { Component } from 'react';
import { Spinner, UncontrolledTooltip, Row, Col, Card, CardHeader, CardBody, Progress } from 'reactstrap';
import numbro from 'numbro';
import AccountCopy from '../components/AccountCopy.jsx';
import LinkIcon from '../components/LinkIcon.jsx';
import Delegations from './Delegations.jsx';
import Unbondings from './Unbondings.jsx';
import AccountTransactions from '../components/TransactionsContainer.js';
import ChainStates from '../components/ChainStatesContainer.js'
import { Helmet } from 'react-helmet';
import { WithdrawButton, TransferButton } from '../ledger/LedgerActions.jsx';
import i18n from 'meteor/universe:i18n';
import Coin from '/both/utils/coins.js'

const T = i18n.createComponent();
export default class AccountDetails extends Component{
    constructor(props){
        super(props);
        this.state = {
            address: props.match.params.address,
            loading: true,
            accountExists: false,
            available: 0,
            delegated: 0,
            unbonding: 0,
            rewards: 0,
            total: 0,
            price: 0,
            user: localStorage.getItem(CURRENTUSERADDR)
        }
    }

    static getDerivedStateFromProps(props, state) {
        if (state.user !== localStorage.getItem(CURRENTUSERADDR)) {
            return {user: localStorage.getItem(CURRENTUSERADDR)};
        }
        return null;
    }

    getBalance(){
        Meteor.call('coinStats.getStats', (error, result) => {
            if (result){
                this.setState({
                    price: result.usd
                })
            }
        });
        Meteor.call('accounts.getBalance', this.props.match.params.address, (error, result) => {
            if (error){
                console.warn(error);
                this.setState({
                    loading:false
                })
            }

            if (result){
                // console.log(result);
                if (result.available){
                    const amount = result.available.amount || '0';
                    this.setState({
                        available: parseFloat(amount),
                        total: parseFloat(this.state.total)+parseFloat(amount)
                    })
                }

                this.setState({delegations: result.delegations || []})
                if (result.delegations && result.delegations.length > 0){
                    result.delegations.forEach((delegation, i) => {
                        const amount = delegation.balance.amount || delegation.balance;
                        this.setState({
                            delegated: this.state.delegated+parseFloat(amount),
                            total: this.state.total+parseFloat(amount)
                        })
                    }, this)
                }

                this.setState({unbondingDelegations: result.unbonding || []})
                if (result.unbonding && result.unbonding.length > 0){
                    result.unbonding.forEach((unbond, i) => {
                        unbond.entries.forEach((entry, j) => {
                            this.setState({
                                unbonding: this.state.unbonding+parseFloat(entry.balance),
                                total: this.state.total+parseFloat(entry.balance)
                            })
                            , this})
                    }, this)
                }

                if (result.rewards && result.rewards.length > 0){
                    result.rewards.forEach((reward, i) => {
                        this.setState({
                            rewards: this.state.rewards+parseFloat(reward.amount),
                            total: this.state.total+parseFloat(reward.amount)
                        })
                    }, this)
                }

                if (result.commission){
                    this.setState({
                        operator_address: result.operator_address,
                        commission: parseFloat(result.commission.amount),
                        total: parseFloat(this.state.total)+parseFloat(result.commission.amount)
                    })
                }


                this.setState({
                    loading:false,
                    accountExists: true
                })
            }
        })
    }

    componentDidMount(){
        this.getBalance();
    }

    componentDidUpdate(prevProps){
        if (this.props.match.params.address !== prevProps.match.params.address){
            this.setState({
                address: this.props.match.params.address,
                loading: true,
                accountExists: false,
                available: 0,
                delegated: 0,
                unbonding: 0,
                commission: 0,
                rewards: 0,
                total: 0,
                price: 0
            }, () => {
                this.getBalance();
            })
        }
    }

    renderShareLink() {
        let primaryLink = `/account/${this.state.address}`
        let otherLinks = [
            {label: 'Transfer', url: `${primaryLink}/send`}
        ]
        return <LinkIcon link={primaryLink} otherLinks={otherLinks} />
    }

    render(){
        if (this.state.loading){
            return <div id="account">
                <h1 className="d-none d-lg-block"><T>accounts.accountDetails</T></h1>
                <Spinner type="grow" color="primary" />
            </div>
        }
        else if (this.state.accountExists){
            return <div id="account">
                <Helmet>
                    <title>Account Details of {this.state.address} on Cosmos Hub | The Big Dipper</title>
                    <meta name="description" content={"Account Details of "+this.state.address+" on Cosmos Hub"} />
                </Helmet>
                <Row>
                    <Col md={3} xs={12}><h1 className="d-none d-lg-block"><T>accounts.accountDetails</T></h1></Col>
                    <Col md={9} xs={12} className="text-md-right"><ChainStates /></Col>
                </Row>
                <Row>
                    <Col><h3 className="text-primary"><AccountCopy address={this.state.address} /></h3></Col>
                </Row>
                <Row>
                    <Col><Card>
                        <CardHeader>
                            Balance
                            <div className="shareLink float-right">{this.renderShareLink()}</div>
                        </CardHeader>
                        <CardBody>
                            <Row className="account-distributions">
                                <Col xs={12}>
                                    <Progress multi>
                                        <Progress bar className="available" value={this.state.available/this.state.total*100} />
                                        <Progress bar className="delegated" value={this.state.delegated/this.state.total*100} />
                                        <Progress bar className="unbonding" value={this.state.unbonding/this.state.total*100} />
                                        <Progress bar className="rewards" value={this.state.rewards/this.state.total*100} />
                                        <Progress bar className="commission" value={this.state.commission/this.state.total*100} />
                                    </Progress>
                                </Col>
                            </Row>
                            <Row>
                                <Col md={6} lg={8}>
                                    <Row>
                                        <Col xs={4} className="label text-nowrap"><div className="available infinity" /><T>accounts.available</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.available).toString(4)}</Col>
                                    </Row>
                                    <Row>
                                        <Col xs={4} className="label text-nowrap"><div className="delegated infinity" /><T>accounts.delegated</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.delegated).toString(4)}</Col>
                                    </Row>
                                    <Row>
                                        <Col xs={4} className="label text-nowrap"><div className="unbonding infinity" /><T>accounts.unbonding</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.unbonding).toString(4)}</Col>
                                    </Row>
                                    <Row>
                                        <Col xs={4} className="label text-nowrap"><div className="rewards infinity" /><T>accounts.rewards</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.rewards).toString(4)}</Col>
                                    </Row>
                                    {this.state.commission?<Row>
                                        <Col xs={4} className="label text-nowrap"><div className="commission infinity" /><T>validators.commission</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.commission).toString(4)}</Col>
                                    </Row>:null}
                                </Col>
                                <Col md={6} lg={4} className="total d-flex flex-column justify-content-end">
                                    {this.state.user?<Row>
                                        <Col xs={12}><TransferButton history={this.props.history} address={this.state.address}/></Col>
                                        {this.state.user===this.state.address?<Col xs={12}><WithdrawButton  history={this.props.history} rewards={this.state.rewards} commission={this.state.commission} address={this.state.operator_address}/></Col>:null}
                                    </Row>:null}
                                    <Row>
                                        <Col xs={4} className="label d-flex align-self-end"><div className="infinity" /><T>accounts.total</T></Col>
                                        <Col xs={8} className="value text-right">{new Coin(this.state.total).toString(4)}</Col>
                                        <Col xs={12} className="dollar-value text-right text-secondary">~{numbro(this.state.total/Meteor.settings.public.stakingFraction*this.state.price).format("$0,0.0000a")} ({numbro(this.state.price).format("$0,0.00")}/{Meteor.settings.public.stakingDenom})</Col>
                                    </Row>
                                </Col>
                            </Row>
                        </CardBody>
                    </Card></Col>
                </Row>
                <Row>
                    <Col md={6}>
                        <Delegations address={this.state.address} delegations={this.state.delegations}/>
                    </Col>
                    <Col md={6}>
                        <Unbondings address={this.state.address} unbonding={this.state.unbondingDelegations}/>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        <AccountTransactions delegator={this.state.address} limit={100}/>
                    </Col>
                </Row>
            </div>
        }
        else{
            return <div id="account">
                <h1 className="d-none d-lg-block"><T>accounts.accountDetails</T></h1>
                <p><T>acounts.notFound</T></p>
            </div>
        }
    }
}
