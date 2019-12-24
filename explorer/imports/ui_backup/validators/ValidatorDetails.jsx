import React, { Component } from 'react';
import { Row, Col } from 'reactstrap';
import { Route, Switch } from 'react-router-dom';
import Validator from './ValidatorContainer.js';
import MissedBlocks from './MissedBlocksContainer.js';
import ChainStates from '../components/ChainStatesContainer.js'
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

export default class ValidatorDetails extends Component{
    constructor(props){
        super(props);
    }

    render() {
        return <div>
            <Row>
            <Col lg={3} xs={12}><h1 className="d-none d-lg-block"><T>validators.validatorDetails</T></h1></Col>
                <Col lg={9} xs={12} className="text-lg-right"><ChainStates /></Col>
          </Row>
            <Row>
                <Col md={12}>
                    <Switch>
                        <Route exact path="/(validator|validators)/:address/missed/blocks" render={(props) => <MissedBlocks {...props} type='voter' />} />
                        <Route exact path="/(validator|validators)/:address/missed/precommits" render={(props) => <MissedBlocks {...props} type='proposer' />} />
                        <Route path="/(validator|validators)/:address" render={(props) => <Validator address={props.match.params.address} {...props}/>} />
                  </Switch>
              </Col>
          </Row>
        </div>
    }

}