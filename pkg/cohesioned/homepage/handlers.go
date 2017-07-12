package homepage

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func HomepageHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h, err := repo.Get()
		if err != nil {
			resp := cohesioned.NewAPIErrorResponse("Failed to retrieve homepage %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
		}

		if h == nil {
			h = cohesioned.NewHomepage(-1)
			h.Header.Title = "Test Title"
			h.Header.Subtitle = "Test Subtitle"
		}

		r.JSON(w, http.StatusOK, h)
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

		h.Header.Title = req.FormValue("header_tagline")
		h.Header.Subtitle = req.FormValue("header_subtext")
		h.Features.Title = req.FormValue("features_tagline")
		h.Features.Subtitle = req.FormValue("features_subtext")

		if _, err := repo.Save(h); err != nil {
			fmt.Printf("Failed to save homepage %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
		}

		http.Redirect(w, req, "/admin/homepage", http.StatusSeeOther)
	}
}
