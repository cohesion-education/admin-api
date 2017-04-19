package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

const gcpKeyfileName = "gcp-keyfile.json"

var (
	privateKey     []byte
	googleAccessID string
)

func init() {
	keyfileBytes, err := ioutil.ReadFile(gcpKeyfileName)
	if err != nil {
		panic(fmt.Errorf("Failed to read %s %v", gcpKeyfileName, err))
	}

	keyfile := make(map[string]string)
	if err := json.Unmarshal(keyfileBytes, &keyfile); err != nil {
		panic(fmt.Errorf("Failed to unmarshall json from keyfile %s %v", gcpKeyfileName, err))
	}

	privateKey = []byte(keyfile["private_key"])
	googleAccessID = keyfile["client_email"]
}

func NewDatastoreClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, projectID, option.WithServiceAccountFile(gcpKeyfileName))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new datastore client %v", err)
	}

	return client, nil
}

func NewStorageClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(gcpKeyfileName))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new storage client %v", err)
	}

	return client, nil
}

func CreateSignedURL(v *cohesioned.Video) (string, error) {
	bucket := v.StorageBucket
	filename := v.StorageObjectName
	method := "GET"
	expires := time.Now().Add(time.Hour * 1)

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		GoogleAccessID: googleAccessID,
		PrivateKey:     privateKey,
		Method:         method,
		Expires:        expires,
	})

	if err != nil {
		return "", fmt.Errorf("Failed to sign url %v", err)
	}

	return url, nil
}
