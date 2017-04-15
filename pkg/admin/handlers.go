package admin

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/common"
	"github.com/cohesion-education/admin-api/pkg/config"
)

func DashboardViewHandler(cfg *config.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard := common.NewDashboardViewWithProfile(req)
		cfg.Renderer.HTML(w, http.StatusOK, "admin/dashboard", dashboard)
	}
}
