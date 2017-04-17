package auth_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
)

func TestLoginViewHandlerWhileNotLoggedInDirectsUserToLoginPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := auth.LoginViewHandler(fakes.FakeRenderer)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := fakes.RenderHTMLWithNoLayout("login/index", nil)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to loginViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestLoginViewHandlerWhileLoggedInDirectsUserToDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := &cohesioned.Profile{}
	profile.PictureURL = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"

	rr := httptest.NewRecorder()
	handler := auth.LoginViewHandler(fakes.FakeRenderer)
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatal(err)
	}

	if location.Path != "/admin/dashboard" {
		t.Errorf("Expected to get redirected to /dashboard but was redirected to %s", location.Path)
	}
}
