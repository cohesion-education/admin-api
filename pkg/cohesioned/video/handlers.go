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

type VideoAPIResponse struct {
	cohesioned.APIResponse
	Video *cohesioned.Video   `json:"video,omitempty"`
	List  []*cohesioned.Video `json:"list,omitempty"`
}

func ListViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
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

func FormViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}

		vars := mux.Vars(req)
		idParam := vars["id"]
		if len(idParam) > 0 {
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

			dashboard.Set("video", video)
		}

		r.HTML(w, http.StatusOK, "video/form", dashboard)
		return
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		video := &cohesioned.Video{}
		resp := &VideoAPIResponse{}

		profile, err := cohesioned.GetProfile(req)
		video.CreatedBy = profile
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			resp.SetErrMsg("Unexpected error %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if err = req.ParseMultipartForm(1000 * 10); err != nil {
			resp.SetErrMsg("Failed to parse form %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		file, fileHeader, err := req.FormFile("video_file")
		if err == http.ErrMissingFile {
			resp.AddValidationError("video_file", "You must select a file for upload")
		} else if err != nil {
			resp.SetErrMsg("Failed to get form file %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if file != nil && fileHeader != nil {
			video.FileName = fileHeader.Filename
			video.FileReader = file
		}

		video.Title = req.FormValue("title")
		if len(video.Title) == 0 {
			resp.AddValidationError("title", "Title is required")
		}

		taxonomyID, err := strconv.ParseInt(req.FormValue("taxonomy_id"), 10, 64)
		video.TaxonomyID = taxonomyID
		if err != nil {
			resp.AddValidationError("taxonomy_id", "Invalid Taxonomy ID "+err.Error())
		}

		if len(resp.ValidationErrors) > 0 {
			resp.SetErrMsg("Hmmm... you seem to be missing some required fields")
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		video, err = repo.Add(file, video)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to save video "+err.Error())
			return
		}

		resp.RedirectURL = fmt.Sprintf("/admin/video/%d", video.ID())
		r.JSON(w, http.StatusSeeOther, resp)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoAPIResponse{}

		profile, err := cohesioned.GetProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			resp.SetErrMsg("Unexpected error %v")
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if err = req.ParseMultipartForm(1000 * 10); err != nil {
			resp.SetErrMsg("Failed to parse form " + err.Error())
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		file, fileHeader, err := req.FormFile("video_file")
		if err != nil && err != http.ErrMissingFile {
			resp.SetErrMsg("Failed to get form file " + err.Error())
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		id, err := strconv.ParseInt(req.FormValue("id"), 10, 64)
		if err != nil {
			resp.SetErrMsg("Invalid video id %s %v", req.FormValue("id"), err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		taxonomyID, err := strconv.ParseInt(req.FormValue("taxonomy_id"), 10, 64)
		if err != nil {
			resp.SetErrMsg("Invalid taxonomy id %s %v", req.FormValue("taxonomy_id"), err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		video, err := repo.Get(id)
		if err != nil {
			resp.SetErr(err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		video.Title = req.FormValue("title")
		video.TaxonomyID = taxonomyID
		video.UpdatedBy = profile

		if fileHeader != nil {
			video.FileName = fileHeader.Filename
		}

		video, err = repo.Update(file, video)
		if err != nil {
			resp.SetErrMsg("Failed to update video %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.RedirectURL = fmt.Sprintf("/admin/video/%d", video.ID())
		r.JSON(w, http.StatusSeeOther, resp)
	}
}

func ShowViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
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
		resp := &VideoAPIResponse{}

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		video, err := repo.Get(videoID)
		if err != nil {
			resp.SetErrMsg("Failed to get video with id %d %v", videoID, err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		signedURL, err := gcp.CreateSignedURL(video, cfg)
		if err != nil {
			resp.SetErrMsg("Failed to generate signed url %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		http.Redirect(w, req, signedURL, http.StatusSeeOther)
	}
}
