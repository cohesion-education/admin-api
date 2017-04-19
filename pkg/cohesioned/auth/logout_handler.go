package auth

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func LogoutHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   cohesioned.AuthSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
