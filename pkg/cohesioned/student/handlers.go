package student

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

type APIResponse struct {
	cohesioned.APIResponse
	Student *cohesioned.Student   `json:"student,omitempty"`
	List    []*cohesioned.Student `json:"list,omitempty"`
}

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &APIResponse{}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		list, err := repo.List(currentUser.ID)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred listing Student entities %v", err)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if list == nil {
			list = []*cohesioned.Student{}
		}

		resp.List = list
		r.JSON(w, http.StatusOK, resp)
	}
}

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &APIResponse{}
		resp.Student = &cohesioned.Student{}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&resp.Student); err != nil {
			resp.SetErrMsg("failed to unmarshall json %v", err)
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

		resp.Student.ParentID = currentUser.ID
		resp.Student.CreatedBy = currentUser.ID
		resp.Student.Created = time.Now()

		id, err := repo.Save(resp.Student)
		resp.Student.ID = id
		if err != nil {
			resp.SetErrMsg("Failed to save student %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &APIResponse{}
		resp.Student = &cohesioned.Student{}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		if err := decoder.Decode(&resp.Student); err != nil {
			resp.SetErrMsg("failed to unmarshall json %v", err)
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

		resp.Student.UpdatedBy = currentUser.ID
		resp.Student.Updated = time.Now()

		if err = repo.Update(resp.Student); err != nil {
			resp.SetErrMsg("Failed to save taxonomy %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		r.JSON(w, http.StatusOK, resp)
	}
}
