import React, { Component } from 'react';
import { FormGroup, Label, Col, Input, Button, Modal, ModalBody, ModalHeader } from 'reactstrap';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

const getFullUrl = (url) => {
	if (url.substr(0, 4) === 'http')
		return url
	else
		return location.origin + url
}

export default class LinkIcon extends Component{
    constructor(props){
        super(props);
        this.state = {
        	isOpen: false,
        	message: 'Click to copy link'
        }
    }

	toggleModal = () => {
		this.setState({isOpen: !this.state.isOpen})
	}

	copyLink = (e) => {
		e.currentTarget.select()
    	document.execCommand('copy');
    	this.setState({message: 'Link copied'})
	}

	renderTextarea(name, value) {
		return <Input type='url' name={name} value={getFullUrl(value)} onClick={this.copyLink} readOnly/>
	}

	renderLink = (link, i) => {
		let { label, url } = link
		return <FormGroup key={i} row>
          <Label for={`sublink-${i}`} sm={2}>{label}</Label>
          <Col sm={10}>
            {this.renderTextarea(`sublink-${i}`, url)}
          </Col>
        </FormGroup>
	}

	render() {
		return <span>
			<button type="button" className='close'><i className='material-icons' onClick={this.toggleModal}>share</i></button>
			<Modal isOpen={this.state.isOpen} toggle={this.toggleModal} className="share-link-modal">
				<ModalHeader toggle={this.toggleModal}>Share</ModalHeader>
            	<ModalBody>
                	<span> {this.state.message} </span>
                	{this.renderTextarea('primary-link', this.props.link)}
                	{this.props.otherLinks && this.props.otherLinks.map(this.renderLink)}
            	</ModalBody>
        	</Modal>
		</span>
	}
}