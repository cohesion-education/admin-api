package fakes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/unrolled/render"
)

var (
	FakeRenderer   = render.New()
	FakeAuthConfig = &config.AuthConfig{}
)

func FakeProfile() *cohesioned.Profile {
	return &cohesioned.Profile{
		ID:         1,
		Created:    time.Now(),
		FullName:   "Test User",
		FirstName:  "Test",
		LastName:   "User",
		Email:      "hello@domain.com",
		Sub:        "abc|123",
		PictureURL: "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg",
	}
}

func FakeAdmin() *cohesioned.Profile {
	return &cohesioned.Profile{
		FullName:      "Test User",
		Email:         "admin@cohesioned.io",
		EmailVerified: true,
		Sub:           "abc|123",
		PictureURL:    "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg",
	}
}

func FakeVideo() *cohesioned.Video {
	return &cohesioned.Video{
		ID:         1,
		Title:      "Test Video",
		FileName:   "test.mp4",
		TaxonomyID: FakeTaxonomy().ID,
		Created:    time.Now(),
		CreatedBy:  FakeProfile().ID,
		SignedURL:  "http://fake-signed-url",
	}
}

func FakeTaxonomy() *cohesioned.Taxonomy {
	return &cohesioned.Taxonomy{
		ID:        1,
		Name:      "Test Taxonomy",
		Created:   time.Now(),
		CreatedBy: FakeProfile().ID,
	}
}

func FakeStudent() *cohesioned.Student {
	return &cohesioned.Student{
		ID:        2,
		Name:      "Little Bobby Tables",
		Grade:     "3rd",
		School:    "Cohesion Ed Elementary",
		ParentID:  FakeProfile().ID,
		Created:   time.Now(),
		CreatedBy: FakeProfile().ID,
	}
}

func NewRequestWithContext(method, urlStr string, body io.Reader, user *cohesioned.Profile) *http.Request {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(fmt.Sprintf("http.NewRequest failed. Method: %s URL: %s", method, urlStr))
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, user)
	req = req.WithContext(ctx)
	return req
}

func WithContext(req *http.Request, user *cohesioned.Profile) {

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
func NewMultipartFileUploadRequest(method string, uri string, params map[string]string, paramName, filePath string) (*http.Request, error) {
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

func NewFileUploadRequestWithContext(method string, uri string, localFilePath string, user *cohesioned.Profile) *http.Request {
	file, err := os.Open(localFilePath)
	if err != nil {
		panic(fmt.Sprintf("%s does not seem to be a valid file path: %v", localFilePath, err))
	}

	return NewRequestWithContext(method, uri, file, user)
}
