package homepage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func HomepageViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h, err := repo.Get()
		if err != nil {
			fmt.Printf("Failed to retrieve homepage %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
		}

		if h == nil {
			h = cohesioned.NewHomepage(-1)
		}

		r.HTML(w, http.StatusOK, "homepage/index", h)
	}
}

func FormViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		h, err := repo.Get()
		if err != nil {
			fmt.Printf("Failed to retrieve homepage %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
		}

		if h == nil {
			h = cohesioned.NewHomepage(-1)
		}

		dashboard.Set("homepage", h)
		r.HTML(w, http.StatusOK, "homepage/form", dashboard)
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h, err := repo.Get()
		if err != nil {
			fmt.Printf("Failed to retrieve homepage %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
		}

		if h == nil {
			h = cohesioned.NewHomepage(-1)
		}

		h.HeaderTagline = req.FormValue("header_tagline")
		h.HeaderSubtext = req.FormValue("header_subtext")
		h.FeaturesHeaderTagline = req.FormValue("features_tagline")
		h.FeaturesHeaderSubtext = req.FormValue("features_subtext")

		if _, err := repo.Save(h); err != nil {
			fmt.Printf("Failed to save homepage %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
		}

		http.Redirect(w, req, "/admin/homepage", http.StatusSeeOther)
	}
}
