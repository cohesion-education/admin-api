package main

import (
	"net/http"
)

func taxonomyListHandler(config *handlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taxonomyRepo := &taxonomyRepository{session: config.mongoSession}
		list, err := taxonomyRepo.list()
		if err != nil {
			config.renderer.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		config.renderer.HTML(w, http.StatusOK, "taxonomy/list", list)
		return
	}
}
