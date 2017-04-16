package auth

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/unrolled/render"
)

func LoginViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if profile := req.Context().Value(config.CurrentUserKey); profile != nil {
			http.Redirect(w, req, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		r.HTML(w, http.StatusOK, "login/index", nil, render.HTMLOptions{Layout: ""})
	}
}
