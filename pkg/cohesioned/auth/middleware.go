package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/urfave/negroni"
)

func IsAuthenticatedHandler(cfg *config.AuthConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		session, err := cfg.GetCurrentSession(req)
		if err != nil {
			fmt.Printf("error getting current user from session %v\n", err)
			http.Redirect(w, req, "/401", http.StatusSeeOther)
			return
		}

		profile, ok := session.Values[cohesioned.CurrentUserSessionKey]
		if !ok {
			http.Redirect(w, req, "/401", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
		next(w, req.WithContext(ctx))
	}
}

func IsAdmin(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	profile, err := cohesioned.GetProfile(req)
	if err != nil {
		fmt.Printf("error getting current user from request context %v\n", err)
		http.Redirect(w, req, "/401", http.StatusSeeOther)
		return
	}

	if profile.HasRole("admin") {
		next(w, req)
		return
	}

	http.Redirect(w, req, "/403", http.StatusSeeOther)
}
