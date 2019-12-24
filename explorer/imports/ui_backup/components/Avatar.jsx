import React from 'react';
export default class Avatar extends React.Component {
    constructor(props) {
        super(props);
    }

    getProfileUrl() {
        return this.props.profileUrl || "https://ui-avatars.com/api/?rounded=true&size=128&name="+this.props.moniker+"&color=fff&background=aaa"
    }
    getColourHex(address){
        // let hex, i;

        // let result = "";
        // let hexString = '1234567890abcde';
        // for (i=0; i<moniker.length; i++) {
        //     hex = moniker.charCodeAt(i).toString(16);
        //     result += hex;
        // }

        // if (result.length < 6){
        //     let tempRes = "";
        //     for (i=0;i<(6-result.length); i++){
        //         tempRes += hexString.charAt(Math.floor((Math.random() * 16)));
        //     }
        //     result += tempRes;
        // }

        return address.substring(0,6);
    }

    render() {
        return (
            <img src={this.getProfileUrl()} alt={this.props.moniker} className={this.props.list?'moniker-avatar-list img-fluid rounded-circle':'img-fluid rounded-circle'} />
        );
    }
}