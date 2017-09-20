package video

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	_ "github.com/go-sql-driver/mysql"
)

type awsRepo struct {
	s3BucketName string
	awsConfig    config.AwsConfig
}

func NewAwsRepo(awsConfig config.AwsConfig, s3BucketName string) Repo {
	return &awsRepo{
		s3BucketName: s3BucketName,
		awsConfig:    awsConfig,
	}
}

func (r *awsRepo) Get(id int64) (*cohesioned.Video, error) {
	return nil, nil
}

func (r *awsRepo) Delete(id int64) error {
	return nil
}

func (r *awsRepo) List() ([]*cohesioned.Video, error) {
	var list []*cohesioned.Video

	db, err := r.awsConfig.DialRDS()
	if err != nil {
		return list, fmt.Errorf("Failed to connect to RDS: %v", err)
	}

	rows, err := db.Query("select * from video")
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var video cohesioned.Video
		if err := rows.Scan(&video); err != nil {
			fmt.Errorf("Failed to scan video row: %v", err)
		}
		fmt.Printf("video: %v\n", video)
		list = append(list, &video)
	}

	if err := rows.Err(); err != nil {
		fmt.Errorf("rows had an error: %v", err)
	}

	return list, nil
}

func (r *awsRepo) Add(v *cohesioned.Video) (*cohesioned.Video, error) {
	return nil, nil
}

func (r *awsRepo) Update(v *cohesioned.Video) (*cohesioned.Video, error) {
	v.Updated = time.Now()

	return nil, nil
}

func (r *awsRepo) SetFile(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	return nil, nil
}

func (r *awsRepo) writeFileToStorage(fileReader io.Reader, objectName string) error {
	sess, err := r.awsConfig.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session %v", err)
	}

	svc := s3.New(sess)

	// Upload input parameters
	params := &s3manager.UploadInput{
		Bucket: aws.String(r.s3BucketName),
		Key:    aws.String(objectName),
		Body:   fileReader,
	}

	uploader := s3manager.NewUploaderWithClient(svc)

	// Perform an upload.
	result, err := uploader.Upload(params, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB part size
		u.LeavePartsOnError = false
	})

	if err != nil {
		return fmt.Errorf("Failed to upload file to s3: %v", err)
	}

	fmt.Printf("uploader result: %v\n", result)
	return nil
}
