package profile

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"
	"github.com/cohesion-education/api/pkg/cohesioned"
)

type Repo interface {
	FindByEmail(email string) (*cohesioned.Profile, error)
	Save(p *cohesioned.Profile) error
	Update(p *cohesioned.Profile) error
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

func (repo *gcpDatastoreRepo) Update(p *cohesioned.Profile) error {
	p.Updated = time.Now()
	_, err := repo.client.Put(repo.ctx, p.Key, p)
	if err != nil {
		return fmt.Errorf("Failed to update Profile %v", err)
	}

	return nil
}

func (repo *gcpDatastoreRepo) FindByEmail(email string) (*cohesioned.Profile, error) {
	q := datastore.NewQuery("Profile").Filter("Email =", email).DistinctOn("Email")

	for t := repo.client.Run(repo.ctx, q); ; {
		var p cohesioned.Profile
		key, err := t.Next(&p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("error getting next; %v\n", err)
		}

		return &p, nil
		fmt.Printf("Key: %v - Profile: %v\n", key, p)
	}

	return nil, nil
}
