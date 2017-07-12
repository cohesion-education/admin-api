package video_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/gcp"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
)

func TestStreamHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/video/stream/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedVideo := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(expectedVideo, nil)

	gcpConfig, err := gcp.NewConfig("../../../testdata/test-gcp-keyfile.json")
	if err != nil {
		t.Fatalf("Failed to get gcp config %v", err)
	}

	renderer := fakes.FakeRenderer
	handler := video.StreamHandler(renderer, repo, gcpConfig)
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

func TestSaveHandler(t *testing.T) {
	v := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})

	params := map[string]string{
		"title":       v.Title,
		"taxonomy_id": fmt.Sprintf("%d", v.TaxonomyID),
	}

	req, err := fakes.NewfileUploadRequest("POST", "/api/video", params, "video_file", "../../../testdata/file-upload.txt")
	if err != nil {
		t.Fatal(err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.AddReturns(v, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.SaveHandler(renderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoAPIResponse{}
	fakeResp.RedirectURL = fmt.Sprintf("/admin/video/%d", v.ID())

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("JSON response not as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestUpdateHandler(t *testing.T) {
	v := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})

	params := map[string]string{
		"title":       v.Title,
		"id":          fmt.Sprintf("%d", v.ID()),
		"taxonomy_id": fmt.Sprintf("%d", v.TaxonomyID),
	}
	req, err := fakes.NewfileUploadRequest("PUT", "/api/video", params, "video_file", "../../../testdata/file-upload.txt")
	if err != nil {
		t.Fatal(err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.AddReturns(v, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.UpdateHandler(renderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoAPIResponse{}
	fakeResp.RedirectURL = fmt.Sprintf("/admin/video/%d", v.ID())

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("JSON response not as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}
