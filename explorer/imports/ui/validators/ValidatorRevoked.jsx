import React, { Component } from 'react';
import { Row, Col } from 'reactstrap';
import List from './ListContainer.js';

export default class ValidatorRevoked extends Component{
    constructor(props){
        super(props);
    }

    render() {
        return <div>
            <h1 className="d-none d-lg-block">Jailed Validators</h1>
            <p className="lead">These validators have been jailed. If you know any of them, please ask them to unjail.</p>
            <Row>
                <Col md={12}>
                    <List jailed={true}/>
                </Col>
            </Row>
        </div>
    }
}