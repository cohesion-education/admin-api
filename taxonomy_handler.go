package main

import (
	"fmt"
	"net/http"

	"github.com/unrolled/render"
)

func taxonomyListHandler(r *render.Render, repo *taxonomyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if repo != nil {
			list, err := repo.list()
			if err != nil {
				r.Text(w, http.StatusInternalServerError, err.Error())
				return
			}

			r.HTML(w, http.StatusOK, "ubold/taxonomy-list", list)
			return
		}

		r.Text(w, http.StatusOK, fmt.Sprintf("Taxonomy Features Coming soon!"))
		return
	}
}
