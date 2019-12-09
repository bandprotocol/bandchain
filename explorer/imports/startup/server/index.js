// Import server startup through a single index entry point

import './util.js';
import './register-api.js';
import './create-indexes.js';

// import React from 'react';
// import { renderToNodeStream } from 'react-dom/server';
// import { renderToString } from "react-dom/server";
import { onPageLoad } from 'meteor/server-render';
// import { StaticRouter } from 'react-router-dom';
// import { ServerStyleSheet } from "styled-components"
import { Helmet } from 'react-helmet';

// import App from '../../ui/App.jsx';

onPageLoad(sink => {
    // const context = {};
    // const sheet = new ServerStyleSheet()

    // const html = renderToString(sheet.collectStyles(
    //     <StaticRouter location={sink.request.url} context={context}>
    //         <App />
    //     </StaticRouter>
    //   ));

    // sink.renderIntoElementById('app', html);

    const helmet = Helmet.renderStatic();
    sink.appendToHead(helmet.meta.toString());
    sink.appendToHead(helmet.title.toString());

    // sink.appendToHead(sheet.getStyleTags());
});