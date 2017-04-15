package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/auth"
	"github.com/cohesion-education/admin-api/pkg/config"
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

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("unable to retreive rr.Result().Location() %v", err)
	}

	if "/" != location.Path {
		t.Errorf("Expected request redirect url to be %s but was %s", "/", req.URL.Path)
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

	profile := make(map[string]interface{})
	profile["picture"] = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"
	session.Values[config.CurrentUserSessionKey] = profile

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true
		ctx := req.Context()
		if ctx.Value(config.CurrentUserKey) == nil {
			t.Errorf("middleware did not set current user in the context as expected")
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
