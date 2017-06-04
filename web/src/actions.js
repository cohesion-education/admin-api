import fetch from 'isomorphic-fetch'

export const REQUEST_HOMEPAGE = 'REQUEST_HOMEPAGE'
export const RECEIVE_HOMEPAGE = 'RECEIVE_HOMEPAGE'
export const RECEIVE_PROFILE = 'RECEIVE_PROFILE'


function requestHomepage() {
  return {
    type: REQUEST_HOMEPAGE
  }
}

function receiveHomepage(json) {
  return {
    type: RECEIVE_HOMEPAGE,
    header:json.header,
    features:json.features,
    testimonials:json.testimonials,
    pricing:json.pricing,
    receivedAt: Date.now()
  }
}

export function receiveProfile(profile) {
  return {
    type: RECEIVE_PROFILE,
    profile:profile,
    receivedAt: Date.now()
  }
}

export function fetchHomepage() {
  return function (dispatch) {
    dispatch(requestHomepage())

    return fetch('/api/homepage')
      .then(response => {
        return response.json()
      })
      .then(json =>
        dispatch(receiveHomepage(json))
      ).catch(function(err) {
        console.log(`an error occurred when making the request: ${err}`)
      })
  }
}
