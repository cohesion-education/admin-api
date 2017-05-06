package fakes

import (
	"bytes"
	"encoding/gob"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var (
	FakeRenderer = render.New()

	FakeAdminDashboardRenderer = render.New(render.Options{
		Layout: "dashboard/admin-layout",
		RenderPartialsWithoutPrefix: true,
		Directory:                   "../../../templates",
	})

	FakeUserDashboardRenderer = render.New(render.Options{
		Layout: "dashboard/user-layout",
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
	return RenderHTML(FakeRenderer, templateFileName, data, render.HTMLOptions{Layout: ""})
}

func RenderHTML(renderer *render.Render, templateFileName string, data interface{}, htmlOpt ...render.HTMLOptions) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0))
	err := renderer.HTML(buffer, http.StatusOK, templateFileName, data, htmlOpt...)
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

//NewfileUploadRequest Creates a new file upload http request with optional extra params
func NewfileUploadRequest(method string, uri string, params map[string]string, paramName, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
