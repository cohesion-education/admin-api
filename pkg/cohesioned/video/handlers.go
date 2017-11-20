package video

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type VideoResponse struct {
	*cohesioned.APIResponse
	*cohesioned.Video
	List    []*cohesioned.Video            `json:"list,omitempty"`
	ByGrade map[string][]*cohesioned.Video `json:"by_grade,omitempty"`
}

func NewAPIResponse(v *cohesioned.Video) *VideoResponse {
	return &VideoResponse{
		Video:       v,
		APIResponse: &cohesioned.APIResponse{},
		ByGrade:     make(map[string][]*cohesioned.Video),
	}
}

func ListHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)
		videos, err := svc.List()
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

func FindByGradeHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

		pathParams := mux.Vars(req)
		grade := pathParams["grade"]

		videosByGrade, err := svc.FindByGrade(grade)
		resp.ByGrade = videosByGrade
		if err != nil {
			resp.SetErrMsg("Failed to list videos by grade %s: %v", grade, err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}

func FindByTaxonomyHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

		pathParams := mux.Vars(req)
		taxonomyIDParam := pathParams["taxonomy_id"]
		taxonomyID, err := strconv.ParseInt(taxonomyIDParam, 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid taxonomy ID; %v", taxonomyIDParam, err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusBadRequest, resp)
			return
		}

		videos, err := svc.FindByTaxonomyID(taxonomyID)
		resp.List = videos
		if err != nil {
			resp.SetErrMsg("Failed to list videos by taxonomy id %d: %v", taxonomyID, err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}

func AddHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

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

		ctx := req.Context()
		if err := svc.Save(ctx, video); err != nil {
			resp.SetErrMsg("Failed to save video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = video
		r.JSON(w, http.StatusOK, resp)
	}
}

func UploadHandler(r *render.Render, svc AdminService) http.HandlerFunc {
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

		video, err := svc.Get(id)
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

		ctx := req.Context()
		if err := svc.SetFile(ctx, req.Body, video); err != nil {
			resp.SetErrMsg("An unknown error occurred when saving the video file: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = video
		r.JSON(w, http.StatusOK, resp)
	}
}

func UpdateHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		existing, err := svc.Get(videoID)
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

		ctx := req.Context()
		if err := svc.Update(ctx, existing); err != nil {
			resp.SetErrMsg("Failed to update video %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = existing
		r.JSON(w, http.StatusOK, resp)
	}
}

func GetByIDHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		video, err := svc.GetWithSignedURL(videoID)
		if err != nil {
			resp.SetErrMsg("Failed to get video with id %d %v", videoID, err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Video = video

		r.JSON(w, http.StatusOK, resp)
	}
}

func DeleteHandler(r *render.Render, svc AdminService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := NewAPIResponse(nil)

		vars := mux.Vars(req)
		videoID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			resp.SetErrMsg("%s is not a valid id %v", vars["id"], err)
			r.JSON(w, http.StatusNotFound, resp)
			return
		}

		if err := svc.Delete(videoID); err != nil {
			resp.SetErrMsg("Failed to delete video with id %d %v", videoID, err)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}
