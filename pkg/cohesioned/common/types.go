package common

import (
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
)

type DashboardView map[string]interface{}

func NewDashboardViewWithProfile(req *http.Request) *DashboardView {
	d := &DashboardView{}
	profile := req.Context().Value(config.CurrentUserKey)
	d.Set("profile", profile)
	return d
}

func (d DashboardView) Set(key string, value interface{}) {
	d[key] = value
}
