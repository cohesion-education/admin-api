package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
)

func ConfigHandler(cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		url, err := url.Parse(cfg.Domain)
		if err != nil {
			message := fmt.Sprintf("AUTH0_DOMAIN misconfigured %s %v", cfg.Domain, err)
			w.Write([]byte(message))
		}

		authVars := fmt.Sprintf(`
			var AUTH0_CLIENT_ID='%s';
			var AUTH0_DOMAIN='%s';
			var AUTH0_CALLBACK_URL='%s';
			var LOGOUT_URL='%s'
			var AUTH0_LOGOUT_URL=AUTH0_DOMAIN+'/v2/logout?returnTo='+encodeURI(LOGOUT_URL)+'&client_id='+AUTH0_CLIENT_ID;
			`, cfg.ClientID, url.Hostname(), cfg.CallbackURL, cfg.LogoutURL)

		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(authVars))
	}
}
