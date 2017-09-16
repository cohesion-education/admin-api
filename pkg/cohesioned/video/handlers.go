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

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		video := &cohesioned.Video{}
		resp := &VideoAPIResponse{}

		if err := decoder.Decode(&video); err != nil {
			resp.SetErrMsg("Unable to process the Video payload. Error: %v\n", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		//TODO - integrate this into the Video struct for consistency
		if len(video.Title) == 0 {
			resp.AddValidationError("title", "Video Title is required")
		}

		if video.TaxonomyID == -1 {
			resp.AddValidationError("taxonomy_id", "Video Taxonomy ID is required")
		}

		if len(resp.ValidationErrors) > 0 {
			resp.SetErrMsg("Hmmm... you seem to be missing some required fields")
			r.JSON(w, http.StatusBadRequest, resp)
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

		fmt.Printf("Headers: %v\n", req.Header)

		video.FileName = req.Header.Get("file-name")
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
		// resp := &VideoAPIResponse{}
		//
		// profile, err := cohesioned.GetCurrentUser(req)
		// if err != nil {
		// 	log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
		// 	resp.SetErrMsg("Unexpected error %v")
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// if err = req.ParseMultipartForm(1000 * 10); err != nil {
		// 	resp.SetErrMsg("Failed to parse form " + err.Error())
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// file, fileHeader, err := req.FormFile("video_file")
		// if err != nil && err != http.ErrMissingFile {
		// 	resp.SetErrMsg("Failed to get form file " + err.Error())
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// id, err := strconv.ParseInt(req.FormValue("id"), 10, 64)
		// if err != nil {
		// 	resp.SetErrMsg("Invalid video id %s %v", req.FormValue("id"), err)
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// taxonomyID, err := strconv.ParseInt(req.FormValue("taxonomy_id"), 10, 64)
		// if err != nil {
		// 	resp.SetErrMsg("Invalid taxonomy id %s %v", req.FormValue("taxonomy_id"), err)
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// video, err := repo.Get(id)
		// if err != nil {
		// 	resp.SetErr(err)
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// video.Title = req.FormValue("title")
		// video.TaxonomyID = taxonomyID
		// video.UpdatedBy = profile
		//
		// if fileHeader != nil {
		// 	video.FileName = fileHeader.Filename
		// }
		//
		// video, err = repo.Update(file, video)
		// if err != nil {
		// 	resp.SetErrMsg("Failed to update video %v", err)
		// 	r.JSON(w, http.StatusInternalServerError, resp)
		// 	return
		// }
		//
		// resp.RedirectURL = fmt.Sprintf("/admin/video/%d", video.ID())
		// r.JSON(w, http.StatusOK, resp)
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
