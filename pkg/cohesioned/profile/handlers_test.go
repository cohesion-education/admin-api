package profile_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
)

func TestSavePreferences(t *testing.T) {
	form := url.Values{}
	form.Add("preferences.newsletter", "on")
	form.Add("preferences.betaprogram", "on")

	req, err := http.NewRequest("POST", "/api/profile/preferences", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	renderer := fakes.FakeRenderer
	repo := new(fakes.FakeProfileRepo)
	repo.SaveReturns(nil)

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
		t.Error("Newletter preferences did not update")
	}

	if result.Preferences.BetaProgram != true {
		t.Error("Beta Program preferences did not update")
	}
}
