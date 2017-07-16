package fakes

import "github.com/cohesion-education/api/pkg/cohesioned"

type FakeProfileRepo struct {
	profile *cohesioned.Profile
	err     error
}

func (r *FakeProfileRepo) SaveReturns(e error) {
	r.err = e
}

func (r *FakeProfileRepo) UpdateReturns(e error) {
	r.err = e
}

func (r *FakeProfileRepo) FindByEmailReturns(p *cohesioned.Profile, e error) {
	r.profile = p
	r.err = e
}

func (r *FakeProfileRepo) Save(p *cohesioned.Profile) error {
	return r.err
}

func (r *FakeProfileRepo) FindByEmail(email string) (*cohesioned.Profile, error) {
	return r.profile, r.err
}

func (r *FakeProfileRepo) Update(p *cohesioned.Profile) error {
	return r.err
}
