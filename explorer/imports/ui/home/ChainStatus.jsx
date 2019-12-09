import React from 'react';
import { Row, Col, Card, CardText,
    CardTitle, UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem, Spinner } from 'reactstrap';
import numbro from 'numbro';
import i18n from 'meteor/universe:i18n';
import TimeStamp from '../components/TimeStamp.jsx';
import Coin from '/both/utils/coins.js';

const T = i18n.createComponent();

export default class ChainStatus extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            blockHeight: 0,
            blockTime: 0,
            averageBlockTime: 0,
            votingPower: 0,
            numValidators: 0,
            totalNumValidators: 0,
            avgBlockTimeType: "",
            avgVotingPowerType: "",
            blockTimeText: <T>chainStatus.all</T>,
            votingPowerText: <T>chainStatus.now</T>
        }
    }

    componentDidUpdate(prevProps){
        if (prevProps != this.props){
            this.setState({
                blockHeight: numbro(this.props.status.latestBlockHeight).format({thousandSeparated: true}),
                blockTime: <TimeStamp time={this.props.status.latestBlockTime}/>,
                delegatedTokens: numbro(this.props.status.totalVotingPower).format('0,0.00a'),
                numValidators: this.props.status.validators,
                totalNumValidators: this.props.status.totalValidators,
                bondedTokens: this.props.states.bondedTokens,
                totalSupply: this.props.states.totalSupply
            })

            switch (this.state.avgBlockTimeType){
            case "":
                this.setState({
                    averageBlockTime: numbro(this.props.status.blockTime/1000).format('0,0.00')
                })
                break;
            case "m":
                this.setState({
                    averageBlockTime: numbro(this.props.status.lastMinuteBlockTime/1000).format('0,0.00')
                })
                break;
            case "h":
                this.setState({
                    averageBlockTime: numbro(this.props.status.lastHourBlockTime/1000).format('0,0.00')
                })
                break;
            case "d":
                this.setState({
                    averageBlockTime: numbro(this.props.status.lastDayBlockTime/1000).format('0,0.00')
                })
                break;
            }

            switch (this.state.avgVotingPowerType){
            case "":
                this.setState({
                    votingPower: numbro(this.props.status.activeVotingPower).format('0,0.00a'),
                });
                break;
            case "h":
                this.setState({
                    votingPower: numbro(this.props.status.lastHourVotingPower).format('0,0.00a'),
                });
                break;
            case "d":
                this.setState({
                    votingPower: numbro(this.props.status.lastDayVotingPower).format('0,0.00a'),
                });
                break;

            }

        }
    }

    handleSwitchBlockTime = (type,e) => {
        e.preventDefault();
        switch (type){
        case "":
            this.setState({
                blockTimeText: <T>chainStatus.all</T>,
                avgBlockTimeType: "",
                averageBlockTime: numbro(this.props.status.blockTime/1000).format('0,0.00')
            })
            break;
        case "m":
            this.setState({
                blockTimeText: "1m",
                avgBlockTimeType: "m",
                averageBlockTime: numbro(this.props.status.lastMinuteBlockTime/1000).format('0,0.00')
            })
            break;
        case "h":
            this.setState({
                blockTimeText: "1h",
                avgBlockTimeType: "h",
                averageBlockTime: numbro(this.props.status.lastHourBlockTime/1000).format('0,0.00')
            })
            break;
        case "d":
            this.setState({
                blockTimeText: "1d",
                avgBlockTimeType: "d",
                averageBlockTime: numbro(this.props.status.lastDayBlockTime/1000).format('0,0.00')
            })
            break;

        }
    }

    handleSwitchVotingPower = (type,e) => {
        e.preventDefault();
        switch (type){
        case "":
            this.setState({
                votingPowerText: <T>chainStatus.now</T>,
                avgVotingPowerType: "",
                votingPower: numbro(this.props.status.activeVotingPower).format('0,0.00a')
            })
            break;
        case "h":
            this.setState({
                votingPowerText: "1h",
                avgVotingPowerType: "h",
                votingPower: numbro(this.props.status.lastHourVotingPower).format('0,0.00a')
            })
            break;
        case "d":
            this.setState({
                votingPowerText: "1d",
                avgVotingPowerType: "d",
                votingPower: numbro(this.props.status.lastDayVotingPower).format('0,0.00a')
            })
            break;

        }
    }

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else {
            if (this.props.statusExist && this.props.status.prevotes){
                return(
                    <Row className="status text-center">
                        <Col lg={3} md={6}>
                            <Card body>
                                <CardTitle><T>chainStatus.latestHeight</T></CardTitle>
                                <CardText>
                                    <span className="display-4 value text-primary">{this.state.blockHeight}</span>
                                    {this.state.blockTime}
                                </CardText>
                            </Card>
                        </Col>
                        <Col lg={3} md={6}>
                            <Card body>
                                <UncontrolledDropdown size="sm" className="more">
                                    <DropdownToggle>
                                        <i className="material-icons">more_vert</i>
                                    </DropdownToggle>
                                    <DropdownMenu>
                                        <DropdownItem onClick={(e) => this.handleSwitchBlockTime("", e)}><T>chainStatus.allTime</T></DropdownItem>
                                        {this.props.status.lastMinuteBlockTime?<DropdownItem onClick={(e) => this.handleSwitchBlockTime("m", e)}><T>chainStatus.lastMinute</T></DropdownItem>:''}
                                        {this.props.status.lastHourBlockTime?<DropdownItem onClick={(e) => this.handleSwitchBlockTime("h", e)}><T>chainStatus.lastHour</T></DropdownItem>:''}
                                        {this.props.status.lastDayBlockTime?<DropdownItem onClick={(e) => this.handleSwitchBlockTime("d", e)}><T>chainStatus.lastDay</T> </DropdownItem>:''}
                                    </DropdownMenu>
                                </UncontrolledDropdown>
                                <CardTitle><T>chainStatus.averageBlockTime</T> ({this.state.blockTimeText})</CardTitle>
                                <CardText>
                                    <span className="display-4 value text-primary">{this.state.averageBlockTime}</span><T>chainStatus.seconds</T>
                                </CardText>
                            </Card>
                        </Col>
                        <Col lg={3} md={6}>
                            <Card body>
                                <CardTitle><T>chainStatus.activeValidators</T></CardTitle>
                                <CardText><span className="display-4 value text-primary">{this.state.numValidators}</span><T totalValidators={this.state.totalNumValidators}>chainStatus.outOfValidators</T></CardText>
                            </Card>
                        </Col>
                        <Col lg={3} md={6}>
                            <Card body>
                                <UncontrolledDropdown size="sm" className="more">
                                    <DropdownToggle>
                                        <i className="material-icons">more_vert</i>
                                    </DropdownToggle>
                                    <DropdownMenu>
                                        <DropdownItem onClick={(e) => this.handleSwitchVotingPower("", e)}><T>chainStatus.now</T></DropdownItem>
                                        {this.props.status.lastHourVotingPower?<DropdownItem onClick={(e) => this.handleSwitchVotingPower("h", e)}><T>chainStatus.lastHour</T></DropdownItem>:''}
                                        {this.props.status.lastDayVotingPower?<DropdownItem onClick={(e) => this.handleSwitchVotingPower("d", e)}><T>chainStatus.lastDay</T></DropdownItem>:''}
                                    </DropdownMenu>
                                </UncontrolledDropdown>
                                <CardTitle><T>chainStatus.onlineVotingPower</T> ({this.state.votingPowerText})</CardTitle>
                                <CardText><span className="display-4 value text-primary">{this.state.votingPower}</span><T percent={numbro(this.state.bondedTokens/this.state.totalSupply).format("0.00%")} totalStakes={numbro(this.state.totalSupply/Coin.StakingFraction).format("0.00a")} denom={Coin.StakingDenom} denomPlural={Coin.StakingDenomPlural}>chainStatus.fromTotalStakes</T></CardText>
                            </Card>
                        </Col>
                    </Row>
                )
            }
            else{
                return <div></div>
            }
        }
    }
}