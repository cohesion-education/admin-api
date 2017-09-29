package video

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/db"
	_ "github.com/go-sql-driver/mysql"
)

type awsRepo struct {
	*sql.DB
	awsConfig config.AwsConfig
}

func NewAwsRepo(db *sql.DB, awsConfig config.AwsConfig) Repo {
	return &awsRepo{
		DB:        db,
		awsConfig: awsConfig,
	}
}

func (repo *awsRepo) Get(id int64) (*cohesioned.Video, error) {
	selectQuery := `select
		id,
		title,
		taxonomy_id,
		file_name,
		file_type,
		file_size,
		bucket,
		object_key,
		key_terms,
		state_standards,
		common_core_standards,
		created,
		created_by,
		updated,
		updated_by
	from
		video
	where
		id = ?`

	row := repo.QueryRow(selectQuery, id)
	video, err := repo.mapRowToObject(row)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, fmt.Errorf("No video with ID %d", id)
		default:
			return nil, fmt.Errorf("Unexpected error querying for user by id %d: %v", id, err)
		}
	}

	return video, nil
}

func (repo *awsRepo) Delete(id int64) error {
	deleteSql := `delete from video where id = ?`
	result, err := repo.Exec(deleteSql, id)

	if err != nil {
		return fmt.Errorf("Failed to delete video with id %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get number of rows affected from result: %v", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("Failed to delete video with id %d (rows affected != 1)", id)
	}

	return nil
}

func (repo *awsRepo) List() ([]*cohesioned.Video, error) {
	var list []*cohesioned.Video

	selectQuery := `select
		id,
		title,
		taxonomy_id,
		file_name,
		file_type,
		file_size,
		bucket,
		object_key,
		key_terms,
		state_standards,
		common_core_standards,
		created,
		created_by,
		updated,
		updated_by
	from
		video`

	rows, err := repo.Query(selectQuery)
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		video, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the result set from the db: %v", err)
		}

		list = append(list, video)
	}

	if err := rows.Err(); err != nil {
		fmt.Errorf("rows had an error: %v", err)
	}

	return list, nil
}

func (repo *awsRepo) Save(v *cohesioned.Video) (int64, error) {
	insertSql := `insert into video
	(
		title,
		taxonomy_id,
		file_name,
		file_type,
		file_size,
		bucket,
		object_key,
		key_terms,
		state_standards,
		common_core_standards,
		created,
		created_by
	)
	values
	(
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`

	stmt, err := repo.Prepare(insertSql)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare statement %s: %v", insertSql, err)
	}

	result, err := stmt.Exec(
		v.Title,
		v.TaxonomyID,
		v.FileName,
		v.FileType,
		v.FileSize,
		v.StorageBucket,
		v.StorageObjectName,
		strings.Join(v.KeyTerms, ","),
		strings.Join(v.StateStandards, ","),
		strings.Join(v.CommonCoreStandards, ","),
		v.Created,
		v.CreatedBy,
	)

	if err != nil {
		return 0, fmt.Errorf("Failed to insert video: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id from result: %v", err)
	}

	return id, nil
}

func (repo *awsRepo) Update(v *cohesioned.Video) error {
	updateSql := `update video set
		title = ?,
		taxonomy_id = ?,
		file_name = ?,
		file_type = ?,
		file_size = ?,
		bucket = ?,
		object_key = ?,
		key_terms = ?,
		state_standards = ?,
		common_core_standards = ?,
		updated = ?,
		updated_by = ?
	where
		id = ?`

	stmt, err := repo.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("Failed to prepare update statement %s: %v", updateSql, err)
	}

	result, err := stmt.Exec(
		v.Title,
		v.TaxonomyID,
		v.FileName,
		v.FileType,
		v.FileSize,
		v.StorageBucket,
		v.StorageObjectName,
		strings.Join(v.KeyTerms, ","),
		strings.Join(v.StateStandards, ","),
		strings.Join(v.CommonCoreStandards, ","),
		v.Updated,
		v.UpdatedBy,
		v.ID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update video: %v", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil || rowsEffected == 0 {
		return fmt.Errorf("Failed to update video: %v", err)
	}

	return nil
}

func (r *awsRepo) SetFile(fileReader io.Reader, video *cohesioned.Video) (*cohesioned.Video, error) {
	if err := r.writeFileToStorage(fileReader, video.StorageObjectName); err != nil {
		return video, fmt.Errorf("Failed to write file to storage: %v", err)
	}

	if err := r.Update(video); err != nil {
		return video, fmt.Errorf("Failed to update video record: %v", err)
	}

	return video, nil
}

func (r *awsRepo) writeFileToStorage(fileReader io.Reader, objectName string) error {
	sess, err := r.awsConfig.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating session %v", err)
	}

	svc := s3.New(sess)

	// Upload input parameters
	params := &s3manager.UploadInput{
		Bucket: aws.String(r.awsConfig.GetVideoBucket()),
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

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.Video, error) {
	video := &cohesioned.Video{}
	var updated db.NullTime
	var fileSize, updatedBy sql.NullInt64
	var fileType, keyTerms, stateStandards, commonCoreStandards sql.NullString

	err := rs.Scan(
		&video.ID,
		&video.Title,
		&video.TaxonomyID,
		&video.FileName,
		&fileType,
		&fileSize,
		&video.StorageBucket,
		&video.StorageObjectName,
		&keyTerms,
		&stateStandards,
		&commonCoreStandards,
		&video.Created,
		&video.CreatedBy,
		&updated,
		&updatedBy,
	)

	if err != nil {
		return video, fmt.Errorf("faled to map row to video: %v", err)
	}

	video.FileType = fileType.String
	video.FileSize = fileSize.Int64
	video.Updated = updated.Time
	video.UpdatedBy = updatedBy.Int64

	if len(keyTerms.String) > 0 {
		video.KeyTerms = strings.Split(keyTerms.String, ",")
	}

	if len(stateStandards.String) > 0 {
		video.StateStandards = strings.Split(stateStandards.String, ",")
	}

	if len(commonCoreStandards.String) > 0 {
		video.CommonCoreStandards = strings.Split(commonCoreStandards.String, ",")
	}

	return video, nil
}
