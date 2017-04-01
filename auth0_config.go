package main

import (
	"log"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/joho/godotenv"
)

type auth0Config struct {
	ClientID     string
	ClientSecret string
	Domain       string
	CallbackURL  string
}

func newAuth0Config() *auth0Config {
	config := &auth0Config{}

	if appEnv, err := cfenv.Current(); err == nil {
		if auth0Service, err := appEnv.Services.WithName("auth0-admin-api"); err == nil {
			if clientID, ok := auth0Service.CredentialString("clientid"); ok {
				config.ClientID = clientID
			}
			if clientSecret, ok := auth0Service.CredentialString("secret"); ok {
				config.ClientSecret = clientSecret
			}
			if domain, ok := auth0Service.CredentialString("domain"); ok {
				config.Domain = domain
			}
			if callbackURL, ok := auth0Service.CredentialString("callbak-url"); ok {
				config.CallbackURL = callbackURL
			}
		}
	}

	if len(config.ClientID) == 0 {
		if err := godotenv.Load(); err == nil {
			config.ClientID = os.Getenv("AUTH0_CLIENT_ID")
		}
	}

	if len(config.ClientSecret) == 0 {
		if err := godotenv.Load(); err == nil {
			config.ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
		}
	}

	if len(config.Domain) == 0 {
		if err := godotenv.Load(); err == nil {
			config.Domain = os.Getenv("AUTH0_DOMAIN")
		}
	}

	if len(config.CallbackURL) == 0 {
		if err := godotenv.Load(); err == nil {
			config.CallbackURL = os.Getenv("AUTH0_CALLBACK_URL")
		}
	}

	var missingConfig []string
	if len(config.ClientID) == 0 {
		missingConfig = append(missingConfig, "ClientID")
	}

	if len(config.ClientSecret) == 0 {
		missingConfig = append(missingConfig, "ClientSecret")
	}

	if len(config.Domain) == 0 {
		missingConfig = append(missingConfig, "Domain")
	}

	if len(config.CallbackURL) == 0 {
		missingConfig = append(missingConfig, "CallbackURL")
	}

	if len(missingConfig) > 0 {
		log.Fatalf("Failed to load auth0 %v from either VCAP_SERVICES or from .env file", missingConfig)
	}

	return config
}
