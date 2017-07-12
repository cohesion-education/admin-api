package profile

import (
	"fmt"
	"net/http"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/unrolled/render"
)

func SavePreferencesHandler(r *render.Render, cfg *config.AuthConfig, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p, err := cohesioned.GetProfile(req)
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

		session, err := cfg.GetCurrentSession(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("unable to Get Current Session %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		session.Values[cohesioned.CurrentUserSessionKey] = p

		if err := session.Save(req, w); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to save Session %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, p)
	}
}
