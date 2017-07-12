package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"

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
			fmt.Printf("error exchanging code for token %s %v\n", code, err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		// Getting now the userInfo
		client := conf.Client(oauth2.NoContext, token)
		resp, err := client.Get(cfg.Domain + "/userinfo")
		if err != nil {
			fmt.Printf("error getting ''%s/userinfo' %v\n", cfg.Domain, err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("error reading userinfo response body %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		var profile cohesioned.Profile
		if err = json.Unmarshal(raw, &profile); err != nil {
			fmt.Printf("error unmarshalling userinfo response body %v\n%s\n", err, string(raw))
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		session, err := cfg.GetCurrentSession(req)
		if err != nil {
			fmt.Printf("unable to initialize Session %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		session.Values["id_token"] = token.Extra("id_token")
		session.Values["access_token"] = token.AccessToken
		session.Values[cohesioned.CurrentUserSessionKey] = profile
		err = session.Save(req, w)
		if err != nil {
			fmt.Printf("failed to save Session %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		if profile.HasRole("admin") {
			http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
	}
}
