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
	"github.com/cohesion-education/api/pkg/cohesioned/taxonomy"
)

type AdminService interface {
	List() ([]*cohesioned.Video, error)
	FindByTaxonomyID(taxonomyID int64) ([]*cohesioned.Video, error)
	FindByGrade(gradeName string) (map[string][]*cohesioned.Video, error)
	Get(id int64) (*cohesioned.Video, error)
	GetWithSignedURL(id int64) (*cohesioned.Video, error)
	Delete(id int64) error
	Save(ctx context.Context, video *cohesioned.Video) error
	Update(ctx context.Context, video *cohesioned.Video) error
	SetFile(ctx context.Context, fileReader io.Reader, video *cohesioned.Video) error
}

type adminService struct {
	videoRepo    Repo
	taxonomyRepo taxonomy.Repo
	cfg          config.AwsConfig
}

func NewService(videoRepo Repo, taxonomyRepo taxonomy.Repo, cfg config.AwsConfig) AdminService {
	return &adminService{
		videoRepo:    videoRepo,
		taxonomyRepo: taxonomyRepo,
		cfg:          cfg,
	}
}

func (s *adminService) List() ([]*cohesioned.Video, error) {
	return s.videoRepo.List()
}

func (s *adminService) FindByTaxonomyID(taxonomyID int64) ([]*cohesioned.Video, error) {
	return s.videoRepo.FindByTaxonomyID(taxonomyID)
}

func (s *adminService) FindByGrade(gradeName string) (map[string][]*cohesioned.Video, error) {
	videosByFlattenedTaxonomy := make(map[string][]*cohesioned.Video)

	grade, err := s.taxonomyRepo.FindGradeByName(gradeName)
	if err != nil {
		return videosByFlattenedTaxonomy, fmt.Errorf("Failed to find grade %s: %v", gradeName, err)
	}

	children, err := s.taxonomyRepo.ListChildren(grade.ID)
	if err != nil {
		return videosByFlattenedTaxonomy, fmt.Errorf("Failed to get taxonomy children for grade with ID %d: %v", grade.ID, err)
	}

	for _, child := range children {
		flattened, err := s.taxonomyRepo.Flatten(child)
		if err != nil {
			fmt.Printf("Failed to flatten child %v: %v\n", child, err)
			continue
		}

		for _, f := range flattened {
			videos, err := s.videoRepo.FindByTaxonomyID(f.ID)
			if err != nil {
				fmt.Printf("Failed to find videos by taxonomy ID %d: %v\n", f.ID, err)
				continue
			}

			videosByFlattenedTaxonomy[f.Name] = videos
		}
	}

	return videosByFlattenedTaxonomy, nil
}

func (s *adminService) Get(id int64) (*cohesioned.Video, error) {
	video, err := s.videoRepo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get video by ID: %v", err)
	}

	video.Taxonomy, err = s.taxonomyRepo.ReverseFlatten(video.Taxonomy)
	if err != nil {
		return nil, fmt.Errorf("Failed to get full taxonomy for video: %v", err)
	}

	return video, nil
}

func (s *adminService) GetWithSignedURL(id int64) (*cohesioned.Video, error) {
	video, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if len(video.StorageBucket) != 0 && len(video.StorageObjectName) != 0 {
		signedURL, err := s.cfg.GetSignedURL(video.StorageBucket, video.StorageObjectName)
		video.SignedURL = signedURL
		if err != nil {
			return nil, fmt.Errorf("Failed to generate signed url %v", err)
		}
	}

	return video, nil
}

func (s *adminService) Delete(id int64) error {
	video, err := s.Get(id)
	if err != nil {
		return fmt.Errorf("Failed to find video with id %d: %v", id, err)
	}

	if len(video.StorageBucket) != 0 && len(video.StorageObjectName) != 0 {
		if err := s.deleteFile(video.StorageBucket, video.StorageObjectName); err != nil {
			return fmt.Errorf("Failed to delete the video file: %v", err)
		}
	}

	return s.videoRepo.Delete(id)
}

//Save saves the given video, looking for the current user in the given context argument. Sets the resulting ID from the save operation on the video instance
func (s *adminService) Save(ctx context.Context, video *cohesioned.Video) error {
	currentUser, _ := cohesioned.FromContext(ctx)
	video.CreatedByID = currentUser.ID

	id, err := s.videoRepo.Save(video)
	if err != nil {
		return err
	}

	video.ID = id
	return nil
}

func (s *adminService) Update(ctx context.Context, video *cohesioned.Video) error {
	currentUser, _ := cohesioned.FromContext(ctx)
	video.Updated = time.Now()
	video.UpdatedByID = currentUser.ID

	return s.videoRepo.Update(video)
}

func (s *adminService) SetFile(ctx context.Context, fileReader io.Reader, video *cohesioned.Video) error {
	//TODO - wrap in transaction
	//TODO - delete existing file
	video.StorageBucket = s.cfg.GetVideoBucket()
	video.StorageObjectName = fmt.Sprintf("%d-%s", video.ID, video.FileName)

	if err := s.writeFileToStorage(fileReader, video.StorageBucket, video.StorageObjectName); err != nil {
		return fmt.Errorf("Failed to write file to storage: %v", err)
	}

	// if err := s.submitTranscodingJobs(video); err != nil {
	// 	return fmt.Errorf("Failed to submit transcoding jobs: %v", err)
	// }

	if err := s.Update(ctx, video); err != nil {
		return fmt.Errorf("Failed to update video record: %v", err)
	}

	return nil
}

func (s *adminService) deleteFile(bucketName, objectName string) error {
	sess, err := s.cfg.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session %v", err)
	}

	svc := s3.New(sess)

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	if _, err := svc.DeleteObject(deleteInput); err != nil {
		return fmt.Errorf("Failed to delete %s/%s: %v", bucketName, objectName, err)
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

	//TODO - grab this config from s.cfg.TranscodingParams
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
