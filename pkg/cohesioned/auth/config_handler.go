package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
)

func ConfigHandler(cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		domain := strings.TrimPrefix(cfg.Domain, "https://")

		authVars := fmt.Sprintf(`
			var AUTH0_CLIENT_ID='%s';
			var AUTH0_DOMAIN='%s';
			var AUTH0_CALLBACK_URL='%s';
			var LOGOUT_URL='%s'
			var AUTH0_LOGOUT_URL=AUTH0_DOMAIN+'/v2/logout?returnTo='+encodeURI(LOGOUT_URL)+'&client_id='+AUTH0_CLIENT_ID;
			`, cfg.ClientID, domain, cfg.CallbackURL, cfg.LogoutURL)

		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(authVars))
	}
}
