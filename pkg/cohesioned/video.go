package cohesioned

import (
	"encoding/json"
	"time"
)

type Video struct {
	Auditable
	Validatable
	Title               string   `datastore:"title" json:"title"`
	FileName            string   `datastore:"file_name" json:"file_name"`
	StorageBucket       string   `datastore:"bucket" json:"bucket"`
	StorageObjectName   string   `datastore:"object_name" json:"object_name"`
	FlattenedTaxonomy   string   `datastore:"flattened_taxonomy" json:"flattened_taxonomy"`
	KeyTerms            []string `datastore:"key_terms" json:"key_terms"`
	StateStandards      []string `datastore:"state_standards" json:"state_standards"`
	CommonCoreStandards []string `datastore:"common_core_standards" json:"common_core_standards"`
	//TODO - Teacher, Related Videos, FAQs

	SignedURL string `json:"signed_url,omitempty"`
}

//NewVideo creates a Video with the Auditable fields initialized
func NewVideo(title, fileName string, id int64, createdBy *Profile) *Video {
	v := &Video{
		Title:    title,
		FileName: fileName,
	}

	v.GCPPersisted.id = id
	v.Auditable.Created = time.Now()
	v.Auditable.CreatedBy = createdBy

	return v
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

func (v *Video) Validate() bool {
	if len(v.Title) == 0 {
		v.AddValidationError("title", "title is required")
	}

	if len(v.FileName) == 0 {
		v.AddValidationError("file_name", "file_name is required")
	}

	if len(v.FlattenedTaxonomy) == -1 {
		v.AddValidationError("flattened_taxonomy", "flattened_taxonomy is required")
	}

	return len(v.ValidationErrors) == 0
}
