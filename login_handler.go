package main

import (
	"net/http"

	"github.com/unrolled/render"
)

func loginViewHandler(hc *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if profile := req.Context().Value(currentUserKey); profile != nil {
			http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
			return
		}

		hc.renderer.HTML(w, http.StatusOK, "login/index", nil, render.HTMLOptions{Layout: ""})
	}
}
