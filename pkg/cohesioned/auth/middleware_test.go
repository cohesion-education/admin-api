package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
)

func TestIsAuthenticatedWhenNotAuthenticatedRedirectsTo401(t *testing.T) {
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

	expectedStatus := http.StatusSeeOther
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Errorf("Failed to get result location from recorder %v", err)
	}

	expectedLocation := "/401"
	if location.String() != expectedLocation {
		t.Errorf("handler returned wrong redirect url: got %s want %s", location.String(), expectedLocation)
	}
}

func TestIsAuthenticatedWhenAuthenticatedProceedsToNext(t *testing.T) {
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
	session.Values[cohesioned.CurrentUserSessionKey] = profile

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true

		profile := req.Context().Value(cohesioned.CurrentUserKey)
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

func TestIsAdminWhenAuthenticatedAsAdminProceedsToNext(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true
	})

	profile := fakes.FakeAdmin()
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)

	rr := httptest.NewRecorder()
	auth.IsAdmin(rr, req.WithContext(ctx), mockNextHandler)

	if !mockNextHandlerCalled {
		t.Errorf("'next' handler was not called")
	}

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

func TestIsAdminWhenAuthenticatedAsNonAdminRedirectsTo403(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true
	})

	profile := fakes.FakeProfile()
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)

	rr := httptest.NewRecorder()
	auth.IsAdmin(rr, req.WithContext(ctx), mockNextHandler)

	if mockNextHandlerCalled {
		t.Errorf("'next' handler was called but it should not have been")
	}

	expectedStatus := http.StatusSeeOther
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Errorf("unable to get result location from request recorder %v", err)
	}

	expectedRedirectTo := "/403"
	if location.String() != expectedRedirectTo {
		t.Errorf("handler did not redirect to the correct page - expected %s got %s", expectedRedirectTo, location.String())
	}
}

func TestIsAdminWhenUnauthenticatedRedirectsTo401(t *testing.T) {
	req, err := http.NewRequest("GET", "/admin/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mockNextHandlerCalled = true
	})

	rr := httptest.NewRecorder()
	auth.IsAdmin(rr, req, mockNextHandler)

	if mockNextHandlerCalled {
		t.Errorf("'next' handler was called but it should not have been")
	}

	expectedStatus := http.StatusSeeOther
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Errorf("unable to get result location from request recorder %v", err)
	}

	expectedRedirectTo := "/401"
	if location.String() != expectedRedirectTo {
		t.Errorf("handler did not redirect to the correct page - expected %s got %s", expectedRedirectTo, location.String())
	}
}
