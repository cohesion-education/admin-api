package taxonomy

import (
	"github.com/cohesion-education/api/pkg/cohesioned"
)

//Repo for interacting with the persistent store for the Taxonomy type
type Repo interface {
	Get(id int64) (*cohesioned.Taxonomy, error)
	FindGradeByName(name string) (*cohesioned.Taxonomy, error)
	List() ([]*cohesioned.Taxonomy, error)
	ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error)
	Save(t *cohesioned.Taxonomy) (int64, error)
	Update(t *cohesioned.Taxonomy) error
	Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error)
	ReverseFlatten(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error)
	ListRecursive() ([]*cohesioned.Taxonomy, error)
	ListChildrenRecursive(parentID int64) ([]*cohesioned.Taxonomy, error)
}
