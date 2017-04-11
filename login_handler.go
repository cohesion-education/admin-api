package main

import (
	"net/http"

	"github.com/unrolled/render"
)

func loginViewHandler(ac *authConfig, hc *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := ac.getCurrentSession(req)
		if err != nil {
			hc.renderer.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		if profile := session.Values["profile"]; profile != nil {
			http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
			return
		}

		hc.renderer.HTML(w, http.StatusOK, "ubold/login", nil, render.HTMLOptions{Layout: ""})
	}
}
