package cohesioned

import (
	"encoding/json"
	"time"
)

type Taxonomy struct {
	Auditable
	Name     string      `datastore:"name" json:"name"`
	ParentID int64       `datastore:"parent_id" schema:"parent_id" json:"parent_id"`
	Children []*Taxonomy `datastore:"children" json:"children"`
}

//NewTaxonomy creates a Taxonomy with the Auditable fields initialized
func NewTaxonomy(name string, id int64, createdBy *Profile) *Taxonomy {
	return NewTaxonomyWithParent(name, id, 0, createdBy)
}

//NewTaxonomyWithParent creates a Taxonomy with the Auditable fields initialized and the parent ID set
func NewTaxonomyWithParent(name string, id int64, parentID int64, createdBy *Profile) *Taxonomy {
	t := &Taxonomy{
		Name:     name,
		ParentID: parentID,
	}

	t.GCPPersisted.id = id
	t.Auditable.Created = time.Now()
	t.Auditable.CreatedBy = createdBy

	return t
}

func (t *Taxonomy) MarshalJSON() ([]byte, error) {
	type Alias Taxonomy
	return json.Marshal(&struct {
		ID int64 `json:"id"`
		*Alias
	}{
		ID:    t.ID(),
		Alias: (*Alias)(t),
	})
}
