import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { CardHeader, CardBody, Card } from 'reactstrap';
import moment from 'moment';
import PChart from '../components/Chart.jsx';

const BATCHSIZE = 15;
const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const displayTimeRange = (time) => {
    let startTime = time.format("h:mm");
    let endTime = time.clone().add(BATCHSIZE, 'minute').format("h:mm");
    return `from ${startTime} to ${endTime} on ${time.format("D MMM YYYY")}`;
}

export default class TimeDistubtionChart extends Component{
    populateChartData() {
        let timeline = [];
        let breakdown = [];
        let i;
        for (i = 0;i < (7 * 24); i ++)
            breakdown.push({
                hour: i % 24,
                day: daysOfWeek[Math.floor(i / 24)],
                count: 0
            });
        let prevBatch = null;
        this.props.missedRecords.forEach((record) => {
            let time = moment.utc(record.time)
            breakdown[time.day() * 24 + time.hour()].count += 1;
            if (prevBatch && time.diff(prevBatch) >= 0) {
                timeline[timeline.length - 1].y = timeline[timeline.length - 1].y + 1
            } else {
                prevBatch = time.minutes(Math.floor(time.minute()/BATCHSIZE) * BATCHSIZE).seconds(0);
                timeline.push({
                    x: prevBatch,
                    y: 1
                })
            }
        })


        return {
            timeline,
            breakdown
        }
    }

    populateTimelineChart(timeline) {
        let layout = [['yAxis','barPlot'],[null, 'xAxis']];
        let scales = [{
            scaleId: 'xScale',
            type: 'Time'
        }, {
            scaleId: 'yScale',
            type: 'Linear'
        }];
        let datasets = [{
            datasetId: 'timeline',
            data: timeline
        }];
        let components = {
            plots: [{
                plotId: 'barPlot',
                type: 'Bar',
                x: {
                    value: (d, i, ds) => d.x.toDate(),
                    scale: 'xScale'
                },
                y: {
                    value: (d, i, ds) => d.y,
                    scale: 'yScale'
                },
                datasets: ['timeline'],
                interactions: {
                    PanZoom: {
                        xScales: ['xScale'],
                        yScales: ['yScale']
                    }
                },
                tooltip: (c, p, data, ds) => `missed ${data.y} ${this.props.type} ${displayTimeRange(data.x)}`
            }],
            axes: [{
                axisId: 'xAxis',
                type: 'Time',
                scale: 'xScale',
                orientation: 'bottom',
                interactions: {
                    PanZoom: {
                        xScales: ['xScale']
                    }
                }
            },{
                axisId: 'yAxis',
                type: 'Numeric',
                scale: 'yScale',
                orientation: 'left',
                interactions: {
                    PanZoom: {
                        yScales: ['yScale']
                    }
                }
            }],
        };
        let config = {
            height:'300px',
            width: '100%',
            margin: 'auto'
        }
        return {layout, datasets, scales, components, config};
    }

    populateBreakDownChart(breakdown) {
        let layout = [
            ['yAxis','heatMap'],
            [null, 'xAxis'],
            [null, 'colorLegend']
        ];
        let scales = [{
            scaleId: 'xScale',
            type: 'Category'
        }, {
            scaleId: 'yScale',
            type: 'Category'
        }, {
            scaleId: 'colorScale',
            type: 'InterpolatedColor'
        }];
        let datasets = [{
            datasetId: 'breakdown',
            data: breakdown
        }];
        let components = {
            plots: [{
                plotId: 'heatMap',
                type: 'Rectangle',
                x: {
                    value: (d, i, ds) => d.hour,
                    scale: 'xScale'
                },
                y: {
                    value: (d, i, ds) => d.day,
                    scale: 'yScale'
                },
                attrs: [{
                    attr: 'fill',
                    value: (d) => d.count,
                    scale: 'colorScale'
                }, {
                    attr: 'stroke',
                    value: 'rgba(200, 200, 200, 0.3)'
                }],
                datasets: ['breakdown'],
                tooltip: (c, p, data, ds) => `missed ${data.count} ${this.props.type} on ${data.day} at ${data.hour}`
            }],
            axes: [{
                axisId: 'xAxis',
                type: 'Category',
                scale: 'xScale',
                orientation: 'bottom',
            },{
                axisId: 'yAxis',
                type: 'Category',
                scale: 'yScale',
                orientation: 'left',
                }
            ],
            legends: [{
                legendId: 'colorLegend',
                type: 'InterpolatedColor',
                plotIds: ['heatMap'],
                colorScaleId: 'colorScale',
            }],
        };
        let config = {
            height:'300px',
            width: '100%',
            margin: 'auto'
        }
        return {layout, datasets, scales, components, config};
    }
    render () {
        let data = this.populateChartData();
        return [
            <Card key='timeilne'>
                <CardHeader className='text-capitalize'>History Missed {this.props.type}</CardHeader>
                <CardBody>
                    <PChart {...this.populateTimelineChart(data.timeline)} />
                </CardBody>
            </Card>,
            <Card key='breakdown'>
                <CardHeader className='text-capitalize'>Missed {this.props.type} By Time of Day</CardHeader>
                <CardBody>
                    <PChart {...this.populateBreakDownChart(data.breakdown)} />
                </CardBody>
            </Card>
        ]
    }
}