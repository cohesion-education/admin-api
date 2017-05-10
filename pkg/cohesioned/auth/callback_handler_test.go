package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/auth"
)

func TestAuthenticatedAsNonAdminUserRedirectsToUserDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/callback?code=abc123", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	fakeAuth0Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/oauth/token":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "access_token=A9CvPwFojaBI&token_type=bearer&id_token=eyJ0eXAiOiJKV1Qi")
		case "/userinfo":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, userInfoPayload)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "No mapping in fake auth0 server for "+req.RequestURI)
		}
	}))
	defer fakeAuth0Server.Close()

	fakes.FakeAuthConfig.Domain = fakeAuth0Server.URL
	fakes.FakeAuthConfig.CallbackURL = req.URL.String()
	fakes.FakeAuthConfig.ClientID = "test-client-id"
	fakes.FakeAuthConfig.ClientSecret = "test-client-secret"
	handler := auth.CallbackHandler(fakes.FakeAuthConfig)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("unable to retreive rr.Result().Location() %v", err)
	}

	expectedRedirectTo := "/dashboard"
	if expectedRedirectTo != location.Path {
		t.Errorf("Expected request redirect url to be %s but was %s", expectedRedirectTo, location.Path)
	}

	session, err := fakes.FakeAuthConfig.GetCurrentSession(req)
	if err != nil {
		t.Fatalf("Failed to get current session %v", err)
	}

	profile, ok := session.Values[cohesioned.CurrentUserSessionKey]
	if !ok {
		t.Errorf("Session did not contain current user")
	}

	if profile == nil {
		t.Errorf("Session contained nil profile")
	}
}

func TestAuthenticatedAsAdminUserRedirectsToAdminDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/callback?code=abc123", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	fakeAuth0Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/oauth/token":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "access_token=A9CvPwFojaBI&token_type=bearer&id_token=eyJ0eXAiOiJKV1Qi")
		case "/userinfo":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, adminUserInfoPayload)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "No mapping in fake auth0 server for "+req.RequestURI)
		}
	}))
	defer fakeAuth0Server.Close()

	fakes.FakeAuthConfig.Domain = fakeAuth0Server.URL
	fakes.FakeAuthConfig.CallbackURL = req.URL.String()
	fakes.FakeAuthConfig.ClientID = "test-client-id"
	fakes.FakeAuthConfig.ClientSecret = "test-client-secret"
	handler := auth.CallbackHandler(fakes.FakeAuthConfig)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("unable to retreive rr.Result().Location() %v", err)
	}

	expectedRedirectTo := "/admin/dashboard"
	if expectedRedirectTo != location.Path {
		t.Errorf("Expected request redirect url to be %s but was %s", expectedRedirectTo, location.Path)
	}

	session, err := fakes.FakeAuthConfig.GetCurrentSession(req)
	if err != nil {
		t.Fatalf("Failed to get current session %v", err)
	}

	profile, ok := session.Values[cohesioned.CurrentUserSessionKey]
	if !ok {
		t.Errorf("Session did not contain current user")
	}

	if profile == nil {
		t.Errorf("Session contained nil profile")
	}
}

func TestExchangeTokenErrorRedirectsTo500(t *testing.T) {
	req, err := http.NewRequest("GET", "/callback?code=abc123", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	fakeAuth0Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/oauth/token":
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "mock failure")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "No mapping in fake auth0 server for "+req.RequestURI)
		}
	}))
	defer fakeAuth0Server.Close()

	fakes.FakeAuthConfig.Domain = fakeAuth0Server.URL
	fakes.FakeAuthConfig.CallbackURL = req.URL.String()
	fakes.FakeAuthConfig.ClientID = "test-client-id"
	fakes.FakeAuthConfig.ClientSecret = "test-client-secret"
	handler := auth.CallbackHandler(fakes.FakeAuthConfig)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusSeeOther
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("unable to retreive rr.Result().Location() %v", err)
	}

	expectedRedirectTo := "/500"
	if expectedRedirectTo != location.Path {
		t.Errorf("Expected request redirect url to be %s but was %s", expectedRedirectTo, location.Path)
	}
}

func TestGetUserInfoErrorRedirectsTo500(t *testing.T) {
	req, err := http.NewRequest("GET", "/callback?code=abc123", nil)
	if err != nil {
		t.Fatalf("Failed to initialize new request %v", err)
	}

	fakeAuth0Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/oauth/token":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "access_token=A9CvPwFojaBI&token_type=bearer&id_token=eyJ0eXAiOiJKV1Qi")
		case "/userinfo":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "mock failure")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "No mapping in fake auth0 server for "+req.RequestURI)
		}
	}))
	defer fakeAuth0Server.Close()

	fakes.FakeAuthConfig.Domain = fakeAuth0Server.URL
	fakes.FakeAuthConfig.CallbackURL = req.URL.String()
	fakes.FakeAuthConfig.ClientID = "test-client-id"
	fakes.FakeAuthConfig.ClientSecret = "test-client-secret"
	handler := auth.CallbackHandler(fakes.FakeAuthConfig)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusSeeOther
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
	}

	location, err := rr.Result().Location()
	if err != nil {
		t.Fatalf("unable to retreive rr.Result().Location() %v", err)
	}

	expectedRedirectTo := "/500"
	if expectedRedirectTo != location.Path {
		t.Errorf("Expected request redirect url to be %s but was %s", expectedRedirectTo, location.Path)
	}
}
