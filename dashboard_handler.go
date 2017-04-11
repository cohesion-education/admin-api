package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

func dashboardHandler(r *render.Render, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "auth-session")
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		profile := session.Values["profile"]
		fmt.Printf("profile: %v", profile)

		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Accel-Expires", "0")

		r.HTML(w, http.StatusOK, "admin/dashboard", profile)
	}
}
