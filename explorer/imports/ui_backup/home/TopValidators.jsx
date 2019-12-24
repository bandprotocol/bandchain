import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Table, Row, Col, Card, CardImg, CardText, CardBody,
    CardTitle, CardSubtitle, Button, Progress, Spinner } from 'reactstrap';
import numbro from 'numbro';
import Avatar from '../components/Avatar.jsx';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();
export default class TopValidators extends Component{
    constructor(props){
        super(props);
        this.timer = 0;
        this.state = {
            validators: null
        }
    }

    componentDidMount(){
        let self = this;
        self.timer = Meteor.setInterval(function(){
            if (self.props.validators.length> 0){
                let validators = self.shuffle(self.props.validators);
                validators.splice(10, validators.length-10);
                // console.log(validators);
                self.setState({
                    validators: validators.map((validator, i ) => {
                        return <tr key={i}>
                            <td><Link to={"/validator/"+validator.address}>
                                <Avatar moniker={validator.description.moniker} profileUrl={validator.profile_url} address={validator.address} list={true} />
                                {validator.description.moniker}
                            </Link></td>
                            <td className="voting-power">{numbro(validator.voting_power).format('0,0')}</td>
                            <td><Progress animated value={validator.uptime}>{validator.uptime?validator.uptime.toFixed(2):0}%</Progress></td>
                        </tr>
                    })
                })
            }
        },5000);
    }

    componentWillUnmount(){
        Meteor.clearInterval(this.timer);
    }

    // componentDidUpdate(prevState){
    //     if (this.props.status != prevState.status){
    //     }
    // }

    shuffle(a) {
        for (let i = a.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [a[i], a[j]] = [a[j], a[i]];
        }
        return a;
    }

    render(){
        if (this.props.loading){
            return <Spinner type="grow" color="primary" />
        }
        else{
            if (this.props.validatorsExist && this.props.status.prevotes){
                return <Card>
                    <div className="card-header"><T>validators.randomValidators</T></div>
                    <CardBody>
                        <Table striped className="random-validators">
                            <thead>
                                <tr>
                                    <th className="moniker"><i className="material-icons">perm_contact_calendar</i><span className="d-none d-sm-inline"><T>validators.moniker</T></span></th>
                                    <th className="voting-power"><i className="material-icons">power</i><span className="d-none d-sm-inline"><T>common.votingPower</T></span></th>
                                    <th className="uptime"><i className="material-icons">flash_on</i><span className="d-none d-sm-inline"><T>validators.uptime</T></span></th>
                                </tr>
                            </thead>
                            <tbody>{this.state.validators}</tbody>
                        </Table>
                    </CardBody>
                </Card>;
            }
            else{
                return <div></div>
            }

        }
    }
}