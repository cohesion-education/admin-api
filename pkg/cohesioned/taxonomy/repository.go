package taxonomy

import (
	"context"
	"fmt"
	"time"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"

	"cloud.google.com/go/datastore"
)

//Repo for interacting with the persistent store for the Taxonomy type
type Repo interface {
	List() ([]*cohesioned.Taxonomy, error)
	ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error)
	Add(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error)
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
