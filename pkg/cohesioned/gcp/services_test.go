package gcp_test

import (
	"net/url"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
	"github.com/cohesion-education/admin-api/pkg/cohesioned/gcp"
)

func TestCreateSignedURL(t *testing.T) {
	video := &cohesioned.Video{
		Key:     datastore.IDKey("Video", 1234, nil),
		Created: time.Now(),
		CreatedBy: &cohesioned.Profile{
			FullName: "Test User",
		},
		Title:             "Test Video",
		FileName:          "test.mp4",
		StorageBucket:     "test-bucket",
		StorageObjectName: "1234-test.mp4",
		TaxonomyID:        1,
	}

	gcpConfig, err := gcp.NewConfig("../../../gcp-keyfile.json")
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
