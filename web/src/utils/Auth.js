import auth0 from 'auth0-js'
import history from '../history'

export default class Auth {
  webAuth = new auth0.WebAuth({
    domain: process.env.REACT_APP_AUTH0_DOMAIN,
    clientID: process.env.REACT_APP_AUTH0_CLIENT_ID,
    redirectUri: process.env.REACT_APP_CALLBACK_URL,
    audience: `https://${process.env.REACT_APP_AUTH0_DOMAIN}/userinfo`,
    responseType: 'token id_token',
    scope: 'openid email'
  })

  userProfile

  constructor() {
    this.login = this.login.bind(this)
    this.logout = this.logout.bind(this)
    this.handleAuthentication = this.handleAuthentication.bind(this)
    this.isAuthenticated = this.isAuthenticated.bind(this)
    this.getAccessToken = this.getAccessToken.bind(this)
    this.getProfile = this.getProfile.bind(this)
  }

  login(e) {
    e.preventDefault()
    this.webAuth.authorize()
  }

  handleAuthentication() {
    this.webAuth.parseHash((err, authResult) => {
      console.log(err)
      console.log(authResult)
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.setSession(authResult)
        history.replace('/dashboard')
      } else if (err) {
        history.replace('/')
        console.log(err)
        alert(`Error: ${err.error}. Check the console for further details.`)
      }
    });
  }


  isAuthenticated() {
    if(!localStorage.getItem('expires_at')){
      return false
    }
    // Check whether the current time is past the
    // access token's expiry time
    let expiresAt = JSON.parse(localStorage.getItem('expires_at'))
    return new Date().getTime() < expiresAt
  }

  setSession(authResult) {
    if (authResult && authResult.accessToken && authResult.idToken) {
      // Set the time that the access token will expire at
      let expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime())
      localStorage.setItem('access_token', authResult.accessToken)
      localStorage.setItem('id_token', authResult.idToken)
      localStorage.setItem('expires_at', expiresAt)
      localStorage.setItem('expires_in', authResult.expiresIn)

      // history.replace('/')
    }
  }

  getAccessToken() {
    const accessToken = localStorage.getItem('access_token')
    if (!accessToken) {
      throw new Error('No access token found')
    }
    return accessToken
  }

  getProfile(cb) {
    try{
      let accessToken = this.getAccessToken()
      this.webAuth.client.userInfo(accessToken, (err, profile) => {
        if (profile) {
          this.userProfile = profile
        }
        cb(err, profile)
      })
    }catch(err){
      cb(err, null)
    }
  }

  logout() {
    // Clear access token and ID token from local storage
    localStorage.removeItem('access_token')
    localStorage.removeItem('id_token')
    localStorage.removeItem('expires_at')
    this.userProfile = null

    history.replace('/')
  }
}
