package video

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/cohesion-education/admin-api/pkg/cohesioned"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

type Repo interface {
	List() ([]*cohesioned.Video, error)
	Get(id int64) (*cohesioned.Video, error)
	Add(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error)
	Update(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error)
}

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

func (r *gcpRepo) Add(fileReader io.Reader, v *cohesioned.Video) (*cohesioned.Video, error) {
	v.Created = time.Now()
	v.StorageBucket = r.storageBucketName

	v.Key = datastore.IncompleteKey("Video", nil)
	fmt.Printf("incomplete key: %v", v.Key)
	key, err := r.datastoreClient.Put(r.ctx, v.Key, v)
	fmt.Printf("key returned from Put: %v", key)
	if err != nil {
		return v, fmt.Errorf("Failed to save video %v", err)
	}

	v, err = r.writeFileToStorage(fileReader, v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (r *gcpRepo) Update(fileReader io.Reader, v *cohesioned.Video) (*cohesioned.Video, error) {
	v.Updated = time.Now()

	_, err := r.datastoreClient.Put(r.ctx, v.Key, v)
	if err != nil {
		return v, fmt.Errorf("Failed to save video %v", err)
	}

	v, err = r.writeFileToStorage(fileReader, v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (r *gcpRepo) writeFileToStorage(fileReader io.Reader, v *cohesioned.Video) (*cohesioned.Video, error) {
	if fileReader == nil {
		return v, nil
	}

	v.StorageObjectName = fmt.Sprintf("%d-%s", v.ID(), v.FileName)
	objectHandle := r.storageClient.Bucket(v.StorageBucket).Object(v.StorageObjectName)
	writer := objectHandle.NewWriter(r.ctx)

	if _, err := io.Copy(writer, fileReader); err != nil {
		return v, fmt.Errorf("Failed to write video file %v", err)
	}

	if err := writer.Close(); err != nil {
		return v, fmt.Errorf("Failed to write video file %v", err)
	}

	if _, err := r.datastoreClient.Put(r.ctx, v.Key, v); err != nil {
		return v, fmt.Errorf("Failed to update video %v", err)
	}

	return v, nil
}
