package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/urfave/negroni"
)

func IsAuthenticatedHandler(cfg *config.AuthConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := cfg.GetCurrentSession(r)
		if err != nil {
			http.Error(w, "Failed to get current session"+err.Error(), http.StatusInternalServerError)
			return
		}

		profile, ok := session.Values[config.CurrentUserSessionKey]
		if !ok {
			http.Error(w, "Failed to get current user from session", http.StatusInternalServerError)
			return
		}

		fmt.Printf("profile: %v\n", profile)
		ctx := context.WithValue(r.Context(), config.CurrentUserKey, profile)
		next(w, r.WithContext(ctx))
	}
}
