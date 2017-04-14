package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginViewHandlerWhileNotLoggedInDirectsUserToLoginPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := loginViewHandler(hc)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := renderHTMLWithNoLayout("login/index", nil)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to loginViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestLoginViewHandlerWhileLoggedInDirectsUserToDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := make(map[string]interface{})
	profile["picture"] = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"

	rr := httptest.NewRecorder()
	handler := loginViewHandler(hc)
	ctx := req.Context()
	ctx = context.WithValue(ctx, currentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatal(err)
	}

	if location.Path != "/dashboard" {
		t.Errorf("Expected to get redirected to /dashboard but was redirected to %s", location.Path)
	}
}
