package taxonomy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
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
