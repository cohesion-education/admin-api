package dashboard

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func AdminViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("dashboard view handler being hit")
		d, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusInternalServerError)
			return
		}

		r.HTML(w, http.StatusOK, "dashboard/admin", d)
	}
}

func UserViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("dashboard view handler being hit")
		d, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusInternalServerError)
			return
		}

		r.HTML(w, http.StatusOK, "dashboard/user", d)
	}
}
