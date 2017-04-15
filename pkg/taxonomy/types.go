package taxonomy

import (
	"time"
)

type Taxonomy struct {
	repo     Repo
	id       int64
	Name     string    `datastore:"name"`
	Created  time.Time `datastore:"created"`
	ParentID int64     `datastore:"parent_id" schema:"parent_id" json:"parent_id"`
	ChildIDs []int64   `datastore:"child_ids"`
}

func (t *Taxonomy) ID() int64 {
	return t.id
}

func (t *Taxonomy) Children() []*Taxonomy {
	children, err := t.repo.ListChildren(t.id)
	if err != nil {
		return nil
	}
	return children
}

func (t *Taxonomy) Parent() *Taxonomy {
	if t.ParentID != -1 {

	}

	return nil
}
