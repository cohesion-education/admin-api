import Auth from './Auth'

it('loads a configured auth0', () => {
  // const reactAppKeys = Object.keys(process.env).filter(key => key.startsWith('REACT_APP_'))
  // reactAppKeys.map(key => console.log(`${key} = ${process.env[key]}`))

  const auth = new Auth()

  // console.log(Object.keys(auth.webAuth.baseOptions))

  expect(auth).toBeDefined()
  expect(auth.webAuth).toBeDefined()
  expect(auth.webAuth.baseOptions.domain).toBe(process.env.REACT_APP_AUTH0_DOMAIN)
  expect(auth.webAuth.baseOptions.clientID).toBe(process.env.REACT_APP_AUTH0_CLIENT_ID)
  expect(auth.webAuth.baseOptions.redirectUri).toBe(process.env.REACT_APP_CALLBACK_URL)
  expect(auth.login).toBeDefined()
  expect(auth.isAuthenticated).toBeDefined()
});
