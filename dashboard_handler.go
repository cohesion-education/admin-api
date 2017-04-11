package main

import (
	"fmt"
	"net/http"
)

func dashboardHandler(config *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := r.Context().Value("current-user")
		fmt.Println("current user: ", profile)
		config.renderer.HTML(w, http.StatusOK, "admin/dashboard", profile)
	}
}
