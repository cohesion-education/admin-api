package video

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/cohesion-education/api/pkg/cohesioned"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

type gcpRepo struct {
	storageBucketName string
	datastoreClient   *datastore.Client
	storageClient     *storage.Client
	ctx               context.Context
}

//NewGCPRepo implementation of video.Repo
func NewGCPRepo(datastoreClient *datastore.Client, storageClient *storage.Client, storageBucketName string) Repo {
	return &gcpRepo{
		datastoreClient:   datastoreClient,
		storageClient:     storageClient,
		storageBucketName: storageBucketName,
		ctx:               context.TODO(),
	}
}

func (r *gcpRepo) Get(id int64) (*cohesioned.Video, error) {
	video := &cohesioned.Video{}

	key := datastore.IDKey("Video", id, nil)
	err := r.datastoreClient.Get(r.ctx, key, video)

	if err == datastore.ErrInvalidEntityType {
		return nil, fmt.Errorf("%d returns an invalid entity type %v", id, err)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to get video by id %d %v", id, err)
	}

	return video, nil
}

func (r *gcpRepo) Delete(id int64) error {
	video := &cohesioned.Video{}

	key := datastore.IDKey("Video", id, nil)
	err := r.datastoreClient.Get(r.ctx, key, video)

	if err == datastore.ErrInvalidEntityType {
		return fmt.Errorf("%d returns an invalid entity type %v", id, err)
	}

	if err != nil {
		return fmt.Errorf("Failed to get video by id %d %v", id, err)
	}

	objectHandle := r.storageClient.Bucket(video.StorageBucket).Object(video.StorageObjectName)
	if err := objectHandle.Delete(r.ctx); err != nil {
		return fmt.Errorf("Failed to delete video storage object %s/%s for video %d: %v", video.StorageBucket, video.StorageObjectName, video.ID(), err)
	}

	if err := r.datastoreClient.Delete(r.ctx, key); err != nil {
		return fmt.Errorf("Failed to delete video with id %d from data store: %v", video.ID(), err)
	}

	return nil
}

func (r *gcpRepo) List() ([]*cohesioned.Video, error) {
	var list []*cohesioned.Video

	q := datastore.NewQuery("Video")
	keys, err := r.datastoreClient.GetAll(r.ctx, q, &list)
	if err != nil {
		return nil, fmt.Errorf("Failed to get video list from Cloud Datastore %v", err)
	}

	for i, key := range keys {
		list[i].SetID(key.ID)
	}

	return list, nil
}

func (r *gcpRepo) Add(v *cohesioned.Video) (*cohesioned.Video, error) {
	key, err := r.datastoreClient.Put(r.ctx, datastore.IncompleteKey("Video", nil), v)
	v.Key = key
	v.Created = time.Now()
	if err != nil {
		return nil, fmt.Errorf("Failed to save video %v", err)
	}

	return v, nil
}

func (r *gcpRepo) Update(v *cohesioned.Video) (*cohesioned.Video, error) {
	v.Updated = time.Now()

	_, err := r.datastoreClient.Put(r.ctx, v.Key, v)
	if err != nil {
		return v, fmt.Errorf("Failed to save video %v", err)
	}

	return v, nil
}

func (r *gcpRepo) SetFile(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	objectName := fmt.Sprintf("%d-%s", video.ID(), video.FileName)
	if err := r.writeFileToStorage(fileReader, objectName); err != nil {
		return video, err
	}

	video.Updated = time.Now()
	video.StorageBucket = r.storageBucketName
	video.StorageObjectName = objectName
	if _, err := r.datastoreClient.Put(r.ctx, video.Key, video); err != nil {
		return video, fmt.Errorf("Failed to save video %v", err)
	}

	return video, nil
}

func (r *gcpRepo) writeFileToStorage(fileReader io.Reader, objectName string) error {
	if fileReader == nil {
		return nil
	}

	objectHandle := r.storageClient.Bucket(r.storageBucketName).Object(objectName)
	writer := objectHandle.NewWriter(r.ctx)

	if _, err := io.Copy(writer, fileReader); err != nil {
		return fmt.Errorf("Failed to write video file %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("Failed to write video file %v", err)
	}

	return nil
}
