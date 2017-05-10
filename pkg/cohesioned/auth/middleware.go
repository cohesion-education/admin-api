package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/urfave/negroni"
)

func IsAuthenticatedHandler(cfg *config.AuthConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := cfg.GetCurrentSession(r)
		if err != nil {
			//TODO - 401
			http.Error(w, "Failed to get current session "+err.Error(), http.StatusInternalServerError)
			return
		}

		profile, ok := session.Values[cohesioned.CurrentUserSessionKey]
		if !ok {
			//TODO - 401
			http.Error(w, "Failed to get current user from session", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), cohesioned.CurrentUserKey, profile)
		next(w, r.WithContext(ctx))
	}
}

func IsAdmin(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	profile, err := cohesioned.GetProfile(req)
	if err != nil {
		//TODO - 401
		http.Error(w, fmt.Sprintf("Failed to get current user %v", err), http.StatusInternalServerError)
		return
	}

	if profile.HasRole("admin") {
		next(w, req)
		return
	}

	//TODO - 403
	http.Error(w, "You are not authorized to view this content", http.StatusUnauthorized)
}
