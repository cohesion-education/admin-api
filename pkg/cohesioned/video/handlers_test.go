package video_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAddHandler(t *testing.T) {
	testVideo := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})

	testJSON, err := testVideo.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshall video json: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/video", bytes.NewReader(testJSON))
	if err != nil {
		t.Fatalf("Failed to initialize create video request %v", err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.AddReturns(testVideo, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.AddHandler(renderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoAPIResponse{
		Video: testVideo,
	}
	fakeResp.ID = testVideo.ID()

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}

func TestUploadHandler(t *testing.T) {
	testVideo := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})

	videoUploadURI := fmt.Sprintf("/api/video/upload/%d", testVideo.ID())
	req, err := fakes.NewFileUploadRequest("POST", videoUploadURI, "../../../testdata/file-upload.txt")
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("file-name", "file-upload.txt")

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(testVideo, nil)
	repo.SetFileReturns(testVideo, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.UploadHandler(renderer, repo)

	router := mux.NewRouter()
	router.HandleFunc("/api/video/upload/{id:[0-9]+}", handler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoAPIResponse{
		Video: testVideo,
	}
	fakeResp.ID = testVideo.ID()

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}

	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&fakeResp); err != nil {
		t.Errorf("Failed to unmarshall response json to VideoAPIResponse: %v", err)
	}

	expectedVideoFileName := "file-upload.txt"
	if fakeResp.Video.FileName != expectedVideoFileName {
		t.Errorf("Video FileName was not set correctly; expected: %s - actual: %s", expectedVideoFileName, fakeResp.Video.FileName)
	}
}

func TestUpdateHandler(t *testing.T) {
	//TODO - this entire function needs to be refactored
	// v := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})
	//
	// params := map[string]string{
	// 	"title":       v.Title,
	// 	"id":          fmt.Sprintf("%d", v.ID()),
	// 	"taxonomy_id": fmt.Sprintf("%d", v.TaxonomyID),
	// }
	// req, err := fakes.NewfileUploadRequest("PUT", "/api/video", params, "video_file", "../../../testdata/file-upload.txt")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// repo := new(fakes.FakeVideoRepo)
	// repo.AddReturns(v, nil)
	//
	// profile := fakes.FakeProfile()
	// ctx := req.Context()
	// ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	// req = req.WithContext(ctx)
	//
	// renderer := fakes.FakeRenderer
	// handler := video.UpdateHandler(renderer, repo)
	// rr := httptest.NewRecorder()
	// handler.ServeHTTP(rr, req)
	//
	// expectedStatus := http.StatusOK
	// if status := rr.Code; status != expectedStatus {
	// 	fmt.Printf("response %s\n", rr.Body.String())
	// 	t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	// }
	//
	// fakeResp := &video.VideoAPIResponse{}
	// fakeResp.RedirectURL = fmt.Sprintf("/admin/video/%d", v.ID())
	//
	// expectedBody := fakes.RenderJSON(fakeResp)
	// if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
	// 	t.Errorf("JSON response not as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	// }
}
