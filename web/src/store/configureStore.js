import { createStore, combineReducers, applyMiddleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { routerReducer, routerMiddleware } from 'react-router-redux'
import history from '../history'
import {
  headerReducer,
  featuresReducer,
  testimonialsReducer,
  pricingReducer
} from './homepageReducers'

const loggerMiddleware = createLogger()

// Add the reducer to your store on the `router` key
// Also apply our middleware for navigating
export default function configureStore(preloadedState) {
  return createStore(
    combineReducers({
      router: routerReducer,
      header:headerReducer,
      features:featuresReducer,
      testimonials:testimonialsReducer,
      pricing:pricingReducer
    }),
    preloadedState,
    applyMiddleware(
      thunkMiddleware, // lets us dispatch() functions
      loggerMiddleware, // neat middleware that logs actions
      routerMiddleware(history) // Build the middleware for intercepting and dispatching navigation actions
    )
  )
}
