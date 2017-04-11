package main

import (
	"fmt"
	"net/http"
)

func logoutHandler(config *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("logging out... goodbye")
		cookie := &http.Cookie{
			Name:   authSessionCookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
