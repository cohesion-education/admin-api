package student

import (
	"database/sql"
	"fmt"

	"github.com/cohesion-education/api/pkg/cohesioned"
	"github.com/cohesion-education/api/pkg/cohesioned/db"
)

type awsRepo struct {
	*sql.DB
}

func NewAwsRepo(db *sql.DB) Repo {
	return &awsRepo{
		DB: db,
	}
}

func (repo *awsRepo) Save(s *cohesioned.Student) (int64, error) {
	sql := `insert into student
	(
		name,
		grade,
		school,
		user_id,
		created,
		created_by
	) values (?, ?, ?, ?, ?, ?)`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		s.Name,
		s.Grade,
		s.School,
		s.ParentID,
		s.Created,
		s.CreatedBy,
	)

	if err != nil {
		return 0, fmt.Errorf("Failed to insert student: %v", err)
	}

	studentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id from result: %v", err)
	}

	return studentID, nil
}

func (repo *awsRepo) Update(s *cohesioned.Student) error {
	sql := `update student set
		name = ?,
		grade = ?,
		school = ?,
		user_id = ?,
		updated = ?,
		updated_by = ?
	where id = ?`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		s.Name,
		s.Grade,
		s.School,
		s.ParentID,
		s.Updated,
		s.UpdatedBy,
		s.ID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update student: %v", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil || rowsEffected == 0 {
		return fmt.Errorf("Failed to update student: %v", err)
	}

	return nil
}

func (repo *awsRepo) List(parentID int64) ([]*cohesioned.Student, error) {
	var list []*cohesioned.Student

	query := `
	select
	  id,
		name,
		grade,
		school,
		user_id,
		created,
		created_by,
		updated,
		updated_by
	from
		student
	where
		user_id = ?`

	rows, err := repo.Query(query, parentID)
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		student, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the result set from the db: %v", err)
		}

		list = append(list, student)
	}

	if err := rows.Err(); err != nil {
		return list, fmt.Errorf("db rows returned unexpected error: %v", err)
	}

	return list, nil
}

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.Student, error) {
	student := &cohesioned.Student{}
	var updated db.NullTime
	var updatedBy sql.NullInt64

	err := rs.Scan(
		&student.ID,
		&student.Name,
		&student.Grade,
		&student.School,
		&student.ParentID,
		&student.Created,
		&student.CreatedBy,
		&updated,
		&updatedBy,
	)

	if err != nil {
		return student, fmt.Errorf("faled to map row to student: %v", err)
	}

	student.Updated = updated.Time
	student.UpdatedBy = updatedBy.Int64

	return student, nil
}
