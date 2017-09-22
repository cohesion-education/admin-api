package video_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/cohesion-education/api/testutils"
	"github.com/joho/godotenv"
)

var (
	repo video.Repo
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		panic("Unexpected error initializing AwsConfig: " + err.Error())
	}

	db, _ := awsConfig.DialRDS()
	repo, err = video.NewAwsRepo(awsConfig, "test-bucket")
	if err != nil {
		panic("Failed to connect to db " + err.Error())
	}

	if err := testutils.SetupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	testResult := m.Run()

	if err := testutils.CleanupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(testResult)
}

func TestGet(t *testing.T) {
	video, err := repo.Get(1)
	if err != nil {
		t.Errorf("Failed to get video by ID: %v", err)
	}

	if video == nil {
		t.Error("Video by ID was null")
	}
}

func TestList(t *testing.T) {
	list, err := repo.List()
	if err != nil {
		t.Errorf("Failed to List videos: %v", err)
	}

	fmt.Printf("video list: %v\n", list)
}
