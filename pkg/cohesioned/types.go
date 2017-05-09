package cohesioned

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"cloud.google.com/go/datastore"
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

type GCPPersisted struct {
	id  int64
	Key *datastore.Key `datastore:"__key__"`
}

func (p *GCPPersisted) ID() int64 {
	if p.Key != nil && p.Key.ID != 0 {
		return p.Key.ID
	}

	if p.id != 0 {
		return p.id
	}

	panic("Unable to load ID")
}

func (p *GCPPersisted) SetID(id int64) {
	fmt.Printf("p nil? %v\n", p)
	p.id = id
}

type Auditable struct {
	GCPPersisted
	Created   time.Time `datastore:"created" json:"created"`
	Updated   time.Time `datastore:"updated" json:"updated"`
	CreatedBy *Profile  `datastore:"created_by" json:"created_by"`
	UpdatedBy *Profile  `datastore:"updated_by" json:"updated_by"`
}

func (a *Auditable) SetCreated(t time.Time) {
	a.Created = t
}

func (a *Auditable) SetUpdated(t time.Time) {
	a.Updated = t
}

func (a *Auditable) SetCreatedBy(p *Profile) {
	a.CreatedBy = p
}

func (a *Auditable) SetUpdatedBy(p *Profile) {
	a.UpdatedBy = p
}
