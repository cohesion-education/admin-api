package video_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/gcp"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
)

func TestListHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	testVideo := cohesioned.NewVideo("Test Video", "test.mp4", 1234, fakeUser)
	videos := []*cohesioned.Video{testVideo}

	req, err := http.NewRequest("GET", "/api/videos", nil)
	if err != nil {
		t.Fatalf("Failed to initialize get videos request %v", err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.ListReturns(videos, nil)

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, fakeUser)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.ListHandler(renderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoAPIResponse{
		List: videos,
	}

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}

func TestGetByIDHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	testVideo := cohesioned.NewVideo("Test Before", "test.mp4", 1234, fakeUser)

	apiURL := fmt.Sprintf("/api/video/%d", testVideo.ID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		t.Fatalf("Failed to initialize create video request %v", err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(testVideo, nil)

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, fakeUser)
	req = req.WithContext(ctx)

	gcpConfig, err := gcp.NewConfig("../../../testdata/test-gcp-keyfile.json")
	if err != nil {
		t.Fatalf("Failed to get gcp config %v", err)
	}

	renderer := fakes.FakeRenderer
	handler := video.GetByIDHandler(renderer, repo, gcpConfig)

	router := mux.NewRouter()
	router.HandleFunc("/api/video/{id:[0-9]+}", handler)

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

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}

	expectedResp := &video.VideoAPIResponse{}

	decoder := json.NewDecoder(rr.Body)
	if err = decoder.Decode(&expectedResp); err != nil {
		t.Errorf("Failed to unmarshall response json to VideoAPIResponse: %v", err)
	}

	if len(expectedResp.Video.SignedURL) == 0 {
		t.Errorf("api response did not contain a signed url")
	}

	expectedHostName := "storage.googleapis.com"
	signedURL, err := url.Parse(expectedResp.Video.SignedURL)
	if err != nil {
		t.Fatalf("%s does not appear to be a valid URL: %v", expectedResp.Video.SignedURL, err)
	}

	if signedURL.Host != expectedHostName {
		t.Errorf("expected signed url host to be %s but was %s", expectedHostName, signedURL.Host)
	}

	q := signedURL.Query()

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

// func TestStreamHandler(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/api/video/stream/1234", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	expectedVideo := cohesioned.NewVideo("Test Video", "test.mp4", "test-bucket", "1234-test.mp4", 1234, 1, &cohesioned.Profile{FullName: "Test User"})
//
// 	repo := new(fakes.FakeVideoRepo)
// 	repo.GetReturns(expectedVideo, nil)
//
// 	gcpConfig, err := gcp.NewConfig("../../../testdata/test-gcp-keyfile.json")
// 	if err != nil {
// 		t.Fatalf("Failed to get gcp config %v", err)
// 	}
//
// 	renderer := fakes.FakeRenderer
// 	handler := video.StreamHandler(renderer, repo, gcpConfig)
// 	rr := httptest.NewRecorder()
//
// 	router := mux.NewRouter()
// 	router.HandleFunc("/api/video/stream/{id:[0-9]+}", handler).Methods("GET")
// 	router.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusSeeOther {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
// 	}
//
// 	expectedHostName := "storage.googleapis.com"
// 	location, err := rr.Result().Location()
// 	if err != nil {
// 		t.Fatalf("Failed to read request recorder result's location %v", err)
// 	}
//
// 	if location.Host != expectedHostName {
// 		t.Errorf("expected signed url host to be %s but was %s", expectedHostName, location.Host)
// 	}
//
// 	q := location.Query()
//
// 	if len(q.Get("GoogleAccessId")) == 0 {
// 		t.Errorf("Signed url did not have GoogleAccessId param")
// 	}
//
// 	if len(q.Get("Expires")) == 0 {
// 		t.Errorf("Signed url did not have Expires param")
// 	}
//
// 	if len(q.Get("Signature")) == 0 {
// 		t.Errorf("Signed url did not have Signature param")
// 	}
// }

func TestAddHandler(t *testing.T) {
	testVideo := cohesioned.NewVideo("Test After", "test.mp4", 1234, &cohesioned.Profile{FullName: "Test User"})
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
	fakeResp.ID = testVideo.ID

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}

func TestUploadHandler(t *testing.T) {
	testVideo := cohesioned.NewVideo("Test After", "test.mp4", 1234, &cohesioned.Profile{FullName: "Test User"})

	videoUploadURI := fmt.Sprintf("/api/video/upload/%d", testVideo.ID)
	req, err := fakes.NewFileUploadRequest("POST", videoUploadURI, "../../../testdata/file-upload.txt")
	if err != nil {
		t.Fatal(err)
	}

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
	fakeResp.ID = testVideo.ID

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}

	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&fakeResp); err != nil {
		t.Errorf("Failed to unmarshall response json to VideoAPIResponse: %v", err)
	}

	if fakeResp.Video.FileName != testVideo.FileName {
		t.Errorf("Video FileName was not set correctly; expected: %s - actual: %s", testVideo.FileName, fakeResp.Video.FileName)
	}
}

func TestUpdateHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	existingVideo := cohesioned.NewVideo("Test Before", "test.mp4", 1234, fakeUser)
	testVideo := cohesioned.NewVideo("Test After", "test.mp4", 1234, fakeUser)

	testJSON, err := testVideo.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshall video json: %v", err)
	}

	apiURL := fmt.Sprintf("/api/video/%d", testVideo.ID)
	req, err := http.NewRequest("PUT", apiURL, bytes.NewReader(testJSON))
	if err != nil {
		t.Fatalf("Failed to initialize create video request %v", err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.GetReturns(existingVideo, nil)
	repo.UpdateReturns(existingVideo, nil)

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, fakeUser)
	req = req.WithContext(ctx)

	renderer := fakes.FakeRenderer
	handler := video.UpdateHandler(renderer, repo)

	router := mux.NewRouter()
	router.HandleFunc("/api/video/{id:[0-9]+}", handler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := video.NewAPIResponse(existingVideo)

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}
