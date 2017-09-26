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
	Profile *cohesioned.Profile   `json:"profile,omitempty"`
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

		resp.Profile = currentUser
		currentUser.Students = list

		r.JSON(w, http.StatusOK, resp)
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := &APIResponse{}
		incomingList := make([]*cohesioned.Student, 0)

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&incomingList); err != nil {
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

		existingStudents, err := repo.List(currentUser.ID)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to retrieve your current list of students: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp.Profile = currentUser
		currentUser.Students = make([]*cohesioned.Student, len(incomingList))

		for _, incomingStudent := range incomingList {
			if existingStudent := findStudent(incomingStudent.Name, existingStudents); existingStudent == nil {
				fmt.Printf("incoming student does not exist - creating %v\n", incomingStudent)
				incomingStudent.ParentID = currentUser.ID
				incomingStudent.CreatedBy = currentUser.ID
				incomingStudent.Created = time.Now()

				id, err := repo.Save(incomingStudent)
				if err != nil {
					resp.SetErrMsg("Failed to save student %v", err)
					fmt.Println(resp.ErrMsg)
					r.JSON(w, http.StatusInternalServerError, resp)
					return
				}

				incomingStudent.ID = id
				currentUser.Students = append(currentUser.Students, incomingStudent)
			} else {
				fmt.Printf("incoming student exists - updating %v\n", existingStudent)
				existingStudent.Name = incomingStudent.Name
				existingStudent.Grade = incomingStudent.Grade
				existingStudent.School = incomingStudent.School
				existingStudent.UpdatedBy = currentUser.ID
				existingStudent.Updated = time.Now()

				if err := repo.Update(existingStudent); err != nil {
					resp.SetErrMsg("Failed to save student %v", err)
					fmt.Println(resp.ErrMsg)
					r.JSON(w, http.StatusInternalServerError, resp)
					return
				}

				currentUser.Students = append(currentUser.Students, existingStudent)
			}
		}

		r.JSON(w, http.StatusOK, resp)
	}
}

func findStudent(name string, students []*cohesioned.Student) *cohesioned.Student {
	for _, existingStudent := range students {
		if existingStudent.Name == name {
			return existingStudent
		}
	}

	return nil
}
