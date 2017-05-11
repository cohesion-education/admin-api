package fakes

import "github.com/cohesion-education/admin-api/pkg/cohesioned"

type FakeProfileRepo struct {
	err error
}

func (r *FakeProfileRepo) SaveReturns(err error) {
	r.err = err
}

func (r *FakeProfileRepo) Save(p *cohesioned.Profile) error {
	return r.err
}
