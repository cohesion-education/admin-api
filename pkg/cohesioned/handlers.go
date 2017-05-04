package cohesioned

import (
	"net/http"

	"github.com/unrolled/render"
)

func HomepageViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "homepage/index", nil, render.HTMLOptions{Layout: "homepage/layout"})
	}
}
