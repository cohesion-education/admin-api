package taxonomy

import (
	"context"
	"fmt"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"

	"cloud.google.com/go/datastore"
)

//Repo for interacting with the persistent store for the Taxonomy type
type Repo interface {
	Get(id int64) (*cohesioned.Taxonomy, error)
	List() ([]*cohesioned.Taxonomy, error)
	ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error)
	Add(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error)
	Save(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error)
	Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error)
}

//NewGCPDatastoreRepo implementation of taxonomy.Repo
func NewGCPDatastoreRepo(client *datastore.Client) Repo {
	return &gcpDatastoreRepo{
		client: client,
		ctx:    context.TODO(),
	}
}

type gcpDatastoreRepo struct {
	client *datastore.Client
	ctx    context.Context
}

func (r *gcpDatastoreRepo) Get(id int64) (*cohesioned.Taxonomy, error) {
	taxonomy := &cohesioned.Taxonomy{}

	key := datastore.IDKey("Taxonomy", id, nil)
	err := r.client.Get(r.ctx, key, taxonomy)

	if err == datastore.ErrInvalidEntityType {
		return nil, fmt.Errorf("%d returns an invalid entity type %v", id, err)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to get Taxonomy by id %d %v", id, err)
	}

	return taxonomy, nil
}

func (repo *gcpDatastoreRepo) List() ([]*cohesioned.Taxonomy, error) {
	return repo.ListChildren(0)
}

func (repo *gcpDatastoreRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	var list []*cohesioned.Taxonomy

	q := datastore.NewQuery("Taxonomy").Filter("parent_id=", parentID)
	if _, err := repo.client.GetAll(context.Background(), q, &list); err != nil {
		return nil, fmt.Errorf("Failed to get taxonomy list from Cloud Datastore %v", err)
	}

	return list, nil
}

func (repo *gcpDatastoreRepo) Add(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	t.Created = time.Now()

	key := datastore.IncompleteKey("Taxonomy", nil)
	key, err := repo.client.Put(repo.ctx, key, t)
	if err != nil {
		return t, fmt.Errorf("Failed to save Taxonomy %v", err)
	}

	t.Key = key
	return t, nil
}

func (repo *gcpDatastoreRepo) Save(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	if t.ID() == -1 {
		t.Created = time.Now()
		t.Key = datastore.IncompleteKey("Taxonomy", nil)
	}

	t.Updated = time.Now()
	key, err := repo.client.Put(repo.ctx, t.Key, t)
	if err != nil {
		return nil, fmt.Errorf("Failed to save Taxonomy %v", err)
	}

	t.Key = key
	return t, nil
}

func (repo *gcpDatastoreRepo) Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error) {
	flattened := []*cohesioned.Taxonomy{}
	if t == nil {
		return flattened, nil
	}

	children, err := repo.ListChildren(t.ID())
	if err != nil {
		return flattened, fmt.Errorf("Failed to get children for %s %v\n", t.Name, err)
	}

	if len(children) == 0 {
		fmt.Printf("Flattened: %s\n", t.Name)
		flattened = append(flattened, t)
		return flattened, nil
	}

	for _, child := range children {
		child.Name = fmt.Sprintf("%s > %s", t.Name, child.Name)
		flattenedChildren, err := repo.Flatten(child)
		if err != nil {
			return flattened, fmt.Errorf("Failed to flatten children of %s %v", child.Name, err)
		}
		flattened = append(flattened, flattenedChildren...)
	}

	return flattened, nil
}
