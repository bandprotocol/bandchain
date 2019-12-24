import React, { Component } from 'react';
import { Button, Spinner, TabContent, TabPane, Row, Col, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap';
import {Ledger} from './ledger.js';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

class LedgerModal extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            loading: false,
            activeTab: '1'
        };
        this.ledger = new Ledger({testModeAllowed: false});
    }

    autoOpenModal = () => {
        if (!this.props.isOpen && this.props.handleLoginConfirmed) {
            this.tryConnect(5000);
            this.props.toggle(true);
        }
    }

    componentDidMount() {
        this.autoOpenModal();
    }

    componentDidUpdate(prevProps, prevState) {
        this.autoOpenModal();
        if (this.props.isOpen && !prevProps.isOpen) {
            this.tryConnect();
        }
    }

    tryConnect = (timeout=undefined) => {
        if (this.state.loading) return
        this.setState({ loading: true, errorMessage: '' })
        this.ledger.getCosmosAddress(timeout).then((res) => {
            let currentUser = localStorage.getItem(CURRENTUSERADDR);
            if (this.props.handleLoginConfirmed && res.address === currentUser) {
                this.closeModal(true)
            } else {
                this.setState({
                    currentUser: currentUser,
                    address: res.address,
                    pubKey: Buffer.from(res.pubKey).toString('base64'),
                    errorMessage: '',
                    loading: false,
                    activeTab: '2'});
                this.trySignIn();
            }
        }, (err) => {
            this.setState({
                errorMessage: err.message,
                loading: false,
                activeTab: '1'
            })
        });
    }

    trySignIn = () => {
        this.setState({ loading: true, errorMessage: '' })
        this.ledger.confirmLedgerAddress().then((res) => {
            localStorage.setItem(CURRENTUSERADDR, this.state.address);
            localStorage.setItem(CURRENTUSERPUBKEY, this.state.pubKey);
            this.props.refreshApp();
            this.closeModal(true);
        }, (err) => {
            this.setState({
                errorMessage: err.message,
                loading: false
            })})
    }

    getActionButton() {
        if (this.state.activeTab === '1' && !this.state.loading)
            return <Button color="primary"  onClick={this.tryConnect}><T>common.retry</T></Button>
        if (this.state.activeTab === '2' && this.state.errorMessage !== '')
            return <Button color="primary"  onClick={this.trySignIn}><T>common.retry</T></Button>
    }

    closeModal = (success) => {
        if (this.props.handleLoginConfirmed) {
            this.props.handleLoginConfirmed(typeof success ==='boolean'?success:false);
        }
        this.setState({
            loading: false,
            errorMessage: '',
            currentUser: null,
            address: null,
            activeTab: '1'
        })
        this.props.toggle(false)
    }

    render() {
        return (
            <Modal isOpen={this.props.isOpen} toggle={this.closeModal} className="ledger-sign-in">
            <ModalHeader><T>accounts.signInWithLedger</T></ModalHeader>
            <ModalBody>
                <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="1">
                    <T _purify={false}>accounts.signInWarning</T>
                </TabPane>
                <TabPane tabId="2">
                    {this.state.currentUser?<span>You are currently logged in as <strong className="text-primary d-block">{this.state.currentUser}.</strong></span>:null}
                    <T>accounts.toLoginAs</T> <strong className="text-primary d-block">{this.state.address}</strong><T>accounts.pleaseAccept</T>
                </TabPane>
            </TabContent>
            {this.state.loading?<Spinner type="grow" color="primary" />:''}
            <p className="error-message">{this.state.errorMessage}</p>
        </ModalBody>
            <ModalFooter>
            {this.getActionButton()}
                <Button color="secondary" onClick={this.closeModal}><T>common.cancel</T></Button>
            </ModalFooter>
            </Modal>
        );
    }
}

export default LedgerModal;