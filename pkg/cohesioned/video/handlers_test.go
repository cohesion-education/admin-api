package video_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/gcp"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
)

func TestListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/video", nil)
	if err != nil {
		t.Fatal(err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.ListReturns([]*cohesioned.Video{}, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler := video.ListHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)
	dashboard.Set("list", []*cohesioned.Video{})

	expectedBody := fakes.RenderHTML("video/list", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestFormHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/video/add", nil)
	if err != nil {
		t.Fatal(err)
	}

	repo := new(fakes.FakeVideoRepo)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler := video.FormHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)

	expectedBody := fakes.RenderHTML("video/form", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestShowHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/video/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedVideo := &cohesioned.Video{
		Key:     datastore.IDKey("Video", 1234, nil),
		Created: time.Now(),
		CreatedBy: &cohesioned.Profile{
			FullName: "Test User",
		},
		Title:             "Test Video",
		FileName:          "test.mp4",
		StorageBucket:     "test-bucket",
		StorageObjectName: "1234-test.mp4",
		TaxonomyID:        1,
	}

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(expectedVideo, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler := video.ShowHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/admin/video/{id:[0-9]+}", handler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)
	dashboard.Set("video", expectedVideo)

	expectedBody := fakes.RenderHTML("video/show", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestStreamHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/video/stream/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedVideo := &cohesioned.Video{
		Key:     datastore.IDKey("Video", 1234, nil),
		Created: time.Now(),
		CreatedBy: &cohesioned.Profile{
			FullName: "Test User",
		},
		Title:             "Test Video",
		FileName:          "test.mp4",
		StorageBucket:     "test-bucket",
		StorageObjectName: "1234-test.mp4",
		TaxonomyID:        1,
	}

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(expectedVideo, nil)

	gcpConfig, err := gcp.NewConfig("../../../gcp-keyfile.json")
	if err != nil {
		t.Fatalf("Failed to get gcp config %v", err)
	}

	handler := video.StreamHandler(fakes.FakeRenderer, repo, gcpConfig)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/video/stream/{id:[0-9]+}", handler).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	expectedHostName := "storage.googleapis.com"
	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("Failed to read request recorder result's location %v", err)
	}

	if location.Host != expectedHostName {
		t.Errorf("expected signed url host to be %s but was %s", expectedHostName, location.Host)
	}

	q := location.Query()

	if len(q.Get("GoogleAccessId")) == 0 {
		t.Errorf("Signed url did not have GoogleAccessId param")
	}

	if len(q.Get("Expires")) == 0 {
		t.Errorf("Signed url did not have Expires param")
	}

	if len(q.Get("Signature")) == 0 {
		t.Errorf("Signed url did not have Signature param")
	}
}
