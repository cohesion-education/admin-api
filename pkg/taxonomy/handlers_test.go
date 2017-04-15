package taxonomy_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/common"
	"github.com/cohesion-education/admin-api/pkg/config"
	"github.com/cohesion-education/admin-api/pkg/taxonomy"
)

func TestListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/taxonomy", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//TODO - this won't work until we either inject the taxonomy repo, or mock the datastore client
	handler := taxonomy.ListHandler(fakes.FakeHandlerConfig)
	profile := fakes.FakeProfile()

	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)
	dashboard.Set("list", []taxonomy.Taxonomy{})

	expectedBody := fakes.RenderHTML("taxonomy/list", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected HTML was not generated in the call to dashboardHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestAddHandler(t *testing.T) {
	// req, err := http.NewRequest("GET", "/dashboard", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// profile := make(map[string]interface{})
	// profile["picture"] = "https://pbs.twimg.com/profile_images/2043299214/Adam_Avatar_Small_400x400.jpg"
	//
	// rr := httptest.NewRecorder()
	// handler := dashboardHandler(hc)
	// ctx := req.Context()
	// ctx = context.WithValue(ctx, currentUserKey, profile)
	// req = req.WithContext(ctx)
	//
	// handler.ServeHTTP(rr, req)
	//
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	// }
	//
	// dashboard := &dashboard{}
	// dashboard.Set("profile", profile)
	// expectedBody := renderHTML("admin/dashboard", dashboard)
	// if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
	// 	t.Errorf("The expected HTML was not generated in the call to dashboardHandler: Expected:\n\n%sActual:\n\n%s", string(expectedBody), rr.Body.String())
	// }
}
