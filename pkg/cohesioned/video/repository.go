package video

import (
	"io"

	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	List() ([]*cohesioned.Video, error)
	Get(id int64) (*cohesioned.Video, error)
	Delete(id int64) error
	Add(video *cohesioned.Video) (*cohesioned.Video, error)
	Update(video *cohesioned.Video) (*cohesioned.Video, error)
	SetFile(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error)
}
