package config_test

import "fmt"

const auth0ServicePayload = `{
  "credentials":{
    "clientid":"test-client-id",
    "secret":"test-secret",
    "domain":"test-domain",
    "callback-url":"test-callback-url",
    "session-auth-key":"abc123"
  },
  "syslog_drain_url": "",
  "volume_mounts": [],
  "label": "user-provided",
  "name": "auth0-admin",
  "tags": []
}`

const auth0ServicePartialPayload = `{
  "credentials":{
    "secret":"test-secret",
    "domain":"test-domain",
    "callback-url":"test-callback-url",
    "session-auth-key":"abc12345"
  },
  "syslog_drain_url": "",
  "volume_mounts": [],
  "label": "user-provided",
  "name": "auth0-admin",
  "tags": []
}`

var vcapApplicationPayload = `{
  "application_name": "test-app",
  "name": "test-app"
}`

var vcapServicesPayload = fmt.Sprintf(`{
  "user-provided": [
    %s
  ]
}`, auth0ServicePayload)

var vcapServicesPartialPayload = fmt.Sprintf(`{
  "user-provided": [
    %s
  ]
}`, auth0ServicePartialPayload)
