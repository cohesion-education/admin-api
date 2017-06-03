import React from 'react'
import { Provider } from 'react-redux'
import { Route } from 'react-router'
import ReactGA from 'react-ga'
import { ConnectedRouter } from 'react-router-redux'
import configureStore from './store/configureStore'
import history from './history'
import Homepage from './containers/Homepage'
import Callback from './views/Callback'

const logPageView = () => {
  ReactGA.set({ page: window.location.pathname + window.location.search });
  ReactGA.pageview(window.location.pathname + window.location.search);
}

const store = configureStore()

const Routes = () => (
  <Provider store={store}>
    <ConnectedRouter history={history} onUpdate={logPageView}>
      <div>
        <Route exact path="/" component={Homepage} />
        <Route path="/callback" comoponent={Callback} />
      </div>
    </ConnectedRouter>
  </Provider>
)

export default Routes
