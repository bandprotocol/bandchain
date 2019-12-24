import React, { Component } from 'react';
import { Meteor } from 'meteor/meteor';
import { Link } from 'react-router-dom';
import { toast } from 'react-toastify';
import Account from './Account.jsx';

export default class AccountCopy extends Component{
    constructor(props){
        super(props);
    }

    handleCopy = (e) => {
        const el = document.createElement('textarea');
        el.value = e.target.dataset.address;
        el.setAttribute('readonly', '');
        el.style.position = 'absolute';
        el.style.left = '-9999px';
        e.target.appendChild(el);
        el.select();
        document.execCommand('copy');
        e.target.removeChild(el);

        toast("👌🏼 Address copied.");
    }

    render(){
        return <span className="text-nowrap" ><Account address={this.props.address} copy/> <i className="material-icons copy-button" data-address={this.props.address} onClick={this.handleCopy}>file_copy</i></span>
    }
}