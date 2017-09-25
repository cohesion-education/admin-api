package config_test

import "fmt"

const auth0ServicePayload = `{
  "credentials":{
    "clientid":"test-client-id",
    "secret":"test-secret",
    "domain":"test-domain",
    "callback-url":"test-callback-url",
    "logout-redirect-to":"test-logout-url",
    "session-auth-key":"abc123"
  },
  "syslog_drain_url": "",
  "volume_mounts": [],
  "label": "user-provided",
  "name": "auth0-admin",
  "tags": []
}`

const awsServicePayload = `{
  "credentials":{
    "region":"us-east-1",
    "access_key_id": "abc123",
    "secret_access_key": "abc123",
    "session_token": "",
    "s3_video_bucket": "test-bucket",
    "rds_username":"user",
    "rds_password":"pass",
    "rds_host":"localhost",
    "rds_port": "3309",
    "rds_dbname": "dbname"
  },
  "syslog_drain_url": "",
  "volume_mounts": [],
  "label": "user-provided",
  "name": "aws",
  "tags": []
}`

const auth0ServicePartialPayload = `{
  "credentials":{
    "secret":"test-secret",
    "domain":"test-domain",
    "callback-url":"test-callback-url",
    "logout-redirect-to":"test-logout-url",
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
    %s,
    %s
  ]
}`, auth0ServicePayload, awsServicePayload)

var vcapServicesPartialPayload = fmt.Sprintf(`{
  "user-provided": [
    %s
  ]
}`, auth0ServicePartialPayload)
