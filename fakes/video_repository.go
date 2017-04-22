package fakes

import (
	"io"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"
)

type FakeVideoRepo struct {
	v    *cohesioned.Video
	list []*cohesioned.Video
	err  error
}

func (r *FakeVideoRepo) ListReturns(list []*cohesioned.Video, err error) {
	r.list = list
	r.err = err
}

func (r *FakeVideoRepo) GetReturns(v *cohesioned.Video, err error) {
	r.v = v
	r.err = err
}

func (r *FakeVideoRepo) AddReturns(v *cohesioned.Video, err error) {
	r.v = v
	r.err = err
}

func (r *FakeVideoRepo) UpdateReturns(v *cohesioned.Video, err error) {
	r.v = v
	r.err = err
}

func (r *FakeVideoRepo) List() ([]*cohesioned.Video, error) {
	return r.list, r.err
}
func (r *FakeVideoRepo) Get(id int64) (*cohesioned.Video, error) {
	return r.v, r.err
}
func (r *FakeVideoRepo) Add(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	return r.v, r.err
}

func (r *FakeVideoRepo) Update(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	return r.v, r.err
}
