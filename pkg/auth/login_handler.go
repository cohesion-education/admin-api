package auth

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/config"
	"github.com/unrolled/render"
)

func LoginViewHandler(hc *config.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if profile := req.Context().Value(config.CurrentUserKey); profile != nil {
			http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
			return
		}

		hc.Renderer.HTML(w, http.StatusOK, "login/index", nil, render.HTMLOptions{Layout: ""})
	}
}
