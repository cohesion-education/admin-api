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

type apiResponse struct {
	ID          int64               `json:"id,omitempty"`
	ErrMsg      string              `json:"error,omitempty"`
	Video       *cohesioned.Video   `json:"video,omitempty"`
	List        []*cohesioned.Video `json:"list,omitempty"`
	RedirectURL string              `json:"redirect_url,omitempty"`
}

func (r *apiResponse) setErr(err error) {
	r.ErrMsg = err.Error()
}

func (r *apiResponse) setErrMsg(format string, a ...interface{}) {
	r.ErrMsg = fmt.Sprintf(format, a...)
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

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &apiResponse{}

		vars := mux.Vars(req)

		profile, err := cohesioned.GetProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			resp.setErrMsg("Unexpected error %v")
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if err = req.ParseMultipartForm(1000 * 10); err != nil {
			resp.setErrMsg("Failed to parse form " + err.Error())
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		file, fileHeader, err := req.FormFile("video_file")
		if err != nil && err != http.ErrMissingFile {
			resp.setErrMsg("Failed to get form file " + err.Error())
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.setErrMsg("Invalid video id %s %v", req.FormValue("id"), err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		taxonomyID, err := strconv.ParseInt(req.FormValue("taxonomy_id"), 10, 64)
		if err != nil {
			resp.setErrMsg("Invalid taxonomy id %s %v", req.FormValue("taxonomy_id"), err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		video, err := repo.Get(id)
		if err != nil {
			resp.setErr(err)
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
			resp.setErrMsg("Failed to update video %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.RedirectURL = fmt.Sprintf("/admin/video/%d", video.ID())
		r.JSON(w, http.StatusOK, resp)
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
		resp := &apiResponse{}

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.setErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		video, err := repo.Get(videoID)
		if err != nil {
			resp.setErrMsg("Failed to get video with id %d %v", videoID, err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		signedURL, err := gcp.CreateSignedURL(video, cfg)
		if err != nil {
			resp.setErrMsg("Failed to generate signed url %v", err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		http.Redirect(w, req, signedURL, http.StatusSeeOther)
	}
}
