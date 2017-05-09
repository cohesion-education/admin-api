package taxonomy_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/taxonomy"
	"github.com/gorilla/mux"
)

var (
	testUser = &cohesioned.Profile{FullName: "Test User"}
)

func TestListViewHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/taxonomy", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	repo := new(fakes.FakeTaxonomyRepo)
	repo.ListReturns([]*cohesioned.Taxonomy{}, nil)

	renderer := fakes.FakeAdminDashboardRenderer
	handler := taxonomy.ListViewHandler(renderer, repo)
	profile := fakes.FakeProfile()

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &cohesioned.DashboardView{}
	dashboard.Set("profile", profile)
	dashboard.Set("list", []*cohesioned.Taxonomy{})

	expectedBody := fakes.RenderHTML(renderer, "taxonomy/list", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestAddHandler(t *testing.T) {
	expectedTaxonomy := cohesioned.NewTaxonomy("Test", 1, testUser)

	testJSON, err := expectedTaxonomy.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshall taxonomy to json %v", err)
	}

	req, err := http.NewRequest("POST", "/api/taxonomy", bytes.NewReader(testJSON))
	if err != nil {
		t.Fatalf("Failed to initialize create taxonomy request %v", err)
	}

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	repo := new(fakes.FakeTaxonomyRepo)
	repo.AddReturns(expectedTaxonomy, err)

	handler := taxonomy.AddHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data := struct {
		ID int64 `json:"id"`
	}{
		expectedTaxonomy.ID(),
	}

	expectedBody := fakes.RenderJSON(data)

	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected:\n%s\n\nActual:\n%s", string(expectedBody), rr.Body.String())
	}
}

func TestListChildrenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/taxonomy/1234/children", nil)
	if err != nil {
		t.Fatal(err)
	}

	fakeList := []*cohesioned.Taxonomy{
		cohesioned.NewTaxonomyWithParent("test-child-1", 1, 1234, testUser),
		cohesioned.NewTaxonomyWithParent("test-child-2", 2, 1234, testUser),
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

func TestFlattenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/taxonomy/flatten", nil)
	if err != nil {
		t.Fatal(err)
	}

	fakeList := []*cohesioned.Taxonomy{
		cohesioned.NewTaxonomy("Parent 1", 1234, testUser),
	}

	fakeFlattened := []*cohesioned.Taxonomy{
		cohesioned.NewTaxonomyWithParent("Parent 1 > Child 1", 5678, 1234, testUser),
		cohesioned.NewTaxonomyWithParent("Parent 1 > Child 2", 9101, 1234, testUser),
	}

	repo := new(fakes.FakeTaxonomyRepo)
	repo.ListReturns(fakeList, nil)
	repo.FlattenReturns(fakeFlattened, nil)

	rr := httptest.NewRecorder()
	handler := taxonomy.FlatListHandler(fakes.FakeRenderer, repo)

	router := mux.NewRouter()
	router.HandleFunc("/api/taxonomy/flatten", handler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedBody := fakes.RenderJSON(fakeFlattened)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("The expected json was not generated.\n\nExpected:\n%s\n\nActual:\n%s", string(expectedBody), rr.Body.String())
	}
}
