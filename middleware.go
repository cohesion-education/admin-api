package main

import (
	"context"
	"net/http"

	"github.com/urfave/negroni"
)

func isAuthenticatedHandler(config *authConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := config.getCurrentSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		profile, ok := session.Values["profile"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "current-user", profile)

			next(w, r.WithContext(ctx))
		}
	}
}
