package fakes

import (
	"bytes"
	"encoding/gob"
	"net/http"

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

func init() {
	gob.Register(&cohesioned.Profile{})
}

func FakeProfile() *cohesioned.Profile {
	return &cohesioned.Profile{
		FullName:   "Test User",
		Email:      "hello@domain.com",
		UserID:     "abc|123",
		PictureURL: "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg",
	}
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

func RenderJSON(data interface{}) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0))
	err := FakeRenderer.JSON(buffer, http.StatusOK, data)
	if err != nil {
		panic("Failed to render JSON: " + err.Error())
	}

	return buffer.Bytes()
}

type FakeTaxonomyRepo struct {
	t    *cohesioned.Taxonomy
	list []*cohesioned.Taxonomy
	err  error
}

func (r *FakeTaxonomyRepo) ListReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) ListChildrenReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) AddReturns(t *cohesioned.Taxonomy, err error) {
	r.t = t
	r.err = err
}

func (r *FakeTaxonomyRepo) List() ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) Add(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	return r.t, r.err
}
