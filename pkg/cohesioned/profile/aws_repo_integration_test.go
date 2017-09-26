//+build integration

package profile_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/profile"
	"github.com/cohesion-education/api/testutils"
)

var (
	repo profile.Repo
)

func TestMain(m *testing.M) {
	// if err := godotenv.Load("../../../.env"); err != nil {
	// 	panic("Failed to load .env file: " + err.Error())
	// }

	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		panic("Unexpected error initializing AwsConfig: " + err.Error())
	}

	db, err := awsConfig.DialRDS()
	if err != nil {
		panic("Failed to connect to db " + err.Error())
	}

	repo = profile.NewAwsRepo(db)

	if err := testutils.SetupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	testResult := m.Run()

	if err := testutils.CleanupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(testResult)
}

func TestRepoSave(t *testing.T) {
	profile := fakes.FakeProfile()
	profile.Preferences.BetaProgram = true
	profile.Preferences.Newsletter = true

	id, err := repo.Save(profile)

	if err != nil {
		t.Errorf("Failed to save profile: %v", err)
	}

	if id == 0 {
		t.Errorf("Profile ID was zero - expected db to generate a profile id")
	}
}

func TestRepoUpdate(t *testing.T) {
	profile := fakes.FakeProfile()
	profile.ID = int64(2)
	profile.FullName = "Updated User"
	profile.FirstName = "Updated"
	profile.LastName = "User"
	profile.Updated = time.Now()

	if err := repo.Update(profile); err != nil {
		t.Errorf("Failed to save profile: %v", err)
	}
}

func TestRepoFindByEmail(t *testing.T) {
	email := "hello@domain.com"
	profile, err := repo.FindByEmail(email)

	if err != nil {
		t.Errorf("Failed to find profile by email %s: %v", email, err)
	}

	if profile == nil {
		t.Fatalf("Profile returned from repo was nil")
	}

	fmt.Printf("Profile: %v\n", profile)

	if len(profile.Email) == 0 {
		t.Errorf("Profile returned from repo had an empty email")
	}

	if len(profile.FullName) == 0 {
		t.Errorf("Profile returned from repo had an empty full name")
	}

	if len(profile.FirstName) == 0 {
		t.Errorf("Profile returned from repo had an empty first name")
	}

	if len(profile.LastName) == 0 {
		t.Errorf("Profile returned from repo had an empty last name")
	}

	if profile.ID == 0 {
		t.Errorf("Profile returned from repo had an empty id")
	}
}
