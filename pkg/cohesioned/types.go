package cohesioned

import (
	"time"
)

type Taxonomy struct {
	id       int64
	Name     string    `datastore:"name"`
	Created  time.Time `datastore:"created"`
	ParentID int64     `datastore:"parent_id" schema:"parent_id" json:"parent_id"`
	ChildIDs []int64   `datastore:"child_ids"`
}

func (t *Taxonomy) ID() int64 {
	return t.id
}

func (t *Taxonomy) SetID(id int64) {
	t.id = id
}
