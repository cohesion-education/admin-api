package auth

import (
	"context"
	"fmt"
	"net/http"

	jose "gopkg.in/square/go-jose.v2"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var emailToProfileIDCache map[string]int64

func IsAdmin(r *render.Render) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		profile, ok := cohesioned.FromRequest(req)
		if !ok {
			resp := cohesioned.NewAPIErrorResponse("failed to get current user from request context")
			fmt.Printf("%s\n", resp.ErrMsg)
			r.JSON(w, http.StatusUnauthorized, resp)
			return
		}

		if !profile.IsAdmin() {
			r.JSON(w, http.StatusForbidden, &cohesioned.APIResponse{ErrMsg: "You are not authorized to access this resource"})
			return
		}

		next(w, req)
	}
}

func CheckJwt(r *render.Render, repo profile.Repo, cfg *config.AuthConfig) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		jwksURI := fmt.Sprintf("%s/.well-known/jwks.json", cfg.Domain)
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: jwksURI})
		audience := []string{cfg.ClientID}
		issuer := fmt.Sprintf("%s/", cfg.Domain)

		//fmt.Printf("CheckJwt config::\n\tjwksURI: %s\n\tissuer: %s\n\taudience: %v\n", jwksURI, issuer, audience)

		configuration := auth0.NewConfiguration(client, audience, issuer, jose.RS256)
		validator := auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(req)
		if err != nil {
			resp := cohesioned.NewAPIErrorResponse("Missing or invalid token. %v", err)
			fmt.Printf("%s\n", resp.ErrMsg)
			r.JSON(w, http.StatusUnauthorized, resp)
		} else {
			profile := &cohesioned.Profile{}
			if err := validator.Claims(req, token, profile); err != nil {
				resp := cohesioned.NewAPIErrorResponse("Failed to read claims: %v", err)
				fmt.Printf("%s\n", resp.ErrMsg)
				r.JSON(w, http.StatusUnauthorized, resp)
				return
			}

			profile.ID = getProfileID(repo, profile.Email)
			ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
			next(w, req.WithContext(ctx))
		}
	}
}

func getProfileID(repo profile.Repo, email string) int64 {
	if val, ok := emailToProfileIDCache[email]; ok {
		return val
	}

	profile, err := repo.FindByEmail(email)
	if err != nil {
		fmt.Printf("Failed to find user by email address: %v\n", err)
		return 0
	}

	return profile.ID
}
