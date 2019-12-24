import React, { Component } from 'react';
import { Container } from 'reactstrap';
import { Link } from 'react-router-dom';

export default class NotFound extends Component{
    constructor(props){
        super(props);
    }

    render() {
        return <Container className="not-found">
            <h1>The page you requested cannot be found.</h1>
            <i className="material-icons">block</i>
            <div><Link to="/">Back to home</Link></div>
        </Container>
    }

}