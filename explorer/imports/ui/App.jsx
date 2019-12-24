import React from 'react'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'

/* Components */
import Header from './components/Header'
import PageContainer from './components/PageContainer'
import Footer from './components/Footer'

/* Pages */
import HomePage from './pages/Home'
import ScriptPage from './pages/ReqScript'
import TxDetailPage from './pages/Transaction'

export default () => (
  <>
    <Header />
    <PageContainer content>
      <Router>
        <Switch>
          <Route exact path="/" component={HomePage} />
          <Route exact path="/script/:codeHash" component={ScriptPage} />
          <Route exact path="/tx/:txHash" component={TxDetailPage} />
        </Switch>
      </Router>
    </PageContainer>
    <Footer />
  </>
)
