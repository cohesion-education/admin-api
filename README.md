# api

[![wercker status](https://app.wercker.com/status/20b14e1ac5148b29db7e619c1ffd9b45/s/master "wercker status")](https://app.wercker.com/project/byKey/20b14e1ac5148b29db7e619c1ffd9b45)

The api back-end part of the Cohesion Education platform written using Golang.

## Setup

This project uses Glide as a package manager. It is recommended that you [install Glide](https://github.com/Masterminds/glide#install) to work with this project.

## Run locally

First, install the project.

  glide install

Then, you can run the project.

  go run cmd/srv/main.go


## Build locally

To run tests:

    go test $(glide novendor) --cover

Build locally using the Wercker CLI:

    wercker build --docker-local



## Deploy to Cloud Foundry

Relies on a user provided service called admin-api-auth0. Can be defined as follows:

    cf cups auth0-admin-api -p '{"clientid":"id","secret":"pa55woRD","domain":"domain.auth0.com","callback-url":"http://localhost:3000/callback"}'
