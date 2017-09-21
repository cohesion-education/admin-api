// +build unit

package auth_test

const userInfoPayload = `{
  "email_verified": false,
  "email": "test.account@userinfo.com",
  "clientID": "q2hnj2iu...",
  "updated_at": "2016-12-05T15:15:40.545Z",
  "name": "test.account@userinfo.com",
  "picture": "https://s.gravatar.com/avatar/dummy.png",
  "user_id": "auth0|58454...",
  "nickname": "test.account",
  "identities": [
    {
      "user_id": "58454...",
      "provider": "auth0",
      "connection": "Username-Password-Authentication",
      "isSocial": false
    }
  ],
  "created_at": "2016-12-05T11:16:59.640Z",
  "sub": "auth0|58454..."
}`

const adminUserInfoPayload = `{
  "email_verified": false,
  "email": "test.account@cohesioned.io",
  "clientID": "q2hnj2iu...",
  "updated_at": "2016-12-05T15:15:40.545Z",
  "name": "test.account@userinfo.com",
  "picture": "https://s.gravatar.com/avatar/dummy.png",
  "user_id": "auth0|58454...",
  "nickname": "test.account",
  "identities": [
    {
      "user_id": "58454...",
      "provider": "auth0",
      "connection": "Username-Password-Authentication",
      "isSocial": false
    }
  ],
  "created_at": "2016-12-05T11:16:59.640Z",
  "sub": "auth0|58454...",
  "app_metadata": {
    "roles": ["admin"]
  }
}`
