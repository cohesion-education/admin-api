package dashboard_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
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
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	d := &cohesioned.DashboardView{}
	d.Set("profile", profile)
	expectedBody := fakes.RenderHTML(renderer, "dashboard/admin", d)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to dashboardHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}
