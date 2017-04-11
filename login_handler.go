package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

func loginViewHandler(r *render.Render, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "auth-session")
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		if profile := session.Values["profile"]; profile != nil {
			http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
			return
		}

		r.HTML(w, http.StatusOK, "ubold/login", nil, render.HTMLOptions{Layout: ""})
	}
}
