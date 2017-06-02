import { RECEIVE_HOMEPAGE } from '../actions'

export const headerReducer = (state = {title:"", subtitle:""}, action) => {
  switch(action.type){
    case RECEIVE_HOMEPAGE:
      return Object.assign({}, state, { ...action.header })
    default:
      return state
  }
}

export const featuresReducer = (state = {title:"", subtitle:"", highlights:[]}, action) => {
  switch(action.type){
    case RECEIVE_HOMEPAGE:
      return Object.assign({}, state, { ...action.features })
    default:
      return state
  }
}
export const testimonialsReducer = (state = {list:[]}, action) => {
  switch(action.type){
    case RECEIVE_HOMEPAGE:
      return Object.assign({}, state, { ...action.testimonials })
    default:
      return state
  }
}
export const pricingReducer = (state = {title:"", subtitle:"", list:[]}, action) => {
  switch(action.type){
    case RECEIVE_HOMEPAGE:
      return Object.assign({}, state, { ...action.pricing })
    default:
      return state
  }
}
