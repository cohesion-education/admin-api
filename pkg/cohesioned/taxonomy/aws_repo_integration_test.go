//+build integration

package taxonomy_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/taxonomy"
	"github.com/cohesion-education/api/testutils"
)

var (
	repo      taxonomy.Repo
	emptyTime = time.Time{}
)

func TestMain(m *testing.M) {
	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		panic("Unexpected error initializing AwsConfig: " + err.Error())
	}

	db, err := awsConfig.DialRDS()
	if err != nil {
		panic("Failed to connect to db " + err.Error())
	}

	repo = taxonomy.NewAwsRepo(db)

	if err := testutils.SetupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	testResult := m.Run()

	if err := testutils.CleanupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(testResult)
}

func TestRepoGet(t *testing.T) {
	taxonomy, err := repo.Get(1)

	if err != nil {
		t.Errorf("Failed to find taxonomy by id: %v", err)
	}

	if taxonomy == nil {
		t.Errorf("taxonomy by id returned nil")
	}

	if len(taxonomy.Name) == 0 {
		t.Error("taxonomy name is empty")
	}

	if taxonomy.Created == emptyTime {
		t.Error("taxonomy created is empty")
	}

	if taxonomy.CreatedBy == 0 {
		t.Error("Taxonomy created by is empty")
	}
}

func TestRepoList(t *testing.T) {
	list, err := repo.List()
	if err != nil {
		t.Errorf("Failed to list taxonomy: %v", err)
	}

	if len(list) == 0 {
		t.Errorf("taxonomy list is empty")
	}

	for _, taxonomy := range list {
		if len(taxonomy.Name) == 0 {
			t.Error("taxonomy name is empty")
		}

		if taxonomy.Created == emptyTime {
			t.Error("taxonomy created is empty")
		}

		if taxonomy.CreatedBy == 0 {
			t.Error("Taxonomy created by is empty")
		}
	}
}

func TestRepoListChildren(t *testing.T) {
	list, err := repo.ListChildren(2)
	if err != nil {
		t.Errorf("Failed to list taxonomy children: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("taxonomy list of children not the expected length")
	}

	for _, taxonomy := range list {
		if len(taxonomy.Name) == 0 {
			t.Error("taxonomy name is empty")
		}

		if taxonomy.Created == emptyTime {
			t.Error("taxonomy created is empty")
		}

		if taxonomy.CreatedBy == 0 {
			t.Error("Taxonomy created by is empty")
		}
	}
}

func TestRepoFlatten(t *testing.T) {
	profile := fakes.FakeProfile()
	taxonomy := cohesioned.NewTaxonomy("test-parent", profile)
	taxonomy.ID = 2

	list, err := repo.Flatten(taxonomy)
	if err != nil {
		t.Errorf("Failed to flatten taxonomy: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("flattened taxonomy not the expected length")
	}

	for _, taxonomy := range list {
		if len(taxonomy.Name) == 0 {
			t.Error("flattened taxonomy name is empty")
		}

		if taxonomy.Name != "test-parent > test-child" {
			t.Errorf("flattened taxonomy name not the expected value. Expected: %s Got: %s", "test-parent > test-child", taxonomy.Name)
		}

		if taxonomy.ID != 3 {
			t.Errorf("flattened taxonomy id not the expected value. Expected: %d Got: %d", 3, taxonomy.ID)
		}

		if taxonomy.Created == emptyTime {
			t.Error("flattened taxonomy created is empty")
		}

		if taxonomy.CreatedBy == 0 {
			t.Error("flattened Taxonomy created by is empty")
		}
	}
}

func TestRepoSaveWithNilParent(t *testing.T) {
	profile := fakes.FakeProfile()
	taxonomy := cohesioned.NewTaxonomy("test", profile)
	id, err := repo.Save(taxonomy)

	if err != nil {
		t.Errorf("Failed to save taxonomy: %v", err)
	}

	if id == 0 {
		t.Errorf("taxonomy ID was zero - expected db to generate a taxonomy id")
	}
}

func TestRepoSaveWithParent(t *testing.T) {
	profile := fakes.FakeProfile()
	taxonomy := cohesioned.NewTaxonomy("test child", profile)
	taxonomy.ParentID = 2
	id, err := repo.Save(taxonomy)

	if err != nil {
		t.Errorf("Failed to save taxonomy: %v", err)
	}

	if id == 0 {
		t.Errorf("taxonomy ID was zero - expected db to generate a taxonomy id")
	}
}

func TestRepoUpdate(t *testing.T) {
	profile := fakes.FakeProfile()
	taxonomy := cohesioned.NewTaxonomy("test", profile)
	taxonomy.ID = 1
	taxonomy.Updated = time.Now()
	taxonomy.UpdatedBy = profile.ID
	err := repo.Update(taxonomy)

	if err != nil {
		t.Errorf("Failed to update taxonomy: %v", err)
	}
}
