import React, { Component } from 'react';
import { Table, Spinner } from 'reactstrap';
import { Link } from 'react-router-dom';
import { DenomSymbol, ProposalStatusIcon } from '../components/Icons.jsx';
import numbro from 'numbro';
import i18n from 'meteor/universe:i18n';
import Coin from '/both/utils/coins.js'
import TimeStamp from '../components/TimeStamp.jsx';
import { SubmitProposalButton } from '../ledger/LedgerActions.jsx';

const T = i18n.createComponent();

const ProposalRow = (props) => {
    return <tr>
        <th className="d-none d-sm-table-cell counter">{props.proposal.proposalId}</th>
        <td className="title"><Link to={"/proposals/"+props.proposal.proposalId}>{props.proposal.content.value.title}</Link></td>
        <td className="status"><ProposalStatusIcon status={props.proposal.proposal_status}/><span className="d-none d-sm-inline">{props.proposal.proposal_status.match(/[A-Z]+[^A-Z]*|[^A-Z]+/g).join(" ")}</span></td>
        <td className="submit-block"><TimeStamp time={props.proposal.submit_time}/></td>
        <td className="voting-start">{(props.proposal.voting_start_time != "0001-01-01T00:00:00Z")?<TimeStamp time={props.proposal.voting_start_time}/>:'Not started'}</td>
        <td className="deposit text-right">{props.proposal.total_deposit?props.proposal.total_deposit.map((deposit, i) => {
            return <div key={i}>{new Coin(deposit.amount).toString()}</div>
        }):'0'}</td>
    </tr>
}

export default class List extends Component{
    constructor(props){
        super(props);
        if (Meteor.isServer){
            if (this.props.proposals.length > 0){
                this.state = {
                    proposals: this.props.proposals.map((proposal, i) => {
                        return <ProposalRow key={i} index={i} proposal={proposal} />
                    })
                }
            }
        }
        else{
            this.state = {
                proposals: null
            }
        }
    }

    static getDerivedStateFromProps(props, state) {
        if (state.user !== localStorage.getItem(CURRENTUSERADDR)) {
            return {user: localStorage.getItem(CURRENTUSERADDR)};
        }
        return null;
    }

    componentDidUpdate(prevState){
        if (this.props.proposals != prevState.proposals){
            if (this.props.proposals.length > 0){
                this.setState({
                    proposals: this.props.proposals.map((proposal, i) => {
                        return <ProposalRow key={i} index={i} proposal={proposal} />
                    })
                })
            }
        }
    }

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            return (
                <div>
                    {this.state.user?<SubmitProposalButton history={this.props.history}/>:null}
                    <Table striped className="proposal-list">
                        <thead>
                            <tr>
                                <th className="d-none d-sm-table-cell counter"><i className="fas fa-hashtag"></i> <T>proposals.proposalID</T></th>
                                <th className="title"><i className="material-icons">view_headline</i> <span className="d-none d-sm-inline"><T>proposals.title</T></span></th>
                                <th className="status"><i className="fas fa-toggle-on"></i> <span className="d-none d-sm-inline"><T>proposals.status</T></span></th>
                                <th className="submit-block"><i className="fas fa-box"></i> <span className="d-none d-sm-inline"><T>proposals.submitTime</T> (UTC)</span></th>
                                <th className="voting-start"><i className="fas fa-box-open"></i> <span className="d-none d-sm-inline"><T>proposals.votingStartTime</T> (UTC)</span></th>
                                <th className="deposit text-right"><i className="material-icons">attach_money</i> <span className="d-none d-sm-inline"><T>proposals.totalDeposit</T></span></th>
                            </tr>
                        </thead>
                        <tbody>{this.state.proposals}</tbody>
                    </Table>
                </div>
            )
        }
    }
}