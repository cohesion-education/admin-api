package taxonomy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ListHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		list, err := repo.List()
		if err != nil {
			r.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		dashboard := common.NewDashboardViewWithProfile(req)
		dashboard.Set("list", list)
		r.HTML(w, http.StatusOK, "taxonomy/list", dashboard)
		return
	}
}

func AddHandler(r *render.Render, repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)

		var t cohesioned.Taxonomy

		if err := decoder.Decode(&t); err != nil {
			r.Text(w, http.StatusInternalServerError, "failed to unmarshall json: "+err.Error())
			return
		}

		profile, ok := req.Context().Value(config.CurrentUserKey).(*cohesioned.Profile)
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

		fmt.Println("adding taxonomy ", t)
		key, err := repo.Add(&t)
		if err != nil {
			r.Text(w, http.StatusInternalServerError, "Failed to add Taxonomy: "+err.Error())
			return
		}

		// http.Redirect(w, req, "/taxonomy", http.StatusSeeOther)
		r.JSON(w, http.StatusOK, key.ID)
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
