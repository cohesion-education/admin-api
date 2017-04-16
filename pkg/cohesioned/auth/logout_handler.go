package auth

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/unrolled/render"
)

func LogoutHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   config.AuthSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
