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
			fmt.Printf("Failed to list taxonomy items %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
			return
		}

		dashboard, err := cohesioned.NewDashboardViewWithProfile(req)
		if err != nil {
			log.Printf("Unexpected error when trying to get dashboard view with profile %v\n", err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
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

		t.CreatedBy = profile

		//TODO - validate taxonomy, and if fail, redirect back to form page with validation failure messages
		// dashboard := newDashboardWithProfile(req)
		// config.renderer.HTML(w, http.StatusOK, "taxonomy/add", dashboard)
		// return

		t, err := repo.Add(t)
		if err != nil {
			fmt.Printf("Failed to save taxonomy %v %v\n", t, err)
			http.Redirect(w, req, "/500", http.StatusSeeOther)
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
