# admin-api

[![wercker status](https://app.wercker.com/status/20b14e1ac5148b29db7e619c1ffd9b45/s/master "wercker status")](https://app.wercker.com/project/byKey/20b14e1ac5148b29db7e619c1ffd9b45)

### Build locally using the Wercker CLI

  wercker --environment wercker.env build --docker-local



### Deploy to Cloud Foundry

Relies on a user provided service called admin-api-auth0. Can be defined as follows:

  cf cups auth0-admin-api -p '{"clientid":"id","secret":"pa55woRD","domain":"cohesioned.auth0.com","callback-url":"http://localhost:3000/callback"}'
