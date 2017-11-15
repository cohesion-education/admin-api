package report

import (
	"fmt"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/billing"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
	"github.com/cohesion-education/api/pkg/cohesioned/student"
	"github.com/unrolled/render"
)

func GetUserList(r *render.Render, repo profile.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to list users: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, list)
	}
}

func GetStudentList(r *render.Render, repo student.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to list students: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, list)
	}
}

func GetPaymentDetailList(r *render.Render, repo billing.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to list payment details: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, list)
	}
}
