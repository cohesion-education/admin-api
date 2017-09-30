package video

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	List() ([]*cohesioned.Video, error)
	Get(id int64) (*cohesioned.Video, error)
	Delete(id int64) error
	Save(video *cohesioned.Video) (int64, error)
	Update(video *cohesioned.Video) error
}
