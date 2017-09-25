package fakes

import (
	"bytes"
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

func NewFileUploadRequest(method string, uri string, localFilePath string) (*http.Request, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s does not seem to be a valid file path: %v", localFilePath, err)
	}

	req, err := http.NewRequest(method, uri, file)
	return req, err
}
