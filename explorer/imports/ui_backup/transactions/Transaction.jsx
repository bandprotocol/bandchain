import React, { Component } from 'react';
import { Container, Row, Col, Card, CardBody, Alert, Spinner } from 'reactstrap';
import { TxIcon } from '../components/Icons.jsx';
import Activities from '../components/Activities.jsx';
import CosmosErrors from '../components/CosmosErrors.jsx';
import { Link } from 'react-router-dom';
import { Markdown } from 'react-showdown';
import numbro from 'numbro';
import { Helmet } from 'react-helmet';
import i18n from 'meteor/universe:i18n';
import Coin from '/both/utils/coins.js';
import TimeStamp from '../components/TimeStamp.jsx';

const T = i18n.createComponent();
export default class Transaction extends Component{
    constructor(props){
        super(props);
        let showdown  = require('showdown');
        showdown.setFlavor('github');
    }

    render(){
        if (this.props.loading){
            return <Container id="transaction">
                <Spinner type="grow" color="primary" />
            </Container>
        }
        else{
            if (this.props.transactionExist){
                // console.log(this.props.transaction);
                let tx = this.props.transaction;
                return <Container id="transaction">
                    <Helmet>
                        <title>Transaction {tx.txhash} on Cosmos Hub | The Big Dipper</title>
                        <meta name="description" content={"Details of transaction "+tx.txhash} />
                    </Helmet>
                    <h4><T>transactions.transaction</T> {(!tx.code)?<TxIcon valid />:<TxIcon />}</h4>
                    {(tx.code)?<Row><Col xs={{size:12, order:"last"}} className="error">
                        <Alert color="danger">
                            <CosmosErrors
                                code={tx.code}
                                logs={tx.logs}
                                gasWanted={tx.gas_wanted}
                                gasUses={tx.gas_used}
                            />
                        </Alert>
                    </Col></Row>:''}
                    <Card>
                        <div className="card-header"><T>common.information</T></div>
                        <CardBody>
                            <Row>
                                <Col md={4} className="label"><T>common.hash</T></Col>
                                <Col md={8} className="value text-nowrap address">{tx.txhash}</Col>
                                <Col md={4} className="label"><T>common.height</T></Col>
                                <Col md={8} className="value">
                                    <Link to={"/blocks/"+tx.height}>{numbro(tx.height).format("0,0")}</Link>
                                    {tx.block()?<span> <TimeStamp time={tx.block().time}/></span>:null}
                                </Col>
                                <Col md={4} className="label"><T>transactions.fee</T></Col>
                                <Col md={8} className="value">{(tx.tx.value.fee.amount.length > 0)?tx.tx.value.fee.amount.map((fee,i) => {
                                    return <span className="text-nowrap" key={i}>{new Coin(fee.amount).toString()}</span>
                                }):<span>No fee</span>}</Col>
                                <Col md={4} className="label"><T>transactions.gasUsedWanted</T></Col>
                                <Col md={8} className="value">{numbro(tx.gas_used).format("0,0")} / {numbro(tx.gas_wanted).format("0,0")}</Col>
                                <Col md={4} className="label"><T>transactions.memo</T></Col>
                                <Col md={8} className="value"><Markdown markup={ tx.tx.value.memo } /></Col>
                            </Row>
                        </CardBody>
                    </Card>
                    <Card>
                        <div className="card-header"><T>transactions.activities</T></div>
                    </Card>
                    {(tx.tx.value.msg && tx.tx.value.msg.length >0)?tx.tx.value.msg.map((msg,i) => {
                        return <Card body key={i}><Activities msg={msg} invalid={(!!tx.code)} events={tx.events} /></Card>
                    }):''}
                </Container>
            }
            else{
                return <Container id="transaction"><div><T>transactions.noTxFound</T></div></Container>
            }
        }
    }
}