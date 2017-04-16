package config_test

import "fmt"

const Auth0ServicePayload = `{
  "credentials":{
    "clientid":"test-client-id",
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

const Auth0ServicePartialPayload = `{
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

var VcapApplicationPayload = `{
  "application_name": "test-app",
  "name": "test-app"
}`

var VcapServicesPayload = fmt.Sprintf(`{
  "user-provided": [
    %s
  ]
}`, Auth0ServicePayload)

var VcapServicesPartialPayload = fmt.Sprintf(`{
  "user-provided": [
    %s
  ]
}`, Auth0ServicePartialPayload)
