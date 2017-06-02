import { combineReducers } from 'redux'
import { RECEIVE_HOMEPAGE } from '../actions'

const headerReducer = (state = {title:"", subtitle:""}, action) => ({ state })
const featuresReducer = (state = {title:"", subtitle:"", highlights:[]}, action) => ({ state })
const testimonialsReducer = (state = {list:[]}, action) => ({ state })
const pricingReducer = (state = {title:"", subtitle:"", list:[]}, action) => ({ state })

const homepageReducer = (state = {
  features: {title:"", subtitle:"", highlights:[]},
  header: {title:"", subtitle:""},
  testimonials: {list:[]},
  pricing: {title:"", subtitle:"", list:[]}
}, action) => {
  switch(action.type){
    case RECEIVE_HOMEPAGE:
      return Object.assign({}, state, { ...action })
    default:
      return state
  }
}

const rootReducer = combineReducers({
  homepage:homepageReducer,
  header:headerReducer,
  features:featuresReducer,
  testimonials:testimonialsReducer,
  pricing:pricingReducer
})

export default rootReducer
