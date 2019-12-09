import React from 'react';
import { UncontrolledTooltip } from 'reactstrap';

export const DenomSymbol = (props) => {
    switch (props.denom){
    case "steak":
        return '🥩';
    default:
        return '🍅';
    }
}


export const ProposalStatusIcon = (props) => {
    switch (props.status){
    case 'Passed':
        return <i className="fas fa-check-circle text-success"></i>;
    case 'Rejected':
        return <i className="fas fa-times-circle text-danger"></i>;
    case 'Removed':
        return <i className="fas fa-trash-alt text-dark"></i>
    case 'DepositPeriod':
        return <i className="fas fa-battery-half text-warning"></i>;
    case 'VotingPeriod':
        return <i className="fas fa-hand-paper text-info"></i>;
    default:
        return <i></i>;
    }
}

export const VoteIcon = (props) => {
    switch (props.vote){
    case 'yes':
        return <i className="fas fa-check text-success"></i>;
    case 'no':
        return <i className="fas fa-times text-danger"></i>;
    case 'abstain':
        return <i className="fas fa-user-slash text-warning"></i>;
    case 'no_with_veto':
        return <i className="fas fa-exclamation-triangle text-info"></i>;
    default:
        return <i></i>;
    }
}

export const TxIcon = (props) => {
    if (props.valid){
        return <span className="text-success text-nowrap"><i className="fas fa-check-circle"></i></span>;
    }
    else{
        return <span className="text-danger text-nowrap"><i className="fas fa-times-circle"></i></span>;
    }
}

export class InfoIcon extends React.Component {
    constructor(props) {
        super(props);
        this.ref = React.createRef();
    }

    render() {
        return [
            <i key='icon' className='material-icons info-icon' ref={this.ref}>info</i>,
            <UncontrolledTooltip key='tooltip' placement='right' target={this.ref}>
                {this.props.children?this.props.children:this.props.tooltipText}
            </UncontrolledTooltip>
        ]
    }
}