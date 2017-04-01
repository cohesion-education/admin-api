package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/urfave/negroni"
)

func isAuthenticatedHandler(store sessions.Store) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values["profile"]; !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			next(w, r)
		}
	}
}
