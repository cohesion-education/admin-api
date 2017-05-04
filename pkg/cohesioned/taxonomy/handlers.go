package taxonomy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ListViewHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		list, err := repo.List()
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			//TODO - direct users to an official Error page
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			r.Text(w, http.StatusInternalServerError, fmt.Sprintf("Unexpected error %v", err))
			return
		}
		dashboard.Set("list", list)
		r.HTML(w, http.StatusOK, "taxonomy/list", dashboard)
		return
	}
}

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
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

		r.JSON(w, http.StatusOK, list)
	}
}

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		t := &cohesioned.Taxonomy{}

		if err := decoder.Decode(&t); err != nil {
			r.Text(w, http.StatusInternalServerError, "failed to unmarshall json: "+err.Error())
			return
		}

		profile, ok := req.Context().Value(cohesioned.CurrentUserKey).(*cohesioned.Profile)
		if profile == nil {
			r.Text(w, http.StatusInternalServerError, "middleware did not set profile in the context as expected")
			return
		}

		if !ok {
			errMsg := fmt.Sprintf("profile not of the proper type: %s", reflect.TypeOf(profile).String())
			r.Text(w, http.StatusInternalServerError, errMsg)
			return
		}

		t.CreatedBy = profile

		//TODO - validate taxonomy, and if fail, redirect back to form page with validation failure messages
		// dashboard := newDashboardWithProfile(req)
		// config.renderer.HTML(w, http.StatusOK, "taxonomy/add", dashboard)
		// return

		t, err := repo.Add(t)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		data := struct {
			ID int64 `json:"id"`
		}{
			t.ID(),
		}

		r.JSON(w, http.StatusOK, data)
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
				fmt.Printf("An unexpected error occured when trying to flatten %d %v", t.ID(), err.Error())
				continue
			}
			flattened = append(flattened, tFlattened...)
		}

		r.JSON(w, http.StatusOK, flattened)
	}
}
