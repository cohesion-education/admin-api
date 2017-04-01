package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

func logoutHandler(r *render.Render, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("logging out... goodbye")
		cookie := &http.Cookie{
			Name:   "auth-session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
