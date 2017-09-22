package profile

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	FindByEmail(email string) (*cohesioned.Profile, error)
	Save(p *cohesioned.Profile) (int64, error)
	Update(p *cohesioned.Profile) error
}
