package video_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/gorilla/mux"
)

func TestListHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	testVideo := fakes.FakeVideo()
	videos := []*cohesioned.Video{testVideo}

	fakeAdminService := new(fakes.FakeVideoAdminService)
	fakeAdminService.ListReturns(videos, nil)

	handler := video.ListHandler(fakes.FakeRenderer, fakeAdminService)
	rr := httptest.NewRecorder()

	req := fakes.NewRequestWithContext("GET", "/api/videos", nil, fakeUser)
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoResponse{
		List: videos,
	}

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}

func TestGetByIDHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	testVideo := fakes.FakeVideo()

	apiURL := fmt.Sprintf("/api/video/%d", testVideo.ID)

	fakeAdminService := new(fakes.FakeVideoAdminService)
	fakeAdminService.GetWithSignedURLReturns(testVideo, nil)

	handler := video.GetByIDHandler(fakes.FakeRenderer, fakeAdminService)
	router := mux.NewRouter()
	router.HandleFunc("/api/video/{id:[0-9]+}", handler)

	rr := httptest.NewRecorder()
	req := fakes.NewRequestWithContext("GET", apiURL, nil, fakeUser)

	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoResponse{
		Video: testVideo,
	}

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}

	expectedResp := &video.VideoResponse{}

	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&expectedResp); err != nil {
		t.Errorf("Failed to unmarshall response json to APIResponse: %v", err)
	}

	if len(expectedResp.Video.SignedURL) == 0 {
		t.Errorf("api response did not contain a signed url")
	}
}

func TestAddHandler(t *testing.T) {
	profile := fakes.FakeProfile()
	testVideo := fakes.FakeVideo()

	testJSON, err := json.Marshal(testVideo)
	if err != nil {
		t.Fatalf("Failed to marshall video json: %v", err)
	}

	fakeAdminService := new(fakes.FakeVideoAdminService)
	fakeAdminService.SaveReturns(nil)

	handler := video.AddHandler(fakes.FakeRenderer, fakeAdminService)
	rr := httptest.NewRecorder()

	req := fakes.NewRequestWithContext("POST", "/api/video", bytes.NewReader(testJSON), profile)
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoResponse{
		Video: testVideo,
	}

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}
}

func TestUploadHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	testVideo := fakes.FakeVideo()

	videoUploadURI := fmt.Sprintf("/api/video/upload/%d", testVideo.ID)

	fakeAdminService := new(fakes.FakeVideoAdminService)
	fakeAdminService.GetReturns(testVideo, nil)
	fakeAdminService.SetFileReturns(nil)

	handler := video.UploadHandler(fakes.FakeRenderer, fakeAdminService)

	router := mux.NewRouter()
	router.HandleFunc("/api/video/upload/{id:[0-9]+}", handler)
	rr := httptest.NewRecorder()

	req := fakes.NewFileUploadRequestWithContext("POST", videoUploadURI, "../../../testdata/file-upload.txt", fakeUser)
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoResponse{
		Video: testVideo,
	}

	expectedBody := fakes.RenderJSON(fakeResp)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected: %s\n\nActual: %s", string(expectedBody), rr.Body.String())
	}

	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&fakeResp); err != nil {
		t.Errorf("Failed to unmarshall response json to APIResponse: %v", err)
	}

	if fakeResp.Video.FileName != testVideo.FileName {
		t.Errorf("Video FileName was not set correctly; expected: %s - actual: %s", testVideo.FileName, fakeResp.Video.FileName)
	}
}

func TestUpdateHandler(t *testing.T) {
	fakeUser := fakes.FakeProfile()
	existingVideo := fakes.FakeVideo()
	testVideo := &cohesioned.Video{
		ID:          existingVideo.ID,
		Title:       "Updated Video Title",
		TaxonomyID:  existingVideo.TaxonomyID,
		FileName:    existingVideo.FileName,
		Created:     existingVideo.Created,
		CreatedByID: existingVideo.CreatedByID,
		Updated:     time.Now(),
		UpdatedByID: existingVideo.CreatedByID,
	}

	testJSON, err := json.Marshal(testVideo)
	if err != nil {
		t.Fatalf("Failed to marshall video json: %v", err)
	}

	apiURL := fmt.Sprintf("/api/video/%d", testVideo.ID)

	fakeAdminService := new(fakes.FakeVideoAdminService)
	fakeAdminService.GetReturns(existingVideo, nil)
	fakeAdminService.UpdateReturns(nil)

	handler := video.UpdateHandler(fakes.FakeRenderer, fakeAdminService)
	router := mux.NewRouter()
	router.HandleFunc("/api/video/{id:[0-9]+}", handler)
	rr := httptest.NewRecorder()

	req := fakes.NewRequestWithContext("PUT", apiURL, bytes.NewReader(testJSON), fakeUser)
	router.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		fmt.Printf("response %s\n", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	fakeResp := &video.VideoResponse{}

	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&fakeResp); err != nil {
		t.Errorf("Failed to unmarshall response json to APIResponse: %v", err)
	}

	if fakeResp.Video.UpdatedByID != fakeUser.ID {
		t.Errorf("Video UpdatedBy was not set correctly; expected: %d - actual: %d", fakeUser.ID, fakeResp.Video.UpdatedByID)
	}
}
