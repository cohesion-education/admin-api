package admin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func DashboardViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("dashboard view handler being hit")
		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		fmt.Printf("dashboard: %v\n", dashboard)

		r.HTML(w, http.StatusOK, "admin/dashboard", dashboard)
	}
}
