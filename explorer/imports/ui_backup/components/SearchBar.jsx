import React,{ Component } from 'react';
import { Input, InputGroup, InputGroupAddon, Button } from 'reactstrap';

export default class SearchBar extends Component {
    constructor(props) {
        super(props);
        this.state = {
            queryString: ""
        }
    }

    redirectSearchResult = (query) => {
        let hashRegEx = new RegExp(/[0-9A-F]{64}$/, 'igm');
        let validatorRegEx = new RegExp(Meteor.settings.public.bech32PrefixValAddr+'.*$', 'igm');
        let accountRegEx = new RegExp(Meteor.settings.public.bech32PrefixAccAddr+'.*$', 'igm');
        if (query != ""){
            if (!isNaN(query)){
                this.props.history.push('/blocks/'+query);
            }
            else if (query.match(hashRegEx)){
                this.props.history.push('/transactions/'+query);
            }
            else if (query.match(validatorRegEx)){
                this.props.history.push('/validator/'+query);
            }
            else if (query.match(accountRegEx)){
                this.props.history.push('/account/'+query);
            }
            else{
                // handle not found
            }
        }
    }

    handleInput = (e) => {
        this.setState({
            queryString: e.target.value
        })
    }

    handleMobileSearch = (e) => {
        this.redirectSearchResult(this.state.queryString);
        this.setState({
            queryString: ""
        })
    }

    handleSearch = (e) => {
        if (e.key === 'Enter') {
            this.redirectSearchResult(e.target.value);
            this.setState({
                queryString: ""
            })
        }
    }

    render(){
        return <InputGroup className={(this.props.mobile)?'d-lg-none':'d-none d-lg-flex'} id={this.props.id}>
            <Input id="queryString" value={this.state.queryString} onChange={this.handleInput} placeholder={i18n.__('common.searchPlaceholder')} onKeyDown={this.handleSearch}/>
            {(this.props.mobile)?<InputGroupAddon addonType="append"><Button><i className="material-icons" onClick={this.handleMobileSearch}>search</i></Button></InputGroupAddon>:''}
        </InputGroup>
    }
}