package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
)

func ConfigHandler(cfg *config.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		domain := strings.TrimPrefix(cfg.Domain, "https://")

		authVars := fmt.Sprintf(`var AUTH0_CLIENT_ID='%s';
			var AUTH0_DOMAIN='%s';
			var AUTH0_CALLBACK_URL='%s';`,
			cfg.ClientID, domain, cfg.CallbackURL)

		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(authVars))
	}
}
