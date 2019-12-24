import React, { Component } from 'react';
import {HorizontalBar} from 'react-chartjs-2';
import { Row, Col, Card, CardImg, CardText, CardBody,
    CardTitle, CardSubtitle, Button, Progress, Spinner } from 'reactstrap';
import numbro from 'numbro';
import i18n from 'meteor/universe:i18n';
import SentryBoundary from '../components/SentryBoundary.jsx';


const T = i18n.createComponent();

export default class VotingPower extends Component{
    constructor(props){
        super(props);
        this.state = {
            data: {},
            options: {}
        }
    }

    componentDidUpdate(prevProps){
        if (prevProps.stats != this.props.stats){
            let self = this;

            let labels = [];
            let data = [];
            let totalVotingPower = 0;
            let accumulatePower = [];
            let backgroundColors = [];
            
            for (let i in this.props.stats){
                totalVotingPower += this.props.stats[i].voting_power;
                if (i > 0){
                    accumulatePower[i] = accumulatePower[i-1] + this.props.stats[i].voting_power;
                }
                else{
                    accumulatePower[i] = this.props.stats[i].voting_power;
                }
            }

            for (let v in this.props.stats){
                labels.push(this.props.stats[v].description?this.props.stats[v].description.moniker:'');
                data.push(this.props.stats[v].voting_power);
                let alpha = (this.props.stats.length+1-v)/this.props.stats.length*0.8+0.2;
                backgroundColors.push('rgba(189, 8, 28,'+alpha+')');
            }
            this.setState({
                data:{
                    labels:labels,
                    datasets: [
                        {
                            label: "Voting Power",
                            data: data,
                            backgroundColor: backgroundColors
                        }
                    ]
                },
                options:{
                    tooltips: {
                        callbacks: {
                            label: function(tooltipItem, data) {
                                return numbro(data.datasets[0].data[tooltipItem.index]).format("0,0")+" ("+(numbro(data.datasets[0].data[tooltipItem.index]/totalVotingPower).format("0.00%")+", "+numbro(accumulatePower[tooltipItem.index]/totalVotingPower).format("0.00%"))+")";
                            }
                        }
                    },
                    maintainAspectRatio: false,
                    scales: {
                        xAxes: [{
                            ticks: {
                                beginAtZero:true,
                                userCallback: function(value, index, values) {
                                    // Convert the number to a string and splite the string every 3 charaters from the end
                                    return numbro(value).format("0,0");
                                }
                            }
                        }]
                    }
                }
            });

            $("#voting-power-chart").height(16*data.length);
        }
    }

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            if (this.props.statsExist && this.props.stats){
                return (                    
                    <Card>
                        <div className="card-header"><T>common.votingPower</T></div>
                        <CardBody id="voting-power-chart">
                            <SentryBoundary><HorizontalBar data={this.state.data} options={this.state.options} /></SentryBoundary>
                        </CardBody>
                    </Card>
                );   
            }
            else{
                return <div></div>
            }
        }
    }
}    
