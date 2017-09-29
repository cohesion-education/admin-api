package taxonomy_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/taxonomy"
	"github.com/gorilla/mux"
)

var (
	testUser = &cohesioned.Profile{FullName: "Test User"}
)

func TestAddHandler(t *testing.T) {
	testTaxonomy := cohesioned.NewTaxonomy("Test", testUser)
	testTaxonomy.ID = 5
	testTaxonomy.ParentID = 1
	testJSON, err := testTaxonomy.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshall taxonomy to json %v", err)
	}

	req, err := http.NewRequest("POST", "/api/taxonomy", bytes.NewReader(testJSON))
	if err != nil {
		t.Fatalf("Failed to initialize create taxonomy request %v", err)
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, cohesioned.CurrentUserKey, testUser)
	req = req.WithContext(ctx)

	repo := new(fakes.FakeTaxonomyRepo)
	repo.SaveReturns(testTaxonomy.ID, err)

	handler := taxonomy.AddHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	actualResp := &taxonomy.TaxonomyResponse{}
	decoder := json.NewDecoder(rr.Body)
	if err = decoder.Decode(&actualResp); err != nil {
		t.Errorf("Failed to unmarshall response json to APIResponse: %v", err)
	}

	if actualResp.Name != testTaxonomy.Name {
		t.Errorf("expected name %s - got %s", testTaxonomy.Name, actualResp.Name)
	}

	if actualResp.CreatedBy != testTaxonomy.CreatedBy {
		t.Errorf("expected created by %d - got %d", testTaxonomy.CreatedBy, actualResp.CreatedBy)
	}

	if actualResp.ParentID != testTaxonomy.ParentID {
		t.Errorf("expected parent id %d - got %d", testTaxonomy.ParentID, actualResp.ParentID)
	}

	if len(actualResp.Children) != len(testTaxonomy.Children) {
		t.Errorf("expected children length of %d - got %d", len(testTaxonomy.Children), len(actualResp.Children))
	}
}

func TestListChildrenHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/taxonomy/1234/children", nil)
	if err != nil {
		t.Fatal(err)
	}

	fakeList := []*cohesioned.Taxonomy{
		cohesioned.NewTaxonomyWithParent("test-child-1", 1234, testUser),
		cohesioned.NewTaxonomyWithParent("test-child-2", 1234, testUser),
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
		cohesioned.NewTaxonomy("Parent 1", testUser),
	}

	fakeFlattened := []*cohesioned.Taxonomy{
		cohesioned.NewTaxonomyWithParent("Parent 1 > Child 1", 1234, testUser),
		cohesioned.NewTaxonomyWithParent("Parent 1 > Child 2", 1234, testUser),
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
