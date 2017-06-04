import React from 'react'
import Auth from './Auth'

export default class Callback extends React.Component {
  auth = new Auth()

  constructor(props){
    super(props)
    console.log('Callback.constructor')
    this.handleAuthentication = this.handleAuthentication.bind(this)
  }

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
