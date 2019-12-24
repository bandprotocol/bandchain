import React from 'react';
import { Tooltip } from 'reactstrap';

export default class Block extends React.Component {
    constructor(props) {
        super(props);

        this.toggle = this.toggle.bind(this);
        this.state = {
            tooltipOpen: false
        };
    }

    toggle() {
        this.setState({
            tooltipOpen: !this.state.tooltipOpen
        });
    }

    render() {
        return (
            <div className="block">
                <div id={"block-"+this.props.height} className={this.props.exists?'full':'empty'}></div>
                <Tooltip placement="top" isOpen={this.state.tooltipOpen} autohide={false} target={"block-"+this.props.height} toggle={this.toggle}>
                    {this.props.height}
                </Tooltip>
            </div>
        );
    }
}