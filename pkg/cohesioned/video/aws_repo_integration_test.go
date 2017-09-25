//+build integration

package video_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/cohesion-education/api/testutils"
)

var (
	repo      video.Repo
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

	repo = video.NewAwsRepo(db, awsConfig)

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

	if len(list) == 0 {
		t.Error("repo returned empty list")
	}

	for _, video := range list {
		if len(video.Title) == 0 {
			t.Error("Video Title is empty")
		}

		if video.TaxonomyID == 0 {
			t.Error("Video TaxonomyID is empty")
		}

		if len(video.FileName) == 0 {
			t.Error("Video FileName is empty")
		}

		if video.Created == emptyTime {
			t.Error("Video created is empty")
		}

		if video.CreatedBy == 0 {
			t.Error("Video created by is empty")
		}
	}
}

func TestRepoSave(t *testing.T) {
	video := fakes.FakeVideo()
	id, err := repo.Save(video)

	if err != nil {
		t.Errorf("Failed to save video: %v", err)
	}

	if id == 0 {
		t.Errorf("Video ID was zero - expected db to generate a video id")
	}
}
