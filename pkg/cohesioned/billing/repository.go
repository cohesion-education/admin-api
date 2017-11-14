package billing

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	FindByCreatedByID(id int64) (*cohesioned.PaymentDetails, error)
	Save(p *cohesioned.PaymentDetails) (int64, error)
	Update(p *cohesioned.PaymentDetails) error
}
