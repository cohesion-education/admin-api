# admin-api

[![wercker status](https://app.wercker.com/status/20b14e1ac5148b29db7e619c1ffd9b45/s/master "wercker status")](https://app.wercker.com/project/byKey/20b14e1ac5148b29db7e619c1ffd9b45)


## Run locally
There are two parts to the app - the api written in Go, and the front-end written in React. The React app must be built first, and is served up by the Go app as a static resource.

  cd web && yarn install && yarn build && cd ..

  go run cmd/srv/main.go


## Build locally
To run tests:

    go test $(glide novendor) --cover

    cd web && yarn test -- --coverage

Build locally using the Wercker CLI:

    wercker --environment wercker.env build --docker-local



## Deploy to Cloud Foundry

Relies on a user provided service called admin-api-auth0. Can be defined as follows:

    cf cups auth0-admin-api -p '{"clientid":"id","secret":"pa55woRD","domain":"domain.auth0.com","callback-url":"http://localhost:3000/callback"}'
