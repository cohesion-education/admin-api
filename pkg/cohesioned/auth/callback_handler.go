package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"

	"golang.org/x/oauth2"
)

func CallbackHandler(cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conf := &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.CallbackURL,
			Scopes:       []string{"openid", "profile"},

			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.Domain + "/authorize",
				TokenURL: cfg.Domain + "/oauth/token",
			},
		}

		code := req.URL.Query().Get("code")

		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Error(w, "error exchanging code for token "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Getting now the userInfo
		client := conf.Client(oauth2.NoContext, token)
		resp, err := client.Get(cfg.Domain + "/userinfo")
		if err != nil {
			http.Error(w, "error getting user info "+err.Error(), http.StatusInternalServerError)
			return
		}

		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, "error reading userinfo response body "+err.Error(), http.StatusInternalServerError)
			return
		}

		var profile cohesioned.Profile
		if err = json.Unmarshal(raw, &profile); err != nil {
			http.Error(w, "error unmarshalling userinfo response body "+err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := cfg.GetCurrentSession(req)
		if err != nil {
			http.Error(w, "unable to initialize Session "+err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["id_token"] = token.Extra("id_token")
		session.Values["access_token"] = token.AccessToken
		session.Values[cohesioned.CurrentUserSessionKey] = profile
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, "failed to save Session "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
	}
}
