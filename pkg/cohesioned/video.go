package cohesioned

import (
	"time"
)

type Video struct {
	Validatable
	ID                  int64     `json:"id"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
	CreatedByID         int64     `json:"created_by_id"`
	CreatedBy           *Profile  `json:"created_by"`
	UpdatedByID         int64     `json:"updated_by_id"`
	UpdatedBy           *Profile  `json:"updated_by"`
	Title               string    `json:"title"`
	TaxonomyID          int64     `json:"taxonomy_id"`
	Taxonomy            *Taxonomy `json:"taxonomy"`
	KeyTerms            []string  `json:"key_terms,omitempty"`
	StateStandards      []string  `json:"state_standards,omitempty"`
	CommonCoreStandards []string  `json:"common_core_standards,omitempty"`
	FileName            string    `json:"file_name"`
	FileType            string    `json:"file_type"`
	FileSize            int64     `json:"file_size"`
	StorageBucket       string    `json:"bucket"`
	StorageObjectName   string    `json:"object_name"`
	SignedURL           string    `json:"signed_url,omitempty"`
	ThumbnailURL        string    `json:thumbnail_url`
	//TODO - Teacher, Related Videos, FAQs
}

//NewVideo creates a Video with the Auditable fields initialized
func NewVideo(title, fileName string, id int64, createdBy *Profile) *Video {
	v := &Video{
		Title:       title,
		FileName:    fileName,
		ID:          id,
		Created:     time.Now(),
		CreatedByID: createdBy.ID,
	}

	return v
}

func (v *Video) Validate() bool {
	if len(v.Title) == 0 {
		v.AddValidationError("title", "title is required")
	}

	if len(v.FileName) == 0 {
		v.AddValidationError("file_name", "file_name is required")
	}

	if v.TaxonomyID == 0 {
		v.AddValidationError("taxonomy_id", "taxonomy_id is required")
	}

	return len(v.ValidationErrors) == 0
}
