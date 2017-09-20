package config_test

import (
	"fmt"
	"testing"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/joho/godotenv"
)

func TestNewAwsConfig(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Errorf("Failed to load .env file: %v", err)
	}

	config, err := config.NewAwsConfig()
	if err != nil {
		t.Errorf("Unexpected error initializing AwsConfig: %v", err)
	}

	fmt.Printf("aws config: %v\n", config)
}
