import React from 'react'
import Auth from './Auth'

export default class Logout extends React.Component {
  auth = new Auth()

  componentWillMount(){
    this.auth.logout()
  }

  render (){
    return (
      <div>Logging out...</div>
    )
  }
}
