package cohesioned

import (
	"encoding/json"
	"time"
)

type Taxonomy struct {
	ID        int64       `json:"id"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
	CreatedBy int64       `json:"created_by"`
	UpdatedBy int64       `json:"updated_by"`
	Name      string      `json:"name"`
	ParentID  int64       `schema:"parent_id" json:"parent_id"`
	Children  []*Taxonomy `json:"children"`
}

//NewTaxonomy creates a Taxonomy with the Auditable fields initialized
func NewTaxonomy(name string, createdBy *Profile) *Taxonomy {
	return &Taxonomy{
		Name:      name,
		Created:   time.Now(),
		CreatedBy: createdBy.ID,
	}
}

//NewTaxonomyWithParent creates a Taxonomy with the Auditable fields initialized and the parent ID set
func NewTaxonomyWithParent(name string, parentID int64, createdBy *Profile) *Taxonomy {
	return &Taxonomy{
		Name:      name,
		ParentID:  parentID,
		Created:   time.Now(),
		CreatedBy: createdBy.ID,
	}
}

func (t *Taxonomy) MarshalJSON() ([]byte, error) {
	type Alias Taxonomy
	return json.Marshal(&struct {
		ID int64 `json:"id"`
		*Alias
	}{
		ID:    t.ID,
		Alias: (*Alias)(t),
	})
}
