package taxonomy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred listing Taxonomy entities %v", err)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if list == nil {
			list = []*cohesioned.Taxonomy{}
		}

		r.JSON(w, http.StatusOK, list)
	}
}

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		t := &cohesioned.Taxonomy{}

		if err := decoder.Decode(&t); err != nil {
			fmt.Printf("failed to unmarshall json %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		profile, ok := req.Context().Value(cohesioned.CurrentUserKey).(*cohesioned.Profile)
		if profile == nil {
			http.Redirect(w, req, "/401", http.StatusInternalServerError)
			return
		}

		if !ok {
			fmt.Printf("profile not of the proper type: %s\n", reflect.TypeOf(profile).String())
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		t.CreatedBy = profile.ID
		t.Created = time.Now()

		//TODO - validate taxonomy, and if fail, redirect back to form page with validation failure messages
		// dashboard := newDashboardWithProfile(req)
		// config.renderer.HTML(w, http.StatusOK, "taxonomy/add", dashboard)
		// return

		id, err := repo.Save(t)
		t.ID = id
		if err != nil {
			fmt.Printf("Failed to save taxonomy %v %v\n", t, err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		data := struct {
			ID int64 `json:"id"`
		}{
			t.ID,
		}

		r.JSON(w, http.StatusOK, data)
	}
}

func UpdateHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			apiResponse := &cohesioned.APIResponse{
				ErrMsg: fmt.Sprintf("%s is not a valid id %v", vars["id"], err),
			}
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusBadRequest, apiResponse)
			return
		}

		existing, err := repo.Get(id)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to find Taxonomy by ID: %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		if existing == nil {
			apiResponse := &cohesioned.APIResponse{
				ErrMsg: fmt.Sprintf("%s is not a valid id", id),
			}
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusNotFound, apiResponse)
			return
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		incoming := &cohesioned.Taxonomy{}
		if err := decoder.Decode(&incoming); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("failed to unmarshall json %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		currentUser, err := cohesioned.GetCurrentUser(req)
		if err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("An unexpected error occurred when trying to get the current user %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		existing.Name = incoming.Name
		existing.Children = incoming.Children
		existing.UpdatedBy = currentUser.ID
		existing.Updated = time.Now()

		if err = repo.Update(existing); err != nil {
			apiResponse := cohesioned.NewAPIErrorResponse("Failed to save taxonomy %v", err)
			fmt.Println(apiResponse.ErrMsg)
			r.JSON(w, http.StatusInternalServerError, apiResponse)
			return
		}

		r.JSON(w, http.StatusOK, existing)
	}
}

func ListChildrenHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		parentID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			data := struct {
				Error string `json:"error"`
			}{
				fmt.Sprintf("%s is not a valid id %v", vars["id"], err),
			}
			r.JSON(w, http.StatusNotFound, data)
			return
		}

		list, err := repo.ListChildren(parentID)
		if err != nil {
			data := struct {
				Error error `json:"error"`
			}{
				err,
			}
			r.JSON(w, http.StatusInternalServerError, data)
			return
		}

		data := struct {
			Children []*cohesioned.Taxonomy `json:"children"`
			ParentID int64                  `json:"parent_id"`
		}{
			list,
			parentID,
		}

		r.JSON(w, http.StatusOK, data)
	}
}

func FlatListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		list, err := repo.List()
		if err != nil {
			data := struct {
				Error error `json:"error"`
			}{
				err,
			}
			r.JSON(w, http.StatusInternalServerError, data)
			return
		}

		flattened := []*cohesioned.Taxonomy{}
		for _, t := range list {
			tFlattened, err := repo.Flatten(t)
			if err != nil {
				fmt.Printf("An unexpected error occured when trying to flatten %d %v", t.ID, err.Error())
				continue
			}
			flattened = append(flattened, tFlattened...)
		}

		r.JSON(w, http.StatusOK, flattened)
	}
}
