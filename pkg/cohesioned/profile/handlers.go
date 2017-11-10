package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/unrolled/render"
)

func GetOrCreateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		incoming := &cohesioned.Profile{
			Created:       time.Now(),
			TrialStart:    time.Now(),
			BillingStatus: cohesioned.BillingStatusTrial,
		}

		if err := decoder.Decode(&incoming); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		existing, err := repo.FindByEmail(incoming.Email)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to find user by email %s: %v", incoming.Email, err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if existing != nil {
			fmt.Println("returning existing user ", existing.Email)
			r.JSON(w, http.StatusOK, existing)
			return
		}

		id, err := repo.Save(incoming)
		incoming.ID = id
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to create user %s: %v", incoming.Email, err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, incoming)
	}
}

func SaveHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		existing, err := repo.FindByEmail(currentUser.Email)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if existing != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("User with that email address already exists: %s", existing.Email)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		incoming := &cohesioned.Profile{}
		if err = decoder.Decode(&incoming); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		incoming.Created = time.Now()

		id, err := repo.Save(existing)
		incoming.ID = id
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to save User %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, incoming)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		existing, err := repo.FindByEmail(currentUser.Email)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to retrieve your profile: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		incoming := &cohesioned.Profile{}
		if err = decoder.Decode(&incoming); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		existing.FullName = incoming.FullName
		existing.Email = incoming.Email
		existing.State = incoming.State
		existing.County = incoming.County
		existing.Onboarded = incoming.Onboarded
		// existing.Students = incoming.Students
		existing.Updated = time.Now()

		if repo.Update(existing); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to update User %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, existing)
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
			apiResponse := cohesioned.NewAPIErrorResponse("We were unable to find your user", nil)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
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

		p.Updated = time.Now()
		p.Preferences.Newsletter = preferences["newsletter"]
		p.Preferences.BetaProgram = preferences["beta_program"]

		if err := repo.Update(p); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to update User %v", err)
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
