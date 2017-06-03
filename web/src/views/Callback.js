import React from 'react'
import Auth from '../utils/Auth'

export default class Callback extends React.Component {
  auth = new Auth()

  handleAuthentication = (nextState, replace) => {
    if (/access_token|id_token|error/.test(nextState.location.hash)) {
      this.auth.handleAuthentication()
    }
  }

  componentDidMount(){
    console.log('Callback.componentDidMount')
    this.handleAuthentication(this.props)
  }

  render (){
    return (
      <div>Loading...</div>
    )
  }
}
