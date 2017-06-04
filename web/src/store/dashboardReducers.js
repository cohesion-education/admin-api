import { RECEIVE_PROFILE } from '../actions'

export const profileReducer = (state = {picture:""}, action) => {
  switch(action.type){
    case RECEIVE_PROFILE:
      return Object.assign({}, state, { ...action.profile })
    default:
      return state
  }
}
