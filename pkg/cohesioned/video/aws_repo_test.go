package video_test

import (
	"fmt"
	"testing"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/joho/godotenv"
)

func TestList(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Errorf("Failed to load .env file: %v", err)
	}

	config, err := config.NewAwsConfig()
	if err != nil {
		t.Errorf("Unexpected error initializing AwsConfig: %v", err)
	}

	repo := video.NewAwsRepo(config, "test-bucket")
	list, err := repo.List()
	if err != nil {
		t.Errorf("Failed to List videos: %v", err)
	}

	fmt.Printf("video list: %v\n", list)
}
