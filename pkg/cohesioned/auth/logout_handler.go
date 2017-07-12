package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
)

func LogoutHandler(cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}
		params := logoutURL.Query()
		params.Set("returnTo", cfg.LogoutRedirectTo)
		params.Set("client_id", cfg.ClientID)
		logoutURL.RawQuery = params.Encode()

		http.Redirect(w, req, logoutURL.String(), http.StatusSeeOther)
	}
}
