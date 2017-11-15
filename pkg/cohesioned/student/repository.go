package student

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	FindByUserID(parentID int64) ([]*cohesioned.Student, error)
	List() ([]*cohesioned.Student, error)
	Save(s *cohesioned.Student) (int64, error)
	Update(s *cohesioned.Student) error
	Delete(id int64) error
}
