package cohesioned

import (
	"encoding/json"
	"time"
)

type Video struct {
	Auditable
	Validatable
	ID                  int64    `json:"id"`
	Title               string   `datastore:"title" json:"title"`
	FlattenedTaxonomy   string   `datastore:"flattened_taxonomy" json:"flattened_taxonomy"`
	KeyTerms            []string `datastore:"key_terms" json:"key_terms,omitempty"`
	StateStandards      []string `datastore:"state_standards" json:"state_standards,omitempty"`
	CommonCoreStandards []string `datastore:"common_core_standards" json:"common_core_standards,omitempty"`
	FileName            string   `datastore:"file_name" json:"file_name"`
	FileType            string   `datastore:"file_type" json:"file_type"`
	FileSize            int64    `datastore:"file_size" json:"file_size"`
	StorageBucket       string   `datastore:"bucket" json:"bucket"`
	StorageObjectName   string   `datastore:"object_name" json:"object_name"`
	SignedURL           string   `json:"signed_url,omitempty"`
	//TODO - Teacher, Related Videos, FAQs
}

//NewVideo creates a Video with the Auditable fields initialized
func NewVideo(title, fileName string, id int64, createdBy *Profile) *Video {
	v := &Video{
		Title:    title,
		FileName: fileName,
		ID:       id,
	}

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
		ID:    v.ID,
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
