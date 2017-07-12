package homepage

import (
	"context"
	"fmt"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"

	"cloud.google.com/go/datastore"
)

//Repo for interacting with the persistent store for the Homepage type
type Repo interface {
	Save(h *cohesioned.Homepage) (*cohesioned.Homepage, error)
	Get() (*cohesioned.Homepage, error)
}

//NewGCPDatastoreRepo implementation of homepage.Repo
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

func (repo *gcpDatastoreRepo) Get() (*cohesioned.Homepage, error) {
	var homepages []*cohesioned.Homepage

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery("Homepage").Limit(1)
	keys, err := repo.client.GetAll(repo.ctx, query, &homepages)
	if err != nil {
		return nil, err
	}

	if len(homepages) != 1 {
		return nil, nil
	}

	homepage := homepages[0]
	homepage.SetID(keys[0].ID)
	return homepage, nil
}

func (repo *gcpDatastoreRepo) Save(h *cohesioned.Homepage) (*cohesioned.Homepage, error) {
	if h.ID() == -1 {
		h.Created = time.Now()
		key := datastore.IncompleteKey("Homepage", nil)
		key, err := repo.client.Put(repo.ctx, key, h)
		if err != nil {
			return h, fmt.Errorf("Failed to save Homepage %v", err)
		}

		h.Key = key
		return h, nil
	}

	h.Updated = time.Now()
	_, err := repo.client.Put(repo.ctx, h.Key, h)
	if err != nil {
		return h, fmt.Errorf("Failed to update Homepage %v", err)
	}

	return nil, nil
}
