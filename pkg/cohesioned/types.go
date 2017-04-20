package cohesioned

import (
	"fmt"
	"net/http"
	"reflect"
)

const profileKey = "profile"

type DashboardView map[string]interface{}

func NewDashboardViewWithProfile(req *http.Request) (*DashboardView, error) {
	d := &DashboardView{}
	profile, err := GetProfile(req)
	if err != nil {
		return nil, err
	}
	d.SetProfile(profile)
	return d, nil
}

func (d DashboardView) Profile() (*Profile, error) {
	profile, ok := d.Get(profileKey).(*Profile)
	if !ok {
		return nil, fmt.Errorf("profile was not of type *cohesion.Profile was %s", reflect.TypeOf(profile).String())
	}

	return profile, nil
}

func (d DashboardView) SetProfile(p *Profile) {
	d.Set(profileKey, p)
}

func (d DashboardView) Set(key string, value interface{}) {
	d[key] = value
}

func (d DashboardView) Get(key string) interface{} {
	return d[key]
}
