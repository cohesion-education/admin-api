package taxonomy_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/taxonomy"
	"github.com/gorilla/mux"
)

func TestListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/taxonomy", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	repo := new(fakes.FakeTaxonomyRepo)
	repo.ListReturns([]*cohesioned.Taxonomy{}, nil)

	handler := taxonomy.ListHandler(fakes.FakeRenderer, repo)
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
	dashboard.Set("list", []*cohesioned.Taxonomy{})

	expectedBody := fakes.RenderHTML("taxonomy/list", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
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

func TestListChildrenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/taxonomy/1234/children", nil)
	if err != nil {
		t.Fatal(err)
	}

	fakeList := []*cohesioned.Taxonomy{
		&cohesioned.Taxonomy{
			Name:     "test-child-1",
			Created:  time.Now(),
			ParentID: 1234,
		},
		&cohesioned.Taxonomy{
			Name:     "test-child-2",
			Created:  time.Now(),
			ParentID: 1234,
		},
	}

	repo := new(fakes.FakeTaxonomyRepo)
	repo.ListChildrenReturns(fakeList, nil)

	rr := httptest.NewRecorder()
	handler := taxonomy.ListChildrenHandler(fakes.FakeRenderer, repo)

	router := mux.NewRouter()
	router.HandleFunc("/api/taxonomy/{id:[0-9]+}/children", handler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data := struct {
		Children []*cohesioned.Taxonomy `json:"children"`
		ParentID int64                  `json:"parent_id"`
	}{
		fakeList,
		1234,
	}

	expectedBody := fakes.RenderJSON(data)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected:\n%s\n\nActual:\n%s", string(expectedBody), rr.Body.String())
	}
}
