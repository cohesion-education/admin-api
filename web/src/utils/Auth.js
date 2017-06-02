import auth0 from 'auth0-js'

// export const webAuth = new auth0.WebAuth({
//   domain: `${process.env.REACT_APP_AUTH0_DOMAIN}`,
//   clientID: `${process.env.REACT_APP_AUTH0_CLIENT_ID}`,
//   redirectUri: `${process.env.REACT_APP_CALLBACK_URL}`,
//   audience: `https://${process.env.REACT_APP_AUTH0_DOMAIN}/userinfo`,
//   scope: 'openid email'
// })

export default class Auth {
  webAuth = new auth0.WebAuth({
    domain: process.env.REACT_APP_AUTH0_DOMAIN,
    clientID: process.env.REACT_APP_AUTH0_CLIENT_ID,
    redirectUri: process.env.REACT_APP_CALLBACK_URL,
    audience: `https://${process.env.REACT_APP_AUTH0_DOMAIN}/userinfo`,
    responseType: 'code',
    scope: 'openid email'
  })

  constructor() {
    this.login = this.login.bind(this)
    this.isAuthenticated = this.isAuthenticated.bind(this)
  }


  login(e) {
    e.preventDefault()
    console.log("auth login; delegating to authorize")
    console.log(this.webAuth.baseOptions.domain)
    console.log(this.webAuth.baseOptions.redirectUri)
    this.webAuth.popup.authorize({
      connection: 'google-oauth2'
    })
    // this.webAuth.authorize()
    // this.webAuth.authorize({
    //   connection: ['Username-Password-Authentication', 'google-oauth2', 'facebook']
    // })
  }

  isAuthenticated() {
    // Check whether the current time is past the
    // access token's expiry time
    let expiresAt = JSON.parse(localStorage.getItem('expires_at'));
    return new Date().getTime() < expiresAt;
  }

  setSession(authResult) {
    //TODO - process result from my api
    // if (authResult && authResult.accessToken && authResult.idToken) {
    //   // Set the time that the access token will expire at
    //   let expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
    //   localStorage.setItem('access_token', authResult.accessToken);
    //   localStorage.setItem('id_token', authResult.idToken);
    //   localStorage.setItem('expires_at', expiresAt);
    //   // navigate to the home route
    //   history.replace('/home');
    // }
  }
}
