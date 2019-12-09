import '/imports/startup/client';
import '/imports/ui/stylesheets/bootstrap.min.css';
import '/imports/ui/stylesheets/pace-theme.css';
import '/imports/ui/stylesheets/flipclock.css';
import '/node_modules/plottable/plottable.css';
import './styles.css';
import App from '/imports/ui/App.jsx';
import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom'
// import ReactDOM from 'react-dom';

import { Meteor } from 'meteor/meteor';
import { render } from 'react-dom';

CURRENTUSERADDR = 'ledgerUserAddress';
CURRENTUSERPUBKEY = 'ledgerUserPubKey';

// import { onPageLoad } from 'meteor/server-render';

Meteor.startup(() => {
    render(<Router><App /></Router>, document.getElementById('app'));
    // render(<Header />, document.getElementById('header'));

    // onPageLoad(async sink => {
    //     const App = (await import('/imports/ui/App.jsx')).default;
    //     ReactDOM.hydrate(
    //         <Router>
    //             <App />
    //         </Router>, document.getElementById('app')
    //     );
    // });
});
