package fakes

import (
	"bytes"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var (
	FakeRenderer = render.New(render.Options{
		Layout: "layout",
		RenderPartialsWithoutPrefix: true,
		Directory:                   "../../../templates",
	})

	FakeAuthConfig = &config.AuthConfig{
		SessionStore: sessions.NewFilesystemStore("/tmp", []byte("oursecret")),
	}
)

func FakeProfile() map[string]interface{} {
	profile := make(map[string]interface{})
	profile["picture"] = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"
	return profile
}

func RenderHTMLWithNoLayout(templateFileName string, data interface{}) []byte {
	return RenderHTML(templateFileName, data, render.HTMLOptions{Layout: ""})
}

func RenderHTML(templateFileName string, data interface{}, htmlOpt ...render.HTMLOptions) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0))
	err := FakeRenderer.HTML(buffer, http.StatusOK, templateFileName, data, htmlOpt...)
	if err != nil {
		panic("Failed to render template " + templateFileName + "; error: " + err.Error())
	}

	return buffer.Bytes()
}

type FakeTaxonomyRepo struct {
	list []*cohesioned.Taxonomy
	err  error
	key  *datastore.Key
}

func (r *FakeTaxonomyRepo) ListReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) ListChildrenReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) AddReturns(key string, err error) {
	r.key = &datastore.Key{Name: key}
	r.err = err
}

func (r *FakeTaxonomyRepo) List() ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) Add(t *cohesioned.Taxonomy) (*datastore.Key, error) {
	return r.key, r.err
}
