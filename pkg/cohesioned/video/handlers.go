package video

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/gcp"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			//TODO - direct users to an official Error page
			r.Text(w, http.StatusInternalServerError, "Failed to list videos "+err.Error())
			return
		}

		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		dashboard.Set("list", list)
		r.HTML(w, http.StatusOK, "video/list", dashboard)
		return
	}
}

func FormHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		r.HTML(w, http.StatusOK, "video/form", dashboard)
		return
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		profile, err := cohesioned.GetProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		if err := req.ParseMultipartForm(1000 * 10); err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to parse form "+err.Error())
			return
		}

		file, fileHeader, err := req.FormFile("video_file")
		if err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to get form file "+err.Error())
			return
		}

		title := req.FormValue("title")
		fileName := fileHeader.Filename
		taxonomyID, err := strconv.ParseInt(req.FormValue("taxonomy_id"), 10, 64)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Invalid taxonomy id %s %v", req.FormValue("taxonomy_id"), err))
			return
		}

		video := &cohesioned.Video{
			Title:      title,
			FileName:   fileName,
			CreatedBy:  profile,
			TaxonomyID: taxonomyID,
		}

		video, err = repo.Add(file, video)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to save video "+err.Error())
			return
		}

		http.Redirect(w, req, fmt.Sprintf("/admin/video/%d", video.ID()), http.StatusSeeOther)
	}
}

func ShowHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("%s is not a valid id %v", vars["id"], err))
			return
		}

		video, err := repo.Get(videoID)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		dashboard.Set("video", video)
		r.HTML(w, http.StatusOK, "video/show", dashboard)
		return
	}
}

func StreamHandler(r *render.Render, repo Repo, cfg *gcp.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			data := struct {
				Error string `json:"error"`
			}{
				fmt.Sprintf("%s is not a valid id %v", vars["id"], err),
			}
			r.JSON(w, http.StatusNotFound, data)
			return
		}

		video, err := repo.Get(videoID)
		if err != nil {
			data := struct {
				Error string `json:"error"`
			}{
				fmt.Sprintf("Failed to get video with id %d %v", videoID, err),
			}
			r.JSON(w, http.StatusInternalServerError, data)
			return
		}

		signedURL, err := gcp.CreateSignedURL(video, cfg)
		if err != nil {
			data := struct {
				Error string `json:"error"`
			}{
				fmt.Sprintf("Failed to generate signed url %v", err),
			}
			r.JSON(w, http.StatusInternalServerError, data)
			return
		}

		http.Redirect(w, req, signedURL, http.StatusSeeOther)
	}
}
