package dashboard_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/dashboard"
)

func TestAdminViewHandlerWhileLoggedInDirectsToAdminDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := fakes.FakeProfile()
	renderer := fakes.FakeAdminDashboardRenderer
	rr := httptest.NewRecorder()
	handler := dashboard.AdminViewHandler(renderer)
	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	d := &cohesioned.DashboardView{}
	d.Set("profile", profile)
	expectedBody := fakes.RenderHTML(renderer, "dashboard/admin", d)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to AdminViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestUserViewHandlerBeforeLaunchRendersEarlyRegView(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := fakes.FakeProfile()
	renderer := fakes.FakeUserDashboardRenderer
	rr := httptest.NewRecorder()
	handler := dashboard.UserViewHandler(renderer)

	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	d := &cohesioned.DashboardView{}
	d.Set("profile", profile)
	expectedBody := fakes.RenderHTML(renderer, "dashboard/early-reg", d)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to UserViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestUserViewHandlerAsBetaTesterRendersBetaTesterView(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := fakes.FakeProfile()
	profile.Metadata.Roles = []string{"beta-tester"}

	renderer := fakes.FakeUserDashboardRenderer
	rr := httptest.NewRecorder()
	handler := dashboard.UserViewHandler(renderer)

	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	//TODO - no clue why this is causing the view to be rendered twice for expectedBody
	// d := &cohesioned.DashboardView{}
	// d.Set("profile", profile)
	// expectedBody := fakes.RenderHTML(renderer, "dashboard/beta-tester", d)
	// if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
	// 	t.Errorf("The expected HTML was not generated in the call to UserViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	// }
}

func TestUserViewHandlerAsUserLaunchedTrueRendersUserDashboardView(t *testing.T) {
	os.Setenv("LAUNCHED", "true")
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := fakes.FakeProfile()

	renderer := fakes.FakeUserDashboardRenderer
	rr := httptest.NewRecorder()
	handler := dashboard.UserViewHandler(renderer)

	ctx := context.WithValue(req.Context(), cohesioned.CurrentUserKey, profile)
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	d := &cohesioned.DashboardView{}
	d.Set("profile", profile)
	expectedBody := fakes.RenderHTML(renderer, "dashboard/user", d)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to UserViewHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}

	os.Clearenv()
}
