package auth_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
)

func TestIsAuthenticatedHandlerWhenNotAuthenticatedRedirectsToRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	rr := httptest.NewRecorder()
	handler := auth.IsAuthenticatedHandler(fakes.FakeAuthConfig)

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true
	})

	handler.ServeHTTP(rr, req, mockNextHandler)

	if mockNextHandlerCalled {
		t.Errorf("next handler was called but it should not have been")
	}

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expectedBody := "Failed to get current user from session"
	if strings.Compare(strings.Trim(rr.Body.String(), "\n"), expectedBody) != 0 {
		t.Errorf("handler returned wrong response body: got %s want %s", rr.Body.String(), expectedBody)
	}
}

func TestIsAuthenticatedHandlerWhenAuthenticatedProceedsToNext(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	rr := httptest.NewRecorder()
	handler := auth.IsAuthenticatedHandler(fakes.FakeAuthConfig)
	session, err := fakes.FakeAuthConfig.GetCurrentSession(req)
	if err != nil {
		t.Fatalf("Failed to get current session %v", err)
	}

	profile := fakes.FakeProfile()
	session.Values[config.CurrentUserSessionKey] = profile

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true

		profile := req.Context().Value(config.CurrentUserKey)
		if profile == nil {
			t.Errorf("middleware did not set current user in the context as expected")
		}

		_, ok := profile.(*cohesioned.Profile)
		if !ok {
			t.Errorf("current user was not of type cohesioned.Profile; type: " + reflect.TypeOf(profile).Name())
		}
	})

	handler.ServeHTTP(rr, req, mockNextHandler)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !mockNextHandlerCalled {
		t.Errorf("next handler was not called as expected")
	}
}
