package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

func userListHandler(r *render.Render, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Text(w, http.StatusOK, fmt.Sprintf("User Page Coming soon!"))
		return
	}
}
