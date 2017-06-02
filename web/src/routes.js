import React from 'react'
import { createStore, combineReducers, applyMiddleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { Provider } from 'react-redux'
import createHistory from 'history/createBrowserHistory'
import { Route } from 'react-router'
import { ConnectedRouter, routerReducer, routerMiddleware } from 'react-router-redux'
import Homepage from './views/Homepage'
import { fetchHomepage } from './actions'
import { headerReducer, featuresReducer, testimonialsReducer, pricingReducer } from './store/homepageReducers'

// Create a history of your choosing (we're using a browser history in this case)
const history = createHistory()

const loggerMiddleware = createLogger()

// Add the reducer to your store on the `router` key
// Also apply our middleware for navigating
const store = createStore(
  combineReducers({
    router: routerReducer,
    header:headerReducer,
    features:featuresReducer,
    testimonials:testimonialsReducer,
    pricing:pricingReducer
  }),
  applyMiddleware(
    thunkMiddleware, // lets us dispatch() functions
    loggerMiddleware, // neat middleware that logs actions
    routerMiddleware(history) // Build the middleware for intercepting and dispatching navigation actions
  )
)

store.dispatch(fetchHomepage())

const Routes = () => (
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <Route exact path="/" component={Homepage} />
    </ConnectedRouter>
  </Provider>
)

export default Routes
