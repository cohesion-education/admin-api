package config_test

import (
	"os"
	"testing"

	"github.com/cohesion-education/admin-api/pkg/config"
)

func TestNewAuthConfigWithoutEnvVarsOrVcapServicesFails(t *testing.T) {
	if _, err := config.NewAuthConfig(); err == nil {
		t.Error("neither VCAP_SERVICES nor env vars were set - an error should've been returned from newAuthConfig() but was not")
	}
}

func TestNewAuthConfigViaVcapServices(t *testing.T) {
	os.Setenv("VCAP_APPLICATION", VcapApplicationPayload)
	os.Setenv("VCAP_SERVICES", VcapServicesPayload)

	if _, err := config.NewAuthConfig(); err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	os.Clearenv()
}

func TestNewAuthConfigViaPartialVcapServices(t *testing.T) {
	os.Setenv("VCAP_APPLICATION", VcapApplicationPayload)
	os.Setenv("VCAP_SERVICES", VcapServicesPartialPayload)

	_, err := config.NewAuthConfig()

	if err == nil {
		t.Error("expected an error")
	}

	expectedError := "Failed to load auth0 service from either VCAP_SERVICES or from environment vars - missing [ClientID]"

	if err.Error() != expectedError {
		t.Errorf("expected error message %s but got %s", expectedError, err.Error())
	}

	os.Clearenv()
}

func TestNewAuthConfigViaEnvVars(t *testing.T) {
	os.Setenv("AUTH0_CLIENT_ID", "test-client")
	os.Setenv("AUTH0_CLIENT_SECRET", "test-secret")
	os.Setenv("AUTH0_DOMAIN", "test-domain")
	os.Setenv("AUTH0_CALLBACK_URL", "test-callback-url")
	os.Setenv("SESSION_AUTH_KEY", "test-key")

	if _, err := config.NewAuthConfig(); err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	os.Clearenv()
}

func TestNewAuthConfigViaPartialEnvVars(t *testing.T) {
	os.Setenv("AUTH0_CLIENT_SECRET", "test-secret")
	os.Setenv("AUTH0_DOMAIN", "test-domain")
	os.Setenv("AUTH0_CALLBACK_URL", "test-callback-url")
	os.Setenv("SESSION_AUTH_KEY", "test-key")

	_, err := config.NewAuthConfig()

	if err == nil {
		t.Error("expected an error")
	}

	expectedError := "Failed to load auth0 service from either VCAP_SERVICES or from environment vars - missing [ClientID]"

	if err.Error() != expectedError {
		t.Errorf("expected error message %s but got %s", expectedError, err.Error())
	}

	os.Clearenv()
}
