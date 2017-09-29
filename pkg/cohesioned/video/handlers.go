package video

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type VideoResponse struct {
	*cohesioned.APIResponse
	*cohesioned.Video
	List []*cohesioned.Video `json:"list,omitempty"`
}

func NewAPIResponse(v *cohesioned.Video) *VideoResponse {
	return &VideoResponse{
		Video:       v,
		APIResponse: &cohesioned.APIResponse{},
	}
}

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoResponse{}
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
		resp := &VideoResponse{}

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

		video.CreatedBy = currentUser.ID
		id, err := repo.Save(video)
		video.ID = id
		if err != nil {
			resp.SetErrMsg("Failed to save video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = video

		r.JSON(w, http.StatusOK, resp)
	}
}

func UploadHandler(r *render.Render, repo Repo, cfg config.AwsConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		resp := NewAPIResponse(nil)

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

		video.StorageBucket = cfg.GetVideoBucket()
		video.StorageObjectName = fmt.Sprintf("%d-%s", video.ID, video.FileName)
		video.UpdatedBy = currentUser.ID
		video.Updated = time.Now()

		if _, err = repo.SetFile(req.Body, video); err != nil {
			resp.SetErrMsg("An unknown error occurred when saving the video file: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = video
		r.JSON(w, http.StatusOK, resp)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoResponse{}

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
		existing.Updated = time.Now()
		existing.UpdatedBy = currentUser.ID
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		if err := repo.Update(existing); err != nil {
			resp.SetErrMsg("Failed to update video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = existing
		r.JSON(w, http.StatusOK, resp)
	}
}

func GetByIDHandler(r *render.Render, repo Repo, cfg config.AwsConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &VideoResponse{}

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

		// signedURL, err := gcp.CreateSignedURL(video, cfg)
		signedURL, err := cfg.GetSignedURL(video.StorageBucket, video.StorageObjectName)
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
		resp := &VideoResponse{}

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
