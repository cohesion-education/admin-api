package taxonomy

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

//Repo for interacting with the persistent store for the Taxonomy type
type Repo interface {
	List() ([]*Taxonomy, error)
	ListChildren(parentID int64) ([]*Taxonomy, error)
	Add(t *Taxonomy) (*datastore.Key, error)
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

func (repo *gcpDatastoreRepo) List() ([]*Taxonomy, error) {
	return repo.ListChildren(0)
}

func (repo *gcpDatastoreRepo) ListChildren(parentID int64) ([]*Taxonomy, error) {
	var list []*Taxonomy

	q := datastore.NewQuery("Taxonomy").Filter("parent_id=", parentID)
	keys, err := repo.client.GetAll(context.Background(), q, &list)
	if err != nil {
		return nil, fmt.Errorf("Failed to get taxonomy list from Cloud Datastore %v", err)
	}

	for i, key := range keys {
		list[i].id = key.ID
		list[i].repo = repo
	}

	return list, nil
}

func (repo *gcpDatastoreRepo) Add(t *Taxonomy) (*datastore.Key, error) {
	t.Created = time.Now()

	key := datastore.IncompleteKey("Taxonomy", nil)
	return repo.client.Put(repo.ctx, key, t)
}
