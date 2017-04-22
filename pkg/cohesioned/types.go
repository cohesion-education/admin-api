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

type ValidationError struct {
	Field string `json:"field_name,omitempty"`
	Err   string `json:"error,omitempty"`
}

type APIResponse struct {
	ID               int64              `json:"id,omitempty"`
	ErrMsg           string             `json:"error,omitempty"`
	RedirectURL      string             `json:"redirect_url,omitempty"`
	ValidationErrors []*ValidationError `json:"validation_errors,omitempty"`
}

func (r *APIResponse) AddValidationError(field, err string) *APIResponse {
	r.ValidationErrors = append(r.ValidationErrors, &ValidationError{
		Field: field,
		Err:   err,
	})

	return r
}

func (r *APIResponse) SetErr(err error) {
	r.ErrMsg = err.Error()
}

func (r *APIResponse) SetErrMsg(format string, a ...interface{}) {
	r.ErrMsg = fmt.Sprintf(format, a...)
}
