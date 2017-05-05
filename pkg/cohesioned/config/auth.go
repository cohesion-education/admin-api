package config

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/gorilla/sessions"
)

type AuthConfig struct {
	SessionStore     sessions.Store
	ClientID         string
	ClientSecret     string
	Domain           string
	CallbackURL      string
	LogoutRedirectTo string
}

func (config *AuthConfig) GetCurrentSession(r *http.Request) (*sessions.Session, error) {
	return config.SessionStore.Get(r, cohesioned.AuthSessionCookieName)
}

func newSessionStore(authKey string) sessions.Store {
	if len(authKey) == 0 {
		return nil
	}
	gob.Register(&cohesioned.Profile{})
	sessionStore := sessions.NewCookieStore([]byte(authKey))
	return sessionStore
}

func NewAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{}

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
			if LogoutRedirectTo, ok := auth0Service.CredentialString("logout-redirect-to"); ok {
				config.LogoutRedirectTo = LogoutRedirectTo
			}
			if sessionStoreAuthKey, ok := auth0Service.CredentialString("session-auth-key"); ok {
				config.SessionStore = newSessionStore(sessionStoreAuthKey)
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
		config.CallbackURL = os.Getenv("CALLBACK_URL")
	}

	if len(config.LogoutRedirectTo) == 0 {
		config.LogoutRedirectTo = os.Getenv("LOGOUT_REDIRECT_TO")
	}

	if config.SessionStore == nil {
		config.SessionStore = newSessionStore(os.Getenv("SESSION_AUTH_KEY"))
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

	if len(config.LogoutRedirectTo) == 0 {
		missingConfig = append(missingConfig, "LogoutRedirectTo")
	}

	if config.SessionStore == nil {
		missingConfig = append(missingConfig, "SessionStoreAuthKey")
	}

	if len(missingConfig) > 0 {
		return nil, fmt.Errorf("Failed to load auth0 service from either VCAP_SERVICES or from environment vars - missing %v", missingConfig)
	}

	return config, nil
}
