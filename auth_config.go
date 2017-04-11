package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const authSessionCookieName = "auth-session"

type authConfig struct {
	sessionStore sessions.Store
	ClientID     string
	ClientSecret string
	Domain       string
	CallbackURL  string
}

func (config *authConfig) getCurrentSession(r *http.Request) (*sessions.Session, error) {
	return config.sessionStore.Get(r, authSessionCookieName)
}

func newAuthConfig() *authConfig {
	config := &authConfig{}

	if appEnv, err := cfenv.Current(); err == nil {
		if auth0Service, err := appEnv.Services.WithName("auth0-admin"); err == nil {
			if clientID, ok := auth0Service.CredentialString("clientid"); ok {
				config.ClientID = clientID
			}
			if clientSecret, ok := auth0Service.CredentialString("secret"); ok {
				config.ClientSecret = clientSecret
			}
			if domain, ok := auth0Service.CredentialString("domain"); ok {
				config.Domain = domain
			}
			if callbackURL, ok := auth0Service.CredentialString("callback-url"); ok {
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

	//TODO - get this from env var
	sessionStore := sessions.NewCookieStore([]byte("todo-inject-me"))
	gob.Register(map[string]interface{}{})

	config.sessionStore = sessionStore

	return config
}
