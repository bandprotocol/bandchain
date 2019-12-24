import React, { Component } from 'react';
import {Line} from 'react-chartjs-2';
import { Row, Col, Card, CardImg, CardText, CardBody,
    CardTitle, CardSubtitle, Button, Progress, Spinner } from 'reactstrap';
import moment from 'moment';
import i18n from 'meteor/universe:i18n';
import TimeStamp from '../components/TimeStamp.jsx';
import SentryBoundary from '../components/SentryBoundary.jsx';


const T = i18n.createComponent();
export default class Chart extends Component{
    constructor(props){
        super(props);
        this.state = {
            vpData: {},
            timeData: {},
            optionsTime: {},
            optionsVP: {}
        }
    }

    componentDidUpdate(prevProps){
        if (prevProps.history != this.props.history){
            let dates = [];
            let heights = [];
            let blockTime = [];
            let timeDiff = [];
            let votingPower = [];
            let validators = [];
            for (let i in this.props.history){
                dates.push(moment.utc(this.props.history[i].time).format("D MMM YYYY, h:mm:ssa z"));
                heights.push(this.props.history[i].height);
                blockTime.push((this.props.history[i].averageBlockTime/1000).toFixed(2));
                timeDiff.push((this.props.history[i].timeDiff/1000).toFixed(2));
                votingPower.push(this.props.history[i].voting_power);
                validators.push(this.props.history[i].precommits);
            }
            this.setState({
                vpData:{
                    labels:dates,
                    datasets: [
                        {
                            label: 'Voting Power',
                            fill: false,
                            lineTension: 0,
                            yAxisID: 'VotingPower',
                            pointRadius: 1,
                            borderColor: 'rgba(255,152,0,0.5)',
                            borderJoinStyle: 'round',
                            backgroundColor: 'rgba(255,193,101,0.5)',
                            data: votingPower
                        },
                        {
                            label: 'No. of Validators',
                            fill: false,
                            lineTension: 0,
                            yAxisID: 'Validators',
                            pointRadius: 1,
                            borderColor: 'rgba(189,28,8,0.5)',
                            borderJoinStyle: 'round',
                            backgroundColor: 'rgba(255,103,109,0.5)',
                            data: validators,
                        }
                    ]
                },
                timeData:{
                    labels:dates,
                    datasets: [
                        {
                            label: 'Average Block Time',
                            fill: false,
                            lineTension: 0,
                            yAxisID: 'Time',
                            pointRadius: 1,
                            borderColor: 'rgba(156,39,176,0.5)',
                            borderJoinStyle: 'round',
                            backgroundColor: 'rgba(229,112,249,0.5)',
                            data: blockTime,
                            tooltips: {
                                callbacks: {
                                    label: function(tooltipItem, data) {
                                        var label = data.datasets[tooltipItem.datasetIndex].label || '';

                                        if (label) {
                                            label += ': ';
                                        }
                                        label += tooltipItem.yLabel+'s';
                                        return label;
                                    }
                                }
                            }
                        },
                        {
                            label: 'Block Interveral',
                            fill: false,
                            lineTension: 0,
                            yAxisID: 'Time',
                            pointRadius: 1,
                            borderColor: 'rgba(189,28,8,0.5)',
                            borderJoinStyle: 'round',
                            backgroundColor: 'rgba(255,103,109,0.5)',
                            data: timeDiff,
                            tooltips: {
                                callbacks: {
                                    label: function(tooltipItem, data) {
                                        var label = data.datasets[tooltipItem.datasetIndex].label || '';

                                        if (label) {
                                            label += ': ';
                                        }
                                        label += tooltipItem.yLabel+'s';
                                        return label;
                                    }
                                }
                            }
                        },
                        {
                            label: 'No. of Validators',
                            fill: false,
                            lineTension: 0,
                            yAxisID: 'Validators',
                            pointRadius: 1,
                            borderColor: 'rgba(255,152,0,0.5)',
                            borderJoinStyle: 'round',
                            backgroundColor: 'rgba(255,193,101,0.5)',
                            data: validators
                        }
                    ]
                },
                optionsVP: {
                    scales: {
                        xAxes: [
                            {
                                display: false,
                            }
                        ],
                        yAxes: [{
                            id: 'VotingPower',
                            type: 'linear',
                            position: 'left',
                            ticks: {
                                stepSize: 1
                            }
                        }, {
                            id: 'Validators',
                            type: 'linear',
                            position: 'right',
                            ticks: {
                                stepSize: 1
                            }
                        }]
                    }
                },
                optionsTime: {
                    scales: {
                        xAxes: [
                            {
                                display: false,
                            }
                        ],
                        yAxes: [{
                            id: 'Validators',
                            type: 'linear',
                            position: 'right',
                            ticks: {
                                stepSize: 1
                            }
                        }, {
                            id: 'Time',
                            type: 'linear',
                            position: 'left',
                            ticks: {
                            // Include a dollar sign in the ticks
                                callback: function(value, index, values) {
                                    return value+'s';
                                }
                            }
                        }]
                    }
                }
            })
        }
    }

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            if (this.props.historyExist && (this.props.history.length > 0)){
                return (
                    <div>
                        <Card>
                            <div className="card-header"><T>analytics.blockTimeHistory</T></div>
                            <CardBody>
                                <SentryBoundary><Line data={this.state.timeData} options={this.state.optionsTime}/></SentryBoundary>
                            </CardBody>
                        </Card>
                        {/* <Card>
                        <div className="card-header">Voting Power History</div>
                        <CardBody>
                        <Line data={this.state.vpData}  options={this.state.optionsVP}/>
                        </CardBody>
                    </Card> */}
                    </div>
                );
            }
            else{
                return <div></div>
            }
        }
    }
}