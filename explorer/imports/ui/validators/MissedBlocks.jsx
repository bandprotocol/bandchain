import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Nav, NavItem, NavLink, Spinner, Card, CardDeck } from 'reactstrap';
import { Meteor } from 'meteor/meteor';
import { Helmet } from 'react-helmet';
import i18n from 'meteor/universe:i18n';
import MissedBlocksTable from './MissedBlocksTable.jsx';
import TimeDistubtionChart from './TimeDistubtionChart.jsx';
import TimeStamp from '../components/TimeStamp.jsx';

const T = i18n.createComponent();
export default class MissedBlocks extends Component{
    isVoter() {
        return this.props.match.path.indexOf("/missed/blocks")>0;
    }

    render() {
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            if (this.props.validatorExist){
                return <div>
                    <Helmet>
                        <title>{ this.props.validator.description.moniker } - Missed Blocks | The Big Dipper</title>
                        <meta name="description" content={"The missed blocks and precommits of "+this.props.validator.description.moniker} />
                    </Helmet>
                    <Link to={"/validator/"+this.props.validator.address} className="btn btn-link"><i className="fas fa-caret-left"></i> <T>validators.backToValidator</T></Link>
                    <h2><T moniker={this.props.validator.description.moniker}>validators.missedBlocksTitle</T></h2>
                    <Nav pills>
                        <NavItem>
                            <NavLink tag={Link} to={"/validator/"+this.props.validator.address+"/missed/blocks"} active={this.isVoter()}><T>validators.missedBlocks</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to={"/validator/"+this.props.validator.address+"/missed/precommits"} active={!this.isVoter()}><T>validators.missedPrecommits</T></NavLink>
                        </NavItem>
                    </Nav>
                    {(this.props.missedRecords&&this.props.missedRecords.length>0)?
                        <div className="mt-3">
                            <p className="lead"><T>validators.totalMissed</T> {this.isVoter()?<T>common.blocks</T>:<T>common.precommits</T>}: {this.props.missedRecords.length}</p>
                            <CardDeck>
                                <TimeDistubtionChart missedRecords={this.props.missedRecords} type={this.isVoter()?'blocks':'precommits'}/>
                            </CardDeck>
                            <MissedBlocksTable missedStats={this.props.missedRecordsStats} missedRecords={this.props.missedRecords} type={this.isVoter()?'proposer':'voter'}/>

                        </div>:<div><T>validators.iDontMiss</T>{this.isVoter()?<T>common.block</T>:<T>common.precommit</T>}.</div>}
                    {this.props.statusExist?<div><em><T>validators.lastSyncTime</T>:<TimeStamp time={this.props.status.lastMissedBlockTime}/></em></div>:''}
                </div>
            }
            else return <div><T>validators.validatorNotExists</T></div>
        }

    }
}