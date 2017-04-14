package main

import (
	"net/http"
)

func dashboardHandler(config *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := r.Context().Value(currentUserKey)
		config.renderer.HTML(w, http.StatusOK, "admin/dashboard", profile)
	}
}
