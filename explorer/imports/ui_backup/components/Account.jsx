import React, { Component } from 'react';
import { Meteor } from 'meteor/meteor';
import { Link } from 'react-router-dom';
import { Validators } from '/imports/api/validators/validators.js';

const AddressLength = 40;

export default class Account extends Component{
    constructor(props){
        super(props);

        this.state = {
            address: `/account/${this.props.address}`,
            moniker: this.props.address,
            validator: null
        }
    }

    getFields() {
        return {address:1, description:1, operator_address:1, delegator_address:1, profile_url:1};
    }

    getAccount = () => {
        let address = this.props.address;
        let validator = Validators.findOne(
            {$or: [{operator_address:address}, {delegator_address:address}, {address:address}]},
            {fields: this.getFields() });
        if (validator)
            this.setState({
                address: `/validator/${validator.address}`,
                moniker: validator.description?validator.description.moniker:validator.operator_address,
                validator: validator
            });
        else
            this.setState({
                address: `/validator/${address}`,
                moniker: address,
                validator: null
            });
    }

    updateAccount = () => {
        let address = this.props.address;
        Meteor.call('Transactions.findUser', this.props.address, this.getFields(), (error, result) => {
            if (result){
                // console.log(result);
                this.setState({
                    address: `/validator/${result.address}`,
                    moniker: result.description?result.description.moniker:result.operator_address,
                    validator: result
                });
            }
        })
    }

    getAccount = () => {
        let address = this.props.address;
        let validator = Validators.findOne(
            {$or: [{operator_address:address}, {delegator_address:address}, {address:address}]},
            {fields: {address:1, description:1, operator_address:1, delegator_address:1}});
        if (validator)
            this.setState({
                address: `/validator/${validator.address}`,
                moniker: validator.description.moniker
            });
        else
            this.setState({
                address: `/validator/${address}`,
                moniker: address
            });
    }

    componentDidMount(){
        if (this.props.sync)
            this.getAccount();
        else
            this.updateAccount();
    }

    componentDidUpdate(prevProps){
        if (this.props.address != prevProps.address){
            if (this.props.sync) {
                this.getAccount();
            }
            else {
                this.setState({
                    address: `/account/${this.props.address}`,
                    moniker: this.props.address,
                    validator: null
                });
                this.updateAccount();
            }
        }
    }

    userIcon(){
        let signedInAddress = localStorage.getItem(CURRENTUSERADDR);
        if (signedInAddress === this.props.address) {
            return <i className="material-icons account-icon">account_box</i>
        }
    }

    render(){
        return <span className={(this.props.copy)?"address overflow-auto d-inline-block copy":"address overflow-auto d-inline"} >
            <Link to={this.state.address}>{this.userIcon()}{this.state.moniker}</Link>
        </span>
    }
}
