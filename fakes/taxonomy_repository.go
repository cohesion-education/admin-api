package fakes

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type FakeTaxonomyRepo struct {
	t         *cohesioned.Taxonomy
	list      []*cohesioned.Taxonomy
	flattened []*cohesioned.Taxonomy
	id        int64
	err       error
}

func (r *FakeTaxonomyRepo) ListReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) ListRecursiveReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) FindGradeByNameReturns(t *cohesioned.Taxonomy, err error) {
	r.t = t
	r.err = err
}

func (r *FakeTaxonomyRepo) GetReturns(t *cohesioned.Taxonomy, err error) {
	r.t = t
	r.err = err
}

func (r *FakeTaxonomyRepo) ListChildrenReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) ListChildrenRecursiveReturns(list []*cohesioned.Taxonomy, err error) {
	r.list = list
	r.err = err
}

func (r *FakeTaxonomyRepo) SaveReturns(id int64, err error) {
	r.id = id
	r.err = err
}

func (r *FakeTaxonomyRepo) UpdateReturns(err error) {
	r.err = err
}

func (r *FakeTaxonomyRepo) ReverseFlattenReturns(flattened *cohesioned.Taxonomy, err error) {
	r.t = flattened
	r.err = err
}

func (r *FakeTaxonomyRepo) FlattenReturns(flattened []*cohesioned.Taxonomy, err error) {
	r.flattened = flattened
	r.err = err
}

func (r *FakeTaxonomyRepo) FindGradeByName(name string) (*cohesioned.Taxonomy, error) {
	return r.t, r.err
}

func (r *FakeTaxonomyRepo) Get(id int64) (*cohesioned.Taxonomy, error) {
	return r.t, r.err
}
func (r *FakeTaxonomyRepo) List() ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
func (r *FakeTaxonomyRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}

func (r *FakeTaxonomyRepo) ListChildrenRecursive(parentID int64) ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}

func (r *FakeTaxonomyRepo) Save(t *cohesioned.Taxonomy) (int64, error) {
	return r.id, r.err
}
func (r *FakeTaxonomyRepo) Update(t *cohesioned.Taxonomy) error {
	return r.err
}
func (r *FakeTaxonomyRepo) Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error) {
	return r.flattened, r.err
}

func (r *FakeTaxonomyRepo) ReverseFlatten(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	return r.t, r.err
}

func (r *FakeTaxonomyRepo) ListRecursive() ([]*cohesioned.Taxonomy, error) {
	return r.list, r.err
}
