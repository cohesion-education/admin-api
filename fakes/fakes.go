package fakes

import (
	"bytes"
	"net/http"

	"github.com/cohesion-education/admin-api/pkg/config"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var (
	FakeRenderer = render.New(render.Options{
		Layout: "layout",
		RenderPartialsWithoutPrefix: true,
		Directory:                   "../../templates",
	})

	FakeAuthConfig = &config.AuthConfig{
		SessionStore: sessions.NewFilesystemStore("/tmp", []byte("oursecret")),
	}

	FakeHandlerConfig = &config.HandlerConfig{
		Renderer:        FakeRenderer,
		DatastoreClient: nil,
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
