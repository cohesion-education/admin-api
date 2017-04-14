package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboardHandlerWhileLoggedInDirectsUserToDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := make(map[string]interface{})
	profile["picture"] = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"

	rr := httptest.NewRecorder()
	handler := dashboardHandler(hc)
	ctx := req.Context()
	ctx = context.WithValue(ctx, currentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := renderHTML("admin/dashboard", profile)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to dashboardHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}
