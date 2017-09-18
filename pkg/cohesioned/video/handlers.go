package video

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/gcp"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type VideoAPIResponse struct {
	cohesioned.APIResponse
	Video *cohesioned.Video   `json:"video,omitempty"`
	List  []*cohesioned.Video `json:"list,omitempty"`
}

func NewAPIResponse(v *cohesioned.Video) *VideoAPIResponse {
	resp := &VideoAPIResponse{
		Video: v,
	}
	resp.ID = v.ID()

	return resp
}

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoAPIResponse{}
		videos, err := repo.List()
		resp.List = videos
		if err != nil {
			resp.SetErrMsg("Failed to list videos %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoAPIResponse{}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		video := &cohesioned.Video{}

		if err := decoder.Decode(&video); err != nil {
			resp.SetErrMsg("Unable to process the Video payload. Error: %v\n", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if video.Validate() != true {
			resp.Video = video
			resp.SetErrMsg("Hmmm... you seem to be missing some required fields")
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		video.CreatedBy = currentUser
		video, err = repo.Add(video)
		if err != nil {
			resp.SetErrMsg("Failed to save video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.ID = video.ID()
		resp.Video = video

		r.JSON(w, http.StatusOK, resp)
	}
}

func UploadHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		resp := &VideoAPIResponse{}

		pathParams := mux.Vars(req)
		idParam := pathParams["id"]
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid video ID; %v", idParam, err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		video, err := repo.Get(id)
		if err != nil {
			resp.SetErrMsg("Unable to retrieve video by id %v - %v", idParam, err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if video == nil {
			resp.SetErrMsg("%v is not a valid video id", idParam)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		video.UpdatedBy = currentUser
		video, err = repo.SetFile(req.Body, video)
		if err != nil {
			resp.SetErrMsg("An unknown error occurred when saving the video file: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.ID = video.ID()
		resp.Video = video

		r.JSON(w, http.StatusOK, resp)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoAPIResponse{}

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		existing, err := repo.Get(videoID)
		if err != nil {
			resp.SetErrMsg("Failed to get video with id %d %v", videoID, err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&existing); err != nil {
			resp.SetErrMsg("Unable to process the Video payload. Error: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		if existing.Validate() != true {
			resp.Video = existing
			resp.SetErrMsg("Hmmm... you seem to be missing some required fields")
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		existing.UpdatedBy = currentUser
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		existing, err = repo.Update(existing)
		if err != nil {
			resp.SetErrMsg("Failed to update video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.ID = existing.ID()
		resp.Video = existing

		r.JSON(w, http.StatusOK, resp)
	}
}

func GetByIDHandler(r *render.Render, repo Repo, cfg *gcp.Config) http.HandlerFunc {
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

		video.SignedURL = signedURL
		resp.Video = video

		r.JSON(w, http.StatusOK, resp)
	}
}

func DeleteHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoAPIResponse{}

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		if err := repo.Delete(videoID); err != nil {
			resp.SetErrMsg("Failed to delete video with id %d %v", videoID, err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}
