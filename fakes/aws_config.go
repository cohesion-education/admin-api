package fakes

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws/session"
)

type FakeAwsConfig struct {
	err         error
	signedURL   string
	videoBucket string
}

func (cfg *FakeAwsConfig) GetVideoBucketReturns(bucket string) {
	cfg.videoBucket = bucket
}
func (cfg *FakeAwsConfig) GetSignedURLReturns(signedURL string, err error) {
	cfg.signedURL = signedURL
	cfg.err = err
}

func (cfg *FakeAwsConfig) NewSession() (*session.Session, error) {
	return nil, cfg.err
}
func (cfg *FakeAwsConfig) DialRDS() (*sql.DB, error) {
	return nil, cfg.err
}
func (cfg *FakeAwsConfig) GetVideoBucket() string {
	return cfg.videoBucket
}
func (cfg *FakeAwsConfig) GetSignedURL(bucket, objectName string) (string, error) {
	return cfg.signedURL, nil
}
