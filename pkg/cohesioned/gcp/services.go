package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

func NewDatastoreClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, projectID, option.WithServiceAccountFile("gcp-keyfile.json"))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new datastore client %v", err)
	}

	return client, nil
}

func NewStorageClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile("gcp-keyfile.json"))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new storage client %v", err)
	}

	return client, nil
}
