import React from 'react';

export default class KeybaseCheck extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            verified: false,
            error: "",
            keybaseUrl: ""
        }
    }

    componentDidMount(){
        if (this.props.identity != ""){
            fetch("https://keybase.io/_/api/1.0/user/lookup.json?key_suffix="+this.props.identity+"&fields=basics")
                .then(response => response.json())
                .then(data => {
                    if (data.status.code > 0){
                        this.setState({
                            error: data.status.desc
                        })
                    }
                    else{
                        this.setState({
                            verified:true
                        })
                        if (data.them.length > 0){
                            this.setState({keybaseUrl:"https://keybase.io/"+data.them[0].basics.username});
                        }
                    }
                });
        }
    }

    render() {
        if (this.props.identity == ""){
            return <span/>;
        }
        else{
            if (this.state.verified){
                return <a href={this.state.keybaseUrl} target="_blank">{(this.props.showKey?this.props.identity+" ":'')}<i className="fas fa-check-circle text-success" title={this.props.identity}></i></a>
            }
            else{
                return <span>{(this.props.showKey?this.props.identity+" ":'')}<i className="fas fa-check-circle text-warning" title="Keybase identity not verified."></i></span>
            }
        
        }
    }
}