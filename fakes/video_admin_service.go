package fakes

import (
	"context"
	"io"

	"github.com/cohesion-education/api/pkg/cohesioned"
)

type FakeVideoAdminService struct {
	v                       *cohesioned.Video
	err                     error
	list                    []*cohesioned.Video
	videosByGrade map[string][]*cohesioned.Video
	videosBySubject map[string][]*cohesioned.Video
}

func (s *FakeVideoAdminService) FindByTaxonomyIDReturns(list []*cohesioned.Video, err error) {
	s.list = list
	s.err = err
}

func (s *FakeVideoAdminService) FindByGradeReturns(videosByGrade map[string][]*cohesioned.Video, err error) {
	s.videosByGrade = videosByGrade
	s.err = err
}

func (s *FakeVideoAdminService) FindBySubjectReturns(videosBySubject map[string][]*cohesioned.Video, err error) {
	s.videosBySubject = videosBySubject
	s.err = err
}

func (s *FakeVideoAdminService) ListReturns(list []*cohesioned.Video, err error) {
	s.list = list
	s.err = err
}
func (s *FakeVideoAdminService) GetReturns(v *cohesioned.Video, err error) {
	s.v = v
	s.err = err
}
func (s *FakeVideoAdminService) GetWithSignedURLReturns(v *cohesioned.Video, err error) {
	s.v = v
	s.err = err
}

func (s *FakeVideoAdminService) DeleteReturns(err error) {
	s.err = err
}
func (s *FakeVideoAdminService) SaveReturns(err error) {
	s.err = err
}
func (s *FakeVideoAdminService) UpdateReturns(err error) {
	s.err = err
}
func (s *FakeVideoAdminService) SetFileReturns(err error) {
	s.err = err
}

func (s *FakeVideoAdminService) List() ([]*cohesioned.Video, error) {
	return s.list, s.err
}
func (s *FakeVideoAdminService) FindByTaxonomyID(taxonomyID int64) ([]*cohesioned.Video, error) {
	return s.list, s.err
}

func (s *FakeVideoAdminService) FindByGrade(gradeName string) (map[string][]*cohesioned.Video, error) {
	return s.videosByGrade, s.err
}

func (s *FakeVideoAdminService) FindBySubject(gradeName, subjectName string) (map[string][]*cohesioned.Video, error) {
	return s.videosBySubject, s.err
}

func (s *FakeVideoAdminService) Get(id int64) (*cohesioned.Video, error) {
	return s.v, s.err
}
func (s *FakeVideoAdminService) GetWithSignedURL(id int64) (*cohesioned.Video, error) {
	return s.v, s.err
}
func (s *FakeVideoAdminService) Delete(id int64) error {
	return s.err
}
func (s *FakeVideoAdminService) Save(ctx context.Context, video *cohesioned.Video) error {
	return s.err
}
func (s *FakeVideoAdminService) Update(ctx context.Context, video *cohesioned.Video) error {
	return s.err
}
func (s *FakeVideoAdminService) SetFile(ctx context.Context, fileReader io.Reader, video *cohesioned.Video) error {
	return s.err
}
