package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
)

func TestNewAwsConfig(t *testing.T) {
	os.Setenv("VCAP_APPLICATION", vcapApplicationPayload)
	os.Setenv("VCAP_SERVICES", vcapServicesPayload)

	config, err := config.NewAwsConfig()
	if err != nil {
		t.Errorf("Unexpected error initializing AwsConfig: %v", err)
	}

	fmt.Printf("aws config: %v\n", config)

	os.Clearenv()
}
