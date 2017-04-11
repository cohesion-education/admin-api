package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

func callbackHandler(config *authConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := config.Domain

		conf := &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.CallbackURL,
			Scopes:       []string{"openid", "profile"},

			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + domain + "/authorize",
				TokenURL: "https://" + domain + "/oauth/token",
			},
		}

		code := r.URL.Query().Get("code")

		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Getting now the userInfo
		client := conf.Client(oauth2.NoContext, token)
		resp, err := client.Get("https://" + domain + "/userinfo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var profile map[string]interface{}
		if err = json.Unmarshal(raw, &profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := config.getCurrentSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["id_token"] = token.Extra("id_token")
		session.Values["access_token"] = token.AccessToken
		session.Values["profile"] = profile
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
