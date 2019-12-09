import React, { Component } from 'react';
import { Row, Col } from 'reactstrap';
import List from './ListContainer.js';
import { LoadMore } from '../components/LoadMore.jsx';
import { Meteor } from 'meteor/meteor';
import { Route, Switch } from 'react-router-dom';
import Transaction from './TransactionContainer.js';
import Sidebar from "react-sidebar";
import ChainStates from '../components/ChainStatesContainer.js'
import { Helmet } from 'react-helmet';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

export default class Transactions extends Component{
    constructor(props){
        super(props);

        this.state = {
            limit: Meteor.settings.public.initialPageSize,
            monikerDir: 1,
            votingPowerDir: -1,
            uptimeDir: -1,
            proposerDir: -1,
            priority: 2,
            loadmore: false,
            sidebarOpen: (props.location.pathname.split("/transactions/").length == 2)
        }

        this.onSetSidebarOpen = this.onSetSidebarOpen.bind(this);
    }

    isBottom(el) {
        return el.getBoundingClientRect().bottom <= window.innerHeight;
    }
      
    componentDidMount() {
        document.addEventListener('scroll', this.trackScrolling);
    }
    
    componentWillUnmount() {
        document.removeEventListener('scroll', this.trackScrolling);
    }
    
    componentDidUpdate(prevProps){
        if (this.props.location.pathname != prevProps.location.pathname){
            this.setState({
                sidebarOpen: (this.props.location.pathname.split("/transactions/").length == 2)
            })
        }
    }

    trackScrolling = () => {
        const wrappedElement = document.getElementById('transactions');
        if (this.isBottom(wrappedElement)) {
            // console.log('header bottom reached');
            document.removeEventListener('scroll', this.trackScrolling);
            this.setState({loadmore:true});
            this.setState({
                limit: this.state.limit+10
            }, (err, result) => {
                if (!err){
                    document.addEventListener('scroll', this.trackScrolling);
                }
                if (result){
                    this.setState({loadmore:false});
                }
            })
        }
    };

    onSetSidebarOpen(open) {
        // console.log(open);
        this.setState({ sidebarOpen: open }, (error, result) =>{
            let timer = Meteor.setTimeout(() => {
                if (!open){
                    this.props.history.push('/transactions');
                }
                Meteor.clearTimeout(timer);
            },500)
        });
        
    }

    render(){
        return <div id="transactions">
            <Helmet>
                <title>Latest Transactions on Cosmos Hub | The Big Dipper</title>
                <meta name="description" content="See what is happening on Cosmos Hub" />
            </Helmet>
            <Row>
                <Col md={3} xs={12}><h1 className="d-none d-lg-block"><T>transactions.transactions</T></h1></Col>
                <Col md={9} xs={12} className="text-md-right"><ChainStates /></Col>
            </Row>
            <Switch>
                <Route path="/transactions/:txId" render={(props)=> <Sidebar 
                    sidebar={<Transaction {...props} />}
                    open={this.state.sidebarOpen}
                    onSetOpen={this.onSetSidebarOpen}
                    styles={{ sidebar: { 
                        background: "white", 
                        position: "fixed",
                        width: '85%',
                        zIndex: 4
                    },overlay: {
                        zIndex: 3
                    } }}
                >
                </Sidebar>} />
            </Switch>
            <List limit={this.state.limit} />
            <LoadMore show={this.state.loadmore} />
        </div>
    }
}