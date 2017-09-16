package cohesioned

import (
	"encoding/json"
	"time"
)

type Video struct {
	Auditable
	Title    string `datastore:"title" json:"title"`
	FileName string `datastore:"file_name" json:"file_name"`
	// FileReader        io.Reader `datastore:"-" json:"-"`
	StorageBucket       string   `datastore:"bucket" json:"bucket"`
	StorageObjectName   string   `datastore:"object_name" json:"object_name"`
	TaxonomyID          int64    `datastore:"taxonomy_id" json:"taxonomy_id"`
	KeyTerms            []string `datastore:"key_terms" json:"key_terms"`
	StateStandards      []string `datastore:"state_standards" json:"state_standards"`
	CommonCoreStandards []string `datastore:"common_core_standards" json:"common_core_standards"`
	//TODO - Teacher, Tags, Related Videos, FAQs
}

//NewVideo creates a Video with the Auditable fields initialized
func NewVideo(title, fileName, storageBucket, storageObjectName string, id, taxonomyID int64, createdBy *Profile) *Video {
	v := &Video{
		Title:             title,
		FileName:          fileName,
		StorageBucket:     storageBucket,
		StorageObjectName: storageObjectName,
		TaxonomyID:        taxonomyID,
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
