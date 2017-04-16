package auth

import (
	"context"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/urfave/negroni"
)

func IsAuthenticatedHandler(cfg *config.AuthConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := cfg.GetCurrentSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		profile, ok := session.Values[config.CurrentUserSessionKey]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, config.CurrentUserKey, profile)

			next(w, r.WithContext(ctx))
		}
	}
}
