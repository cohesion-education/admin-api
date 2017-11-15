package student

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

type studentAPIResponse struct {
	*cohesioned.APIResponse
	Student *cohesioned.Student   `json:"student,omitempty"`
	List    []*cohesioned.Student `json:"students,omitempty"`
}

func newStudentAPIResponse() *studentAPIResponse {
	return &studentAPIResponse{
		APIResponse: &cohesioned.APIResponse{},
	}
}

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := newStudentAPIResponse()

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		list, err := repo.FindByUserID(currentUser.ID)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred listing Student entities %v", err)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if list == nil {
			list = []*cohesioned.Student{}
		}

		// resp.Profile = currentUser
		// currentUser.Students = list
		resp.List = list

		r.JSON(w, http.StatusOK, resp)
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := newStudentAPIResponse()
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

		//TODO - move this mess into a profile service...

		existingStudents, err := repo.FindByUserID(currentUser.ID)
		if err != nil {
			resp.SetErrMsg("An unexpected error occurred when trying to retrieve your current list of students: %v", err)
			fmt.Println(resp.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, resp)
			return
		}

		// resp.Profile = currentUser
		resp.List = make([]*cohesioned.Student, 0)

		// checks the existing student against the incoming student list for any students that have been removed and deletes those students
		for _, existingStudent := range existingStudents {
			if removedStudent := findStudent(existingStudent.Name, incomingList); removedStudent == nil {
				fmt.Printf("%s was not present in the incoming students list - deleting student", existingStudent.Name)
				if err := repo.Delete(existingStudent.ID); err != nil {
					resp.SetErrMsg("Failed to remove student %v", err)
					fmt.Println(resp.ErrMsg)
					r.JSON(w, http.StatusInternalServerError, resp)
					return
				}
			}
		}

		// checks incoming student list against existing students, and creates any newly added students or updates existing students
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
				resp.List = append(resp.List, incomingStudent)
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

				resp.List = append(resp.List, existingStudent)
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
