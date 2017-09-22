package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/auth"
)

func TestIsAdminWhenAuthenticatedAsAdminProceedsToNext(t *testing.T) {
	req, err := http.NewRequest("GET", "/endpoint-that-requires-admin", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockNextHandlerCalled = true
	})

	adminUser := fakes.FakeAdmin()
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, adminUser)

	rr := httptest.NewRecorder()
	handler := auth.IsAdmin(fakes.FakeRenderer)
	handler.ServeHTTP(rr, req.WithContext(ctx), mockNextHandler)

	if !mockNextHandlerCalled {
		t.Errorf("'next' handler was not called")
	}

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

func TestIsAdminWhenAuthenticatedAsNonAdminReturnsUnauthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/endpoint-that-requires-admin", nil)
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
	handler := auth.IsAdmin(fakes.FakeRenderer)
	handler.ServeHTTP(rr, req.WithContext(ctx), mockNextHandler)

	if mockNextHandlerCalled {
		t.Errorf("'next' handler was called but it should not have been")
	}

	expectedStatus := http.StatusForbidden
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}

func TestIsAdminWhenUnauthenticatedReturnsNotAuthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/endpoint-that-requires-admin", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	mockNextHandlerCalled := false
	mockNextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockNextHandlerCalled = true
	})

	rr := httptest.NewRecorder()
	handler := auth.IsAdmin(fakes.FakeRenderer)
	handler.ServeHTTP(rr, req, mockNextHandler)

	if mockNextHandlerCalled {
		t.Errorf("'next' handler was called but it should not have been")
	}

	expectedStatus := http.StatusUnauthorized
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}
}
