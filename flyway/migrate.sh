#!/bin/bash

source ../.env

 flyway -user=$AWS_RDS_USERNAME -password=$AWS_RDS_PASSWORD -url=jdbc:mysql://$AWS_RDS_HOST/$AWS_RDS_DBNAME -locations=filesystem:./sql/ migrate
