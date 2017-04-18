package cohesioned

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

type Video struct {
	id                int64
	Key               *datastore.Key `datastore:"__key__"`
	Created           time.Time      `datastore:"created" json:"created"`
	CreatedBy         *Profile       `datastore:"created_by" json:"created_by"`
	Title             string         `datastore:"title" json:"title"`
	FileName          string         `datastore:"file_name" json:"file_name"`
	StorageBucket     string         `datastore:"bucket" json:"bucket"`
	StorageObjectName string         `datastore:"object_name" json:"object_name"`
	TaxonomyID        int64          `datastore:"taxonomy_id" json:"taxonomy_id"`
	//TODO - Teacher, Tags, Related Videos, FAQs
}

func (v *Video) MarshalJSON() ([]byte, error) {
	type Alias Video
	return json.Marshal(&struct {
		ID int64 `json:"id"`
		*Alias
	}{
		ID:    v.ID(),
		Alias: (*Alias)(v),
	})
}

func (v *Video) ID() int64 {
	if v.Key.ID != 0 {
		return v.Key.ID
	}

	if v.id != 0 {
		return v.id
	}

	panic("Unable to load video ID")
}

func (v *Video) SetID(id int64) {
	v.id = id
}

func (v *Video) SignedURL() (string, error) {
	bucket := v.StorageBucket
	filename := v.StorageObjectName
	method := "GET"
	//TODO - make this configurable
	expires := time.Now().Add(time.Hour * 1)

	//TODO - load key from env var
	privateKey, err := ioutil.ReadFile("key.pem")
	if err != nil {
		return "", fmt.Errorf("Failed to read key.pem %v", err)
	}

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		//TODO - load service account id from env var
		GoogleAccessID: "cohesion-storage-admin@cohesion-education-164614.iam.gserviceaccount.com",
		PrivateKey:     privateKey,
		Method:         method,
		Expires:        expires,
	})

	if err != nil {
		return "", fmt.Errorf("Failed to sign url %v", err)
	}

	return url, nil
}
