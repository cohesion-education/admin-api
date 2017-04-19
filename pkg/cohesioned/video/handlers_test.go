package video_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-education/admin-api/fakes"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/common"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/config"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/video"
)

func TestListHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/video", nil)
	if err != nil {
		t.Fatal(err)
	}

	repo := new(fakes.FakeVideoRepo)
	repo.ListReturns([]*cohesioned.Video{}, nil)

	profile := fakes.FakeProfile()
	ctx := req.Context()
	ctx = context.WithValue(ctx, config.CurrentUserKey, profile)
	req = req.WithContext(ctx)

	handler := video.ListHandler(fakes.FakeRenderer, repo)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	dashboard := &common.DashboardView{}
	dashboard.Set("profile", profile)
	dashboard.Set("list", []*cohesioned.Video{})

	expectedBody := fakes.RenderHTML("video/list", dashboard)
	if bytes.Compare(expectedBody, rr.Body.Bytes()) != 0 {
		t.Errorf("HTML response was not generated as expected. Expected:\n\n%s\n\nActual:\n\n%s", string(expectedBody), rr.Body.String())
	}
}
