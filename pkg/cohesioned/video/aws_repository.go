package video

import (
	"database/sql"
	"fmt"
	"strings"

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
		v.id,
		v.title,
		v.taxonomy_id,
		v.file_name,
		v.file_type,
		v.file_size,
		v.bucket,
		v.object_key,
		v.key_terms,
		v.state_standards,
		v.common_core_standards,
		v.created,
		v.created_by,
		v.updated,
		v.updated_by,
		u.full_name,
		t.name
	from
		video v, user u, taxonomy t
	where
		v.id = ?
	and
		v.taxonomy_id = t.id
	and
		v.created_by = u.id`

	row := repo.QueryRow(selectQuery, id)
	video, err := repo.mapRowToObject(row)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, fmt.Errorf("No video with ID %d", id)
		default:
			return nil, fmt.Errorf("Unexpected error querying for video by id %d: %v", id, err)
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
		v.id,
		v.title,
		v.taxonomy_id,
		v.file_name,
		v.file_type,
		v.file_size,
		v.bucket,
		v.object_key,
		v.key_terms,
		v.state_standards,
		v.common_core_standards,
		v.created,
		v.created_by,
		v.updated,
		v.updated_by,
		u.full_name,
		t.name
		from
			video v, user u, taxonomy t
		where
			v.taxonomy_id = t.id
		and
			v.created_by = u.id`

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
		return list, fmt.Errorf("rows had an error: %v", err)
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
		v.CreatedByID,
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
		v.UpdatedByID,
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

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.Video, error) {
	video := &cohesioned.Video{}
	var updated db.NullTime
	var fileSize, updatedBy sql.NullInt64
	var createdByFullName, taxonomyName, fileType, keyTerms, stateStandards, commonCoreStandards sql.NullString

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
		&video.CreatedByID,
		&updated,
		&updatedBy,
		&createdByFullName,
		&taxonomyName,
	)

	if err != nil {
		return video, fmt.Errorf("failed to map row to video: %v", err)
	}

	video.FileType = fileType.String
	video.FileSize = fileSize.Int64
	video.CreatedBy = &cohesioned.Profile{ID: video.CreatedByID, FullName: createdByFullName.String}
	video.Taxonomy = &cohesioned.Taxonomy{ID: video.TaxonomyID, Name: taxonomyName.String}
	video.Updated = updated.Time
	video.UpdatedByID = updatedBy.Int64

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
