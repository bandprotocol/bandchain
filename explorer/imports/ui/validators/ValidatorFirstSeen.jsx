import React, { Component } from 'react';
import { Row, Col } from 'reactstrap';
import FirstSeenListContainer from './FirstSeen.js';

export default class ValidatorFirstSeen extends Component{
    constructor(props){
        super(props);
    }

    render() {
        return <div>
            <h1 className="d-none d-lg-block">Validators</h1>
            <p className="lead">First seen time of validators.</p>
            <Row>
                <Col md={12}>
                    <FirstSeenListContainer />
                    {/* <List jailed={true}/> */}
                </Col>
            </Row>
        </div>
    }
}