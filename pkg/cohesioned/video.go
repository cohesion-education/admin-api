package cohesioned

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/datastore"
)

type Video struct {
	id                int64
	Key               *datastore.Key `datastore:"__key__"`
	Created           time.Time      `datastore:"created" json:"created"`
	Updated           time.Time      `datastore:"updated" json:"updated"`
	CreatedBy         *Profile       `datastore:"created_by" json:"created_by"`
	UpdatedBy         *Profile       `datastore:"updated_by" json:"updated_by"`
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
