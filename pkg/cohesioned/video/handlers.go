package video

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to list videos "+err.Error())
			return
		}
		dashboard := common.NewDashboardViewWithProfile(req)
		dashboard.Set("list", list)
		r.HTML(w, http.StatusOK, "video/list", dashboard)
		return
	}
}

func FormHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dashboard := common.NewDashboardViewWithProfile(req)
		r.HTML(w, http.StatusOK, "video/form", dashboard)
		return
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		profile, ok := req.Context().Value(config.CurrentUserKey).(*cohesioned.Profile)
		if profile == nil {
			r.Text(w, http.StatusInternalServerError, "middleware did not set profile in the context as expected")
			return
		}

		if !ok {
			errMsg := fmt.Sprintf("profile not of the proper type: %s", reflect.TypeOf(profile).String())
			r.Text(w, http.StatusInternalServerError, errMsg)
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

		contentType := fileHeader.Header.Get("Content-Type")
		fmt.Println("content type ", contentType)

		title := req.FormValue("title")
		fileName := fileHeader.Filename

		video := &cohesioned.Video{
			Title:    title,
			FileName: fileName,
			//TODO - inject this as config
			StorageBucket: "cohesion-dev",
			CreatedBy:     profile,
		}

		if _, err := repo.Add(file, video); err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to save video "+err.Error())
			return
		}

		//TODO - provide a better response since this is an API call
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

		dashboard := common.NewDashboardViewWithProfile(req)
		dashboard.Set("video", video)

		r.HTML(w, http.StatusOK, "video/show", dashboard)
		return
	}
}

func StreamHandler(r *render.Render, repo Repo) http.HandlerFunc {
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

		signedURL, err := video.SignedURL()
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
		return
	}
}
