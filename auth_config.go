package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/sessions"
)

const authSessionCookieName = "auth-session"

type authConfig struct {
	sessionStore        sessions.Store
	sessionStoreAuthKey string
	ClientID            string
	ClientSecret        string
	Domain              string
	CallbackURL         string
}

func (config *authConfig) getCurrentSession(r *http.Request) (*sessions.Session, error) {
	return config.sessionStore.Get(r, authSessionCookieName)
}

func newAuthConfig() (*authConfig, error) {
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
			if sessionStoreAuthKey, ok := auth0Service.CredentialString("session-auth-key"); ok {
				config.sessionStoreAuthKey = sessionStoreAuthKey
			}
		}
	}

	if len(config.ClientID) == 0 {
		config.ClientID = os.Getenv("AUTH0_CLIENT_ID")
	}

	if len(config.ClientSecret) == 0 {
		config.ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
	}

	if len(config.Domain) == 0 {
		config.Domain = os.Getenv("AUTH0_DOMAIN")
	}

	if len(config.CallbackURL) == 0 {
		config.CallbackURL = os.Getenv("AUTH0_CALLBACK_URL")
	}

	if len(config.sessionStoreAuthKey) == 0 {
		config.sessionStoreAuthKey = os.Getenv("SESSION_AUTH_KEY")
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

	if len(config.sessionStoreAuthKey) == 0 {
		missingConfig = append(missingConfig, "SessionAuthKey")
	}

	if len(missingConfig) > 0 {
		return nil, fmt.Errorf("Failed to load auth0 service from either VCAP_SERVICES or from environment vars - missing %v", missingConfig)
	}

	sessionStore := sessions.NewCookieStore([]byte(config.sessionStoreAuthKey))
	gob.Register(map[string]interface{}{})
	config.sessionStore = sessionStore

	return config, nil
}
