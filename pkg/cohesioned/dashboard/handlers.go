package dashboard

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func AdminViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		r.HTML(w, http.StatusOK, "dashboard/admin", d)
	}
}

func UserViewHandler(r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		profile, err := d.Profile()
		if err != nil {
			fmt.Printf("error getting current user from request context %v\n", err)
			http.Redirect(w, req, "/401", http.StatusSeeOther)
			return
		}

		if profile.HasRole("beta-tester") {
			//TODO - need newer, more limited menu options
			r.HTML(w, http.StatusOK, "dashboard/beta-tester", d)
		}

		if os.Getenv("LAUNCHED") != "true" {
			r.HTML(w, http.StatusOK, "dashboard/early-reg", d)
			return
		}

		r.HTML(w, http.StatusOK, "dashboard/user", d)
	}
}
