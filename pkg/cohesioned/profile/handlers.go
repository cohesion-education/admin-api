package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
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

		if p == nil {
			p = currentUser
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		nextState := map[string]string{}

		if err = decoder.Decode(&nextState); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p.FullName = nextState["name"]
		p.Email = nextState["email"]
		p.State = nextState["state"]
		p.County = nextState["county"]

		if err = repo.Save(p); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to save User %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, p)
	}
}

func SavePreferencesHandler(r *render.Render, repo Repo) http.HandlerFunc {
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

		if p == nil {
			p = currentUser
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		preferences := map[string]bool{}

		if err = decoder.Decode(&preferences); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		p.Preferences.Newsletter = preferences["newsletter"]
		p.Preferences.BetaProgram = preferences["beta_program"]

		if err = repo.Save(p); err != nil {
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
