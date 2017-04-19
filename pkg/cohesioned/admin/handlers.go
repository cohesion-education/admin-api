package admin

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func DashboardViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard := cohesioned.NewDashboardViewWithProfile(req)
		r.HTML(w, http.StatusOK, "admin/dashboard", dashboard)
	}
}
