package config

import (
	"net/http"
	"os"

	"github.com/unrolled/render"
)

type configResponse struct {
	Auth0Domain   string `json:"auth0_domain"`
	Auth0ClientID string `json:"auth0_client_id"`
	CallbackURL   string `json:"callback_url"`
	GATrackingID  string `json:"ga_tracking_id"`
}

//Handler the config necessary to init the front-end
func Handler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &configResponse{
			Auth0Domain:   os.Getenv("AUTH0_DOMAIN"),
			Auth0ClientID: os.Getenv("AUTH0_CLIENT_ID"),
			CallbackURL:   os.Getenv("CALLBACK_URL"),
			GATrackingID:  os.Getenv("GA_TRACKING_ID"),
		}
		r.JSON(w, http.StatusOK, resp)
	}
}
