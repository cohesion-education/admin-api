package profile_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
)

func TestSave(t *testing.T) {
	profilePayload := `{
		"email":"john@doe.com",
		"name":"John Doe",
		"state":"FL",
		"county":"Monroe County",
		"students":[
			{
				"name":"Billy",
				"grade":"1",
				"school":"Key West Elementary"
			}
		]
	}`

	req, err := http.NewRequest("POST", "/api/profile", strings.NewReader(profilePayload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	renderer := fakes.FakeRenderer
	repo := new(fakes.FakeProfileRepo)
	repo.SaveReturns(5, nil)

	handler := profile.SaveHandler(renderer, repo)
	rr := httptest.NewRecorder()

	p := fakes.FakeProfile()
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, p)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	result := &cohesioned.Profile{}
	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("failed to unmarshall json response %v", err)
	}

	if result.Email != "john@doe.com" {
		t.Error("Email did not update")
	}

	if result.FullName != "John Doe" {
		t.Error("Name did not update")
	}

	if result.State != "FL" {
		t.Error("State did not update")
	}

	if result.County != "Monroe County" {
		t.Error("County did not update")
	}

	if len(result.Students) != 1 {
		t.Error("There should be exactly 1 Student in Students")
	}

	if result.Students[0].Name != "Billy" {
		t.Error("Students[0].Name did not update")
	}

	if result.Students[0].Grade != "1" {
		t.Error("Students[0].Grade did not update")
	}

	if result.Students[0].School != "Key West Elementary" {
		t.Error("Students[0].School did not update")
	}
}

func TestSavePreferences(t *testing.T) {
	prefs := `{"newsletter":true, "beta_program":true}`

	req, err := http.NewRequest("POST", "/api/profile/preferences", strings.NewReader(prefs))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	renderer := fakes.FakeRenderer
	repo := new(fakes.FakeProfileRepo)
	repo.SaveReturns(5, nil)

	handler := profile.SavePreferencesHandler(renderer, repo)
	rr := httptest.NewRecorder()

	p := fakes.FakeProfile()
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, p)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	result := &cohesioned.Profile{}
	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("failed to unmarshall json response %v", err)
	}

	if result.Preferences.Newsletter != true {
		t.Error("Newletter preferences did not save")
	}

	if result.Preferences.BetaProgram != true {
		t.Error("Beta Program preferences did not save")
	}
}

func TestGetCurrentUserHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	fakeProfile := fakes.FakeProfile()

	renderer := fakes.FakeRenderer
	repo := new(fakes.FakeProfileRepo)
	repo.FindByEmailReturns(fakeProfile, nil)

	handler := profile.GetCurrentUserHandler(renderer, repo)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, fakeProfile)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	result := &cohesioned.Profile{}
	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("failed to unmarshall json response %v", err)
	}

	if result.Email != fakeProfile.Email {
		t.Error("Email was not set correctly on the result")
	}

	if result.FirstName != fakeProfile.FirstName {
		t.Error("FirstName was not set correctly on the result")
	}
}
