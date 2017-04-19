package cohesioned

import (
	"net/http"
)

type DashboardView map[string]interface{}

func NewDashboardViewWithProfile(req *http.Request) *DashboardView {
	d := &DashboardView{}
	profile := req.Context().Value(CurrentUserKey)
	d.Set("profile", profile)
	return d
}

func (d DashboardView) Set(key string, value interface{}) {
	d[key] = value
}

func (d DashboardView) Get(key string) interface{} {
	return d[key]
}
