package video

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
)

type AdminService interface {
	List() ([]*cohesioned.Video, error)
	Get(id int64) (*cohesioned.Video, error)
	GetWithSignedURL(id int64) (*cohesioned.Video, error)
	Delete(id int64) error
	Save(ctx context.Context, video *cohesioned.Video) error
	Update(ctx context.Context, video *cohesioned.Video) error
	SetFile(ctx context.Context, fileReader io.Reader, video *cohesioned.Video) error
}

type adminService struct {
	repo Repo
	cfg  config.AwsConfig
}

func NewService(repo Repo, cfg config.AwsConfig) AdminService {
	return &adminService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *adminService) List() ([]*cohesioned.Video, error) {
	return s.repo.List()
}

func (s *adminService) Get(id int64) (*cohesioned.Video, error) {
	video, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get video by ID: %v", err)
	}

	return video, nil
}

func (s *adminService) GetWithSignedURL(id int64) (*cohesioned.Video, error) {
	video, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	signedURL, err := s.cfg.GetSignedURL(video.StorageBucket, video.StorageObjectName)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate signed url %v", err)
	}

	video.SignedURL = signedURL
	return video, nil
}

func (s *adminService) Delete(id int64) error {
	//TODO - also delete file!
	return s.repo.Delete(id)
}

//Save saves the given video, looking for the current user in the given context argument. Sets the resulting ID from the save operation on the video instance
func (s *adminService) Save(ctx context.Context, video *cohesioned.Video) error {
	currentUser, _ := cohesioned.FromContext(ctx)
	video.CreatedBy = currentUser.ID

	id, err := s.repo.Save(video)
	if err != nil {
		return err
	}

	video.ID = id
	return nil
}

func (s *adminService) Update(ctx context.Context, video *cohesioned.Video) error {
	currentUser, _ := cohesioned.FromContext(ctx)
	video.Updated = time.Now()
	video.UpdatedBy = currentUser.ID

	return s.repo.Update(video)
}

func (s *adminService) SetFile(ctx context.Context, fileReader io.Reader, video *cohesioned.Video) error {
	//TODO - wrap in transaction
	//TODO - delete existing file
	video.StorageBucket = s.cfg.GetVideoBucket()
	video.StorageObjectName = fmt.Sprintf("%d-%s", video.ID, video.FileName)

	if err := s.writeFileToStorage(fileReader, video.StorageBucket, video.StorageObjectName); err != nil {
		return fmt.Errorf("Failed to write file to storage: %v", err)
	}

	if err := s.submitTranscodingJobs(video); err != nil {
		return fmt.Errorf("Failed to submit transcoding jobs: %v", err)
	}

	if err := s.Update(ctx, video); err != nil {
		return fmt.Errorf("Failed to update video record: %v", err)
	}

	return nil
}

func (s *adminService) writeFileToStorage(fileReader io.Reader, bucketName, objectName string) error {
	sess, err := s.cfg.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session %v", err)
	}

	svc := s3.New(sess)

	// Upload input parameters
	params := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   fileReader,
	}

	uploader := s3manager.NewUploaderWithClient(svc)

	// Perform an upload.
	if _, err := uploader.Upload(params, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB part size
		u.LeavePartsOnError = false
	}); err != nil {
		return fmt.Errorf("Failed to upload file to s3: %v", err)
	}

	return nil
}

func (s *adminService) submitTranscodingJobs(v *cohesioned.Video) error {
	sess, err := s.cfg.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session %v", err)
	}

	svc := elastictranscoder.New(sess)

	pipelineID := "1507471006678-grs7px"
	outputFolder := "outputs/"
	outputKeyPrefix := "480p-16x9"
	outputPresetID := "1351620000001-000020"

	input := &elastictranscoder.CreateJobInput{}
	input.SetPipelineId(pipelineID)
	input.SetOutputKeyPrefix(outputFolder)
	input.SetInput(&elastictranscoder.JobInput{
		Key: aws.String(v.StorageObjectName),
	})
	input.SetOutput(&elastictranscoder.CreateJobOutput{
		Key:              aws.String(fmt.Sprintf("%s-%s", outputKeyPrefix, v.StorageObjectName)),
		PresetId:         aws.String(outputPresetID),
		ThumbnailPattern: aws.String(fmt.Sprintf("%s-{resolution}-{count}", v.StorageObjectName)),
	})

	resp, err := svc.CreateJob(input)
	if err != nil {
		return fmt.Errorf("Failed to create job: %v", err)
	}

	fmt.Printf("create job response: %v\n", resp)
	return nil
}
