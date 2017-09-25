package fakes

import (
	"io"

	"github.com/cohesion-education/api/pkg/cohesioned"
)

type FakeVideoRepo struct {
	v    *cohesioned.Video
	id   int64
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

func (r *FakeVideoRepo) DeleteReturns(err error) {
	r.err = err
}

func (r *FakeVideoRepo) SaveReturns(id int64, err error) {
	r.id = id
	r.err = err
}

func (r *FakeVideoRepo) UpdateReturns(err error) {
	r.err = err
}

func (r *FakeVideoRepo) SetFileReturns(v *cohesioned.Video, err error) {
	r.v = v
	r.err = err
}

func (r *FakeVideoRepo) List() ([]*cohesioned.Video, error) {
	return r.list, r.err
}
func (r *FakeVideoRepo) Get(id int64) (*cohesioned.Video, error) {
	return r.v, r.err
}
func (r *FakeVideoRepo) Delete(id int64) error {
	return r.err
}
func (r *FakeVideoRepo) Save(video *cohesioned.Video) (int64, error) {
	return r.id, r.err
}

func (r *FakeVideoRepo) Update(video *cohesioned.Video) error {
	return r.err
}

func (r *FakeVideoRepo) SetFile(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	return r.v, r.err
}
