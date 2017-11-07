package billing

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	Save(p *cohesioned.PaymentDetails) (int64, error)
	Update(p *cohesioned.PaymentDetails) error
}
