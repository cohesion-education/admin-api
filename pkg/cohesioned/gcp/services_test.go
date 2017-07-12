package gcp_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/gcp"
)

func TestCreateSignedURL(t *testing.T) {
	video := &cohesioned.Video{
		Title:             "Test Video",
		FileName:          "test.mp4",
		StorageBucket:     "test-bucket",
		StorageObjectName: "1234-test.mp4",
		TaxonomyID:        1,
	}

	video.SetID(1234)
	video.SetCreatedBy(&cohesioned.Profile{
		FullName: "Test User",
	})
	video.SetCreated(time.Now())

	gcpConfig, err := gcp.NewConfig("../../../testdata/test-gcp-keyfile.json")
	if err != nil {
		t.Fatalf("Failed to get gcp config %v", err)
	}

	result, err := gcp.CreateSignedURL(video, gcpConfig)
	if err != nil {
		t.Fatalf("Failed to create signed url %v", err)
	}

	expectedHostName := "storage.googleapis.com"
	signedURL, err := url.Parse(result)
	if err != nil {
		t.Fatalf("Parse result %s into net.URL %v", result, err)
	}

	if signedURL.Host != expectedHostName {
		t.Errorf("expected signed url host to be %s but was %s", expectedHostName, signedURL.Host)
	}

	q := signedURL.Query()

	if len(q.Get("GoogleAccessId")) == 0 {
		t.Errorf("Signed url did not have GoogleAccessId param")
	}

	if len(q.Get("Expires")) == 0 {
		t.Errorf("Signed url did not have Expires param")
	}

	if len(q.Get("Signature")) == 0 {
		t.Errorf("Signed url did not have Signature param")
	}
}
