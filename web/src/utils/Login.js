import React from 'react'
import Auth from './Auth'

export default class Login extends React.Component {
  auth = new Auth()

  componentWillMount(){
    this.auth.login()
  }

  render (){
    return (
      <div>Redirecting to login...</div>
    )
  }
}
