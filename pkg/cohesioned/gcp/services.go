package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

func NewDatastoreClient(ctx context.Context, cfg *Config) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, cfg.projectID, option.WithServiceAccountFile(cfg.keyfileLocation))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new datastore client %v", err)
	}

	return client, nil
}

func NewStorageClient(ctx context.Context, cfg *Config) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(cfg.keyfileLocation))
	if err != nil {
		return nil, fmt.Errorf("Failed to make new storage client %v", err)
	}

	return client, nil
}

func CreateSignedURL(v *cohesioned.Video, cfg *Config) (string, error) {
	bucket := v.StorageBucket
	filename := v.StorageObjectName
	method := "GET"
	expires := time.Now().Add(time.Hour * 1)

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		GoogleAccessID: cfg.googleAccessID,
		PrivateKey:     cfg.privateKey,
		Method:         method,
		Expires:        expires,
	})

	if err != nil {
		return "", fmt.Errorf("Failed to sign url %v", err)
	}

	return url, nil
}
