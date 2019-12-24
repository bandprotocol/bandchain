
import React, { Component } from 'react';
import GoogleTagManager from '/imports/ui/components/GoogleTagManager.jsx';
import { Route, Switch, withRouter } from 'react-router-dom'
import { createMemoryHistory } from 'history';
import { Container } from 'reactstrap';
import Header from '/imports/ui/components/Header.jsx';
import Footer from '/imports/ui/components/Footer.jsx';
import Home from '/imports/ui/home/Home.jsx';
import Validators from '/imports/ui/validators/ValidatorsList.jsx';
import Account from '/imports/ui/accounts/Account.jsx';
import BlocksTable from '/imports/ui/blocks/BlocksTable.jsx';
import Proposals from '/imports/ui/proposals/Proposals.jsx';
import ValidatorDetails from '/imports/ui/validators/ValidatorDetails.jsx';
import Transactions from '/imports/ui/transactions/TransactionsList.jsx';
import Distribution from '/imports/ui/voting-power/Distribution.jsx';
import SearchBar from '/imports/ui/components/SearchBar.jsx';
import moment from 'moment';
import SentryBoundary from '/imports/ui/components/SentryBoundary.jsx';
import NotFound from '/imports/ui/pages/NotFound.jsx';

import { ToastContainer, toast } from 'react-toastify';

if (Meteor.isClient)
    import 'react-toastify/dist/ReactToastify.min.css';

// import './App.js'

const RouteHeader = withRouter( (props) => <Header {...props}/>)
const MobileSearchBar = withRouter( ({history}) => <SearchBar history={history} id="mobile-searchbar" mobile />)

function getLang () {
    return (
        navigator.languages && navigator.languages[0] ||
        navigator.language ||
        navigator.browserLanguage ||
        navigator.userLanguage ||
        'en-US'
    );
}

class App extends Component {
    constructor(props){
        super(props);
    }

    componentDidMount(){
        let lastDay = moment("2019-02-10");
        let now = moment();
        if (now.diff(lastDay) < 0 ){
            toast.error("🐷 Gung Hei Fat Choi! 恭喜發財！");
        }

        let lang = getLang();

        if ((lang.toLowerCase() == 'zh-tw') || (lang.toLowerCase() == 'zh-hk')){
            i18n.setLocale('zh-Hant');
        }
        else if ((lang.toLowerCase() == 'zh-cn') || (lang.toLowerCase() == 'zh-hans-cn') || (lang.toLowerCase() == 'zh')){
            i18n.setLocale('zh-Hans');
        }
        else{
            i18n.setLocale(lang);
        }

    }

    propagateStateChange = () => {
        this.forceUpdate();
    }

    render() {
        const history = createMemoryHistory();

        return(
            // <Router history={history}>
                <div>
                    {(Meteor.settings.public.gtm)?<GoogleTagManager gtmId={Meteor.settings.public.gtm} />:''}
                    <RouteHeader refreshApp={this.propagateStateChange}/>
                    <Container fluid id="main">
                        <ToastContainer />
                        <SentryBoundary>
                            <MobileSearchBar />
                            <Switch>
                                <Route exact path="/" component={Home} />
                                <Route path="/blocks" component={BlocksTable} />
                                <Route path="/transactions" component={Transactions} />
                                <Route path="/account/:address" render={(props)=><Account {...props} />} />
                                <Route path="/validators" exact component={Validators} />
                                <Route path="/validators/inactive" render={(props) => <Validators {...props} inactive={true} />} />
                                <Route path="/voting-power-distribution" component={Distribution} />
                                <Route path="/(validator|validators)" component={ValidatorDetails} />
                                <Route path="/proposals" component={Proposals} />
                                <Route component={NotFound} />
                            </Switch>
                        </SentryBoundary>
                    </Container>
                    <Footer />
                </div>
            // </Router>
        );
    }
}

export default App