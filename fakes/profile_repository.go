package fakes

import "github.com/cohesion-education/api/pkg/cohesioned"

type FakeProfileRepo struct {
	id      int64
	profile *cohesioned.Profile
	list    []*cohesioned.Profile
	err     error
}

func (r *FakeProfileRepo) SaveReturns(id int64, e error) {
	r.id = id
	r.err = e
}

func (r *FakeProfileRepo) ListReturns(list []*cohesioned.Profile, e error) {
	r.list = list
	r.err = e
}

func (r *FakeProfileRepo) UpdateReturns(e error) {
	r.err = e
}

func (r *FakeProfileRepo) FindByEmailReturns(p *cohesioned.Profile, e error) {
	r.profile = p
	r.err = e
}

func (r *FakeProfileRepo) Save(p *cohesioned.Profile) (int64, error) {
	return r.id, r.err
}

func (r *FakeProfileRepo) FindByEmail(email string) (*cohesioned.Profile, error) {
	return r.profile, r.err
}

func (r *FakeProfileRepo) Update(p *cohesioned.Profile) error {
	return r.err
}

func (r *FakeProfileRepo) List() ([]*cohesioned.Profile, error) {
	return r.list, r.err
}
