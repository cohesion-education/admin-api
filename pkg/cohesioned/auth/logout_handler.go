package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/unrolled/render"
)

func LogoutHandler(r *render.Render, cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   cohesioned.AuthSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)

		logoutURL, err := url.Parse(fmt.Sprintf("%s/v2/logout", cfg.Domain))
		if err != nil {
			fmt.Printf("Failed to parse logout url %v\n", err)
			//TODO - redirect to 500 page
		}
		params := logoutURL.Query()
		params.Set("returnTo", cfg.LogoutRedirectTo)
		params.Set("client_id", cfg.ClientID)
		logoutURL.RawQuery = params.Encode()

		fmt.Printf("Logout redirect %s\n", logoutURL.String())
		http.Redirect(w, r, logoutURL.String(), http.StatusSeeOther)
	}
}
