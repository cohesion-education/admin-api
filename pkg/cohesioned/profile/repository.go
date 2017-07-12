package profile

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	Save(p *cohesioned.Profile) error
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

func (repo *gcpDatastoreRepo) Save(p *cohesioned.Profile) error {
	if p.ID() == -1 {
		p.Created = time.Now()
		key := datastore.IncompleteKey("Profile", nil)
		key, err := repo.client.Put(repo.ctx, key, p)
		if err != nil {
			return fmt.Errorf("Failed to save Profile %v", err)
		}

		p.Key = key
		return nil
	}

	p.Updated = time.Now()
	_, err := repo.client.Put(repo.ctx, p.Key, p)
	if err != nil {
		return fmt.Errorf("Failed to update Profile %v", err)
	}

	return nil
}
