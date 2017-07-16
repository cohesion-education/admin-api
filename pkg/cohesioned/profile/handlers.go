package profile

import (
	"fmt"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func SavePreferencesHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}
		if err = req.ParseForm(); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Error processing form %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p.Preferences.Newsletter = (req.PostFormValue("preferences.newsletter") == "on")
		p.Preferences.BetaProgram = (req.PostFormValue("preferences.betaprogram") == "on")

		if err = repo.Save(p); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to save User %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, p)
	}
}

func UpdatePreferencesHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p, err := repo.FindByEmail(currentUser.Email)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to retrieve your profile: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if err = req.ParseForm(); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Error processing form %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p.Preferences.Newsletter = (req.PostFormValue("preferences.newsletter") == "on")
		p.Preferences.BetaProgram = (req.PostFormValue("preferences.betaprogram") == "on")

		if err = repo.Update(p); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to save User %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, p)
	}
}

func GetCurrentUserHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p, err := repo.FindByEmail(currentUser.Email)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to lookup your user: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, p)
	}
}
