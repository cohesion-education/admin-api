package admin_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/admin"
	"github.com/cohesion-education/admin-api/pkg/common"
	"github.com/cohesion-education/admin-api/pkg/config"
)

func TestDashboardHandlerWhileLoggedInDirectsUserToDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	profile := fakes.FakeProfile()

	rr := httptest.NewRecorder()
	handler := admin.DashboardViewHandler(fakes.FakeHandlerConfig)
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)
	expectedBody := fakes.RenderHTML("admin/dashboard", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to dashboardHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}
