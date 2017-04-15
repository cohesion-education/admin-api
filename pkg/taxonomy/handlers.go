package taxonomy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/common"
	"github.com/cohesion-education/admin-api/pkg/config"
)

func ListHandler(cfg *config.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repo := NewGCPDatastoreRepo(cfg.DatastoreClient)
		list, err := repo.List()
		if err != nil {
			cfg.Renderer.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		dashboard := common.NewDashboardViewWithProfile(req)
		dashboard.Set("list", list)
		cfg.Renderer.HTML(w, http.StatusOK, "taxonomy/list", dashboard)
		return
	}
}

func AddHandler(cfg *config.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		repo := NewGCPDatastoreRepo(cfg.DatastoreClient)
		decoder := json.NewDecoder(req.Body)

		var t Taxonomy

		if err := decoder.Decode(&t); err != nil {
			cfg.Renderer.Text(w, http.StatusInternalServerError, "failed to unmarshall json: "+err.Error())
			return
		}

		//TODO - validate taxonomy, and if fail, redirect back to form page with validation failure messages
		// dashboard := newDashboardWithProfile(req)
		// config.renderer.HTML(w, http.StatusOK, "taxonomy/add", dashboard)
		// return

		fmt.Println("adding taxonomy ", t)
		key, err := repo.Add(&t)
		if err != nil {
			cfg.Renderer.Text(w, http.StatusInternalServerError, "Failed to add Taxonomy: "+err.Error())
			return
		}

		// http.Redirect(w, req, "/taxonomy", http.StatusSeeOther)
		cfg.Renderer.JSON(w, http.StatusOK, key.ID)
	}
}
