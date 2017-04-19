package gcp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	projectID       string
	keyfileLocation string
	privateKey      []byte
	googleAccessID  string
}

func NewConfig(gcpKeyfileLocation string) (*Config, error) {
	if len(gcpKeyfileLocation) == 0 {
		return nil, fmt.Errorf("gcp keyfile location is required")
	}

	keyfileBytes, err := ioutil.ReadFile(gcpKeyfileLocation)
	if err != nil {
		return nil, fmt.Errorf("Failed to load gcp keyfile %s %v", gcpKeyfileLocation, err)
	}

	keyfile := make(map[string]string)
	if err := json.Unmarshal(keyfileBytes, &keyfile); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall keyfile json %s %v", gcpKeyfileLocation, err)
	}

	cfg := &Config{
		keyfileLocation: gcpKeyfileLocation,
		privateKey:      []byte(keyfile["private_key"]),
		googleAccessID:  keyfile["client_email"],
		projectID:       keyfile["project_id"],
	}

	return cfg, nil
}
