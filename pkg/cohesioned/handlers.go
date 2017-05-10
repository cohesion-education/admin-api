package cohesioned

import (
	"net/http"

	"github.com/unrolled/render"
)

func HomepageViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "homepage/index", nil)
	}
}

func UnauthorizedViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "401", nil, render.HTMLOptions{Layout: "empty-layout"})
	}
}

func ForbiddenViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "403", nil, render.HTMLOptions{Layout: "empty-layout"})
	}
}

func NotFoundViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "404", nil, render.HTMLOptions{Layout: "empty-layout"})
	}
}

func InternalServerErrorViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "500", nil, render.HTMLOptions{Layout: "empty-layout"})
	}
}
