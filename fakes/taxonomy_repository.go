package fakes

import (
	"github.com/cohesion-education/admin-api/pkg/cohesioned"
)

type FakeTaxonomyRepo struct {
	t         *cohesioned.Taxonomy
	list      []*cohesioned.Taxonomy
	flattened []*cohesioned.Taxonomy
	err       error
}

func (r *FakeTaxonomyRepo) ListReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) ListChildrenReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) AddReturns(t *cohesioned.Taxonomy, err error) {
	r.t = t
	r.err = err
}

func (r *FakeTaxonomyRepo) FlattenReturns(flattened []*cohesioned.Taxonomy, err error) {
	r.flattened = flattened
	r.err = err
}

func (r *FakeTaxonomyRepo) List() ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) Add(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	return r.t, r.err
}

func (r *FakeTaxonomyRepo) Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error) {
	return r.flattened, r.err
}
