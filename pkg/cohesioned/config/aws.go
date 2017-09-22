package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

type AwsConfig interface {
	NewSession() (*session.Session, error)
	DialRDS() (*sql.DB, error)
}

type rds struct {
	username string
	password string
	host     string
	port     int16
	dbname   string
}

func (r *rds) String() string {
	return fmt.Sprintf("host: %s port: %d db: %s", r.host, r.port, r.dbname)
}

func (r *rds) DialRDS() (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&parseTime=true",
		r.username,
		r.password,
		r.host,
		r.port,
		r.dbname,
	)

	return sql.Open("mysql", connStr)
}

type config struct {
	rds
	region          string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

func (c *config) String() string {
	return fmt.Sprintf("Region: %s\tAccess Key ID: %s", c.region, c.accessKeyID)
}

func (c *config) NewSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.region),
		Credentials: credentials.NewStaticCredentials(c.accessKeyID, c.secretAccessKey, c.sessionToken),
	})

	return sess, err
}

func NewAwsConfig() (AwsConfig, error) {
	config := &config{}

	if appEnv, err := cfenv.Current(); err == nil {
		if awsService, err := appEnv.Services.WithName("aws"); err == nil {
			if region, ok := awsService.CredentialString("region"); ok {
				config.region = region
			}
			if accessKeyID, ok := awsService.CredentialString("access_key_id"); ok {
				config.accessKeyID = accessKeyID
			}
			if secretAccessKey, ok := awsService.CredentialString("secret_access_key"); ok {
				config.secretAccessKey = secretAccessKey
			}
			if sessionToken, ok := awsService.CredentialString("session_token"); ok {
				config.sessionToken = sessionToken
			}
			if rdsUsername, ok := awsService.CredentialString("rds_username"); ok {
				config.rds.username = rdsUsername
			}
			if rdsPassword, ok := awsService.CredentialString("rds_password"); ok {
				config.rds.password = rdsPassword
			}
			if rdsHost, ok := awsService.CredentialString("rds_host"); ok {
				config.rds.host = rdsHost
			}
			if rdsPort, ok := awsService.CredentialString("rds_port"); ok {
				if port, err := strconv.Atoi(rdsPort); err == nil {
					config.rds.port = int16(port)
				}
			}
			if rdsDbname, ok := awsService.CredentialString("rds_dbname"); ok {
				config.rds.dbname = rdsDbname
			}
		}
	}

	if len(config.region) == 0 {
		config.region = os.Getenv("AWS_REGION")
	}

	if len(config.accessKeyID) == 0 {
		config.accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if len(config.secretAccessKey) == 0 {
		config.secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	if len(config.sessionToken) == 0 {
		config.sessionToken = os.Getenv("AWS_SESSION_TOKEN")
	}

	if len(config.rds.username) == 0 {
		config.rds.username = os.Getenv("AWS_RDS_USERNAME")
	}

	if len(config.rds.password) == 0 {
		config.rds.password = os.Getenv("AWS_RDS_PASSWORD")
	}

	if len(config.rds.host) == 0 {
		config.rds.host = os.Getenv("AWS_RDS_HOST")
	}

	if config.rds.port == 0 {
		rdsPort := os.Getenv("AWS_RDS_PORT")
		if port, err := strconv.Atoi(rdsPort); err == nil {
			config.rds.port = int16(port)
		}
	}

	if len(config.rds.dbname) == 0 {
		config.rds.dbname = os.Getenv("AWS_RDS_DBNAME")
	}

	var missingConfig []string
	if len(config.region) == 0 {
		missingConfig = append(missingConfig, "Region")
	}

	if len(config.secretAccessKey) == 0 {
		missingConfig = append(missingConfig, "SecretAccessKey")
	}

	if len(config.accessKeyID) == 0 {
		missingConfig = append(missingConfig, "AccessKeyID")
	}

	if len(config.rds.username) == 0 {
		missingConfig = append(missingConfig, "RDS.Username")
	}

	// if len(config.rds.password) == 0 {
	// 	missingConfig = append(missingConfig, "RDS.Password")
	// }

	if len(config.rds.host) == 0 {
		missingConfig = append(missingConfig, "RDS.Host")
	}

	if config.rds.port == 0 {
		missingConfig = append(missingConfig, "RDS.Port")
	}

	if len(config.rds.dbname) == 0 {
		missingConfig = append(missingConfig, "RDS.Dbname")
	}

	if len(missingConfig) > 0 {
		return nil, fmt.Errorf("Failed to load aws service from either VCAP_SERVICES or from environment vars - missing %v", missingConfig)
	}

	return config, nil
}
