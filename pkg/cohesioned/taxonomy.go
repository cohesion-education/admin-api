package cohesioned

import (
	"encoding/json"
	"time"
)

type Taxonomy struct {
	id        int64
	Created   time.Time `datastore:"created" json:"created"`
	CreatedBy *Profile  `datastore:"created_by" json:"created_by"`
	Name      string    `datastore:"name" json:"name"`
	ParentID  int64     `datastore:"parent_id" schema:"parent_id" json:"parent_id"`
}

func (t *Taxonomy) MarshalJSON() ([]byte, error) {
	type Alias Taxonomy
	return json.Marshal(&struct {
		ID int64 `json:"id"`
		*Alias
	}{
		ID:    t.id,
		Alias: (*Alias)(t),
	})
}

func (t *Taxonomy) ID() int64 {
	return t.id
}

func (t *Taxonomy) SetID(id int64) {
	t.id = id
}
