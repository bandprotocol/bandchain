import qs from 'querystring';
import React,{ Component } from 'react';
import { HTTP } from 'meteor/http'
import {
    Badge,
    Button,
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    NavLink,
    // Input,
    // InputGroup,
    // InputGroupAddon,
    // Button,
    UncontrolledDropdown,
    UncontrolledPopover,
    PopoverBody,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';
import { Link } from 'react-router-dom';
import SearchBar from './SearchBar.jsx';
import i18n from 'meteor/universe:i18n';
import LedgerModal from '../ledger/LedgerModal.jsx';
import Account from './Account.jsx';

const T = i18n.createComponent();

// Firefox does not support named group yet
// const SendPath = new RegExp('/account/(?<address>\\w+)/(?<action>send)')
// const DelegatePath = new RegExp('/validators?/(?<address>\\w+)/(?<action>delegate)')
// const WithdrawPath = new RegExp('/account/(?<action>withdraw)')

const SendPath = new RegExp('/account/(\\w+)/(send)')
const DelegatePath = new RegExp('/validators?/(\\w+)/(delegate)')
const WithdrawPath = new RegExp('/account/(withdraw)')

const getUser = () => localStorage.getItem(CURRENTUSERADDR)

export default class Header extends Component {
    constructor(props) {
        super(props);

        this.toggle = this.toggle.bind(this);
        this.state = {
            isOpen: false,
            networks: ""
        };
    }

    toggle() {
        this.setState({
            isOpen: !this.state.isOpen
        }, ()=>{
            // console.log(this.state.isOpen);
        });
    }

    toggleSignIn = (value) => {
        this.setState(( prevState) => {
            return {isSignInOpen: value!=undefined?value:!prevState.isSignInOpen}
        })
    }

    handleLanguageSwitch(lang, e) {
        i18n.setLocale(lang)
    }

    componentDidMount(){
        let url = Meteor.settings.public.networks
        if (!url)
            return
        try{
            HTTP.get(url, null, (error, result) => {
                if (result.statusCode == 200){
                    let networks = JSON.parse(result.content);
                    if (networks.length > 0){
                        this.setState({
                            networks: <DropdownMenu>{
                                networks.map((network, i) => {
                                    return <span key={i}>
                                        <DropdownItem header><img src={network.logo} /> {network.name}</DropdownItem>
                                        {network.links.map((link, k) => {
                                            return <DropdownItem key={k} disabled={link.chain_id == Meteor.settings.public.chainId}>
                                                <a href={link.url} target="_blank">{link.chain_id} <Badge size="xs" color="secondary">{link.name}</Badge></a>
                                            </DropdownItem>})}
                                        {(i < networks.length - 1)?<DropdownItem divider />:''}
                                    </span>

                                })
                            }</DropdownMenu>
                        })
                    }
                }
            })
        }
        catch(e){
            console.warn(e);
        }
    }

    signOut () {
        localStorage.removeItem(CURRENTUSERADDR);
        localStorage.removeItem(CURRENTUSERPUBKEY);
        this.props.refreshApp();
    }

    shouldLogin = () => {
        let pathname = this.props.location.pathname
        let groups;
        let match = pathname.match(SendPath) || pathname.match(DelegatePath)|| pathname.match(WithdrawPath);
        if (match) {
            if (match[0] === '/account/withdraw') {
                groups = {action: 'withdraw'}
            } else {
                groups = {address: match[1], action: match[2]}
            }
        }
        let params = qs.parse(this.props.location.search.substr(1))
        return groups || params.signin != undefined
    }

    handleLoginConfirmed = (success) => {
        let groups = this.shouldLogin()
        if (!groups) return
        let redirectUrl;
        let params;
        if (groups) {
            let { action, address } = groups;
            params = {action}
            switch (groups.action) {
                case 'send':
                    params.transferTarget = address
                    redirectUrl = `/account/${address}`
                    break
                case 'withdraw':
                    redirectUrl = `/account/${getUser()}`
                    break;
                case 'delegate':
                    redirectUrl = `/validators/${address}`
                    break;
            }
        } else {
            let location = this.props.location;
            params = qs.parse(location.search.substr(1))
            redirectUrl = params.redirect?params.redirect:location.pathname;
            delete params['redirectUrl']
            delete params['signin']
        }

        let query = success?`?${qs.stringify(params)}`:'';
        this.props.history.push(redirectUrl + query)
    }

    render() {
        let signedInAddress = getUser();
        return (
            <Navbar color="primary" dark expand="lg" fixed="top" id="header">
                <NavbarBrand tag={Link} to="/"><img src="/img/big-dipper.svg" className="img-fluid logo"/> <span className="d-none d-xl-inline-block"><T>navbar.siteName</T>&nbsp;</span><Badge color="secondary"><T>navbar.version</T></Badge> </NavbarBrand>
                <UncontrolledDropdown className="d-inline text-nowrap">
                    <DropdownToggle caret={(this.state.networks !== "")} tag="span" size="sm" id="network-nav">{Meteor.settings.public.chainId}</DropdownToggle>
                    {this.state.networks}
                </UncontrolledDropdown>
                <SearchBar id="header-search" history={this.props.history} />
                <NavbarToggler onClick={this.toggle} />
                <Collapse isOpen={this.state.isOpen} navbar>
                    <Nav className="ml-auto text-nowrap" navbar>
                        <NavItem>
                            <NavLink tag={Link} to="/validators"><T>navbar.validators</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/blocks"><T>navbar.blocks</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/transactions"><T>navbar.transactions</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/proposals"><T>navbar.proposals</T></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/voting-power-distribution"><T>navbar.votingPower</T></NavLink>
                        </NavItem>
                        <NavItem id="user-acconut-icon">
                            {!signedInAddress?<Button className="sign-in-btn" color="link" size="lg" onClick={() => {this.setState({isSignInOpen: true})}}><i className="material-icons">vpn_key</i></Button>:
                                <span>
                                    <span className="d-lg-none">
                                        <i className="material-icons large d-inline">account_circle</i>
                                        <Link to={`/account/${signedInAddress}`}> {signedInAddress}</Link>
                                        <Button className="float-right" color="link" size="sm" onClick={this.signOut.bind(this)}><i className="material-icons">exit_to_app</i></Button>
                                    </span>
                                    <span className="d-none d-lg-block">
                                        <i className="material-icons large">account_circle</i>
                                        <UncontrolledPopover className="d-none d-lg-block" trigger="legacy" placement="bottom" target="user-acconut-icon">
                                            <PopoverBody>
                                                <div><T>accounts.signInText</T></div>
                                                <div><Link className="text-nowrap" to={`/account/${signedInAddress}`}>{signedInAddress}</Link></div>
                                                <Button className="float-right" color="link" onClick={this.signOut.bind(this)}><i className="material-icons">exit_to_app</i><span> <T>accounts.signOut</T></span></Button>
                                            </PopoverBody>
                                        </UncontrolledPopover>
                                    </span>
                                </span>}
                            <LedgerModal isOpen={this.state.isSignInOpen} toggle={this.toggleSignIn} refreshApp={this.props.refreshApp} handleLoginConfirmed={this.shouldLogin()?this.handleLoginConfirmed:null}/>
                        </NavItem>
                        <NavItem>
                            <UncontrolledDropdown inNavbar>
                                <DropdownToggle nav caret>
                                    <T>navbar.lang</T>
                                </DropdownToggle>
                                <DropdownMenu right>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('en-US', e)}><T>navbar.english</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('zh-Hant', e)}><T>navbar.chinese</T></DropdownItem>
                                    <DropdownItem onClick={(e) => this.handleLanguageSwitch('zh-Hans', e)}><T>navbar.simChinese</T></DropdownItem>
                                </DropdownMenu>
                            </UncontrolledDropdown>
                        </NavItem>
                    </Nav>
                </Collapse>
            </Navbar>
        );
    }
}