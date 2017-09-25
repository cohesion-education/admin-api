package taxonomy

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

func (repo *awsRepo) Get(id int64) (*cohesioned.Taxonomy, error) {
	taxonomy := new(cohesioned.Taxonomy)

	query := `select
		id,
		name,
		parent_id,
		created,
		created_by,
		updated,
		updated_by
	from
		taxonomy
	where
		id = ?`

	row := repo.QueryRow(query, id)
	taxonomy, err := repo.mapRowToObject(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%d is not a valid taxonomy id", id)
		}

		return nil, err
	}

	return taxonomy, nil
}

func (repo *awsRepo) List() ([]*cohesioned.Taxonomy, error) {
	var list []*cohesioned.Taxonomy

	query := `select
		id,
		name,
		parent_id,
		created,
		created_by,
		updated,
		updated_by
	from
		taxonomy`

	rows, err := repo.Query(query)
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		taxonomy, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the result set from the db: %v", err)
		}

		list = append(list, taxonomy)
	}

	if err := rows.Err(); err != nil {
		return list, fmt.Errorf("db rows returned unexpected error: %v", err)
	}

	return list, nil
}

func (repo *awsRepo) ListChildren(parentID int64) ([]*cohesioned.Taxonomy, error) {
	var list []*cohesioned.Taxonomy

	query := `select
		id,
		name,
		parent_id,
		created,
		created_by,
		updated,
		updated_by
	from
		taxonomy
	where
		parent_id = ?`

	rows, err := repo.Query(query, parentID)
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		taxonomy, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the result set from the db: %v", err)
		}

		list = append(list, taxonomy)
	}

	if err := rows.Err(); err != nil {
		return list, fmt.Errorf("db rows returned unexpected error: %v", err)
	}

	return list, nil
}

func (repo *awsRepo) Save(t *cohesioned.Taxonomy) (int64, error) {
	insertSql := `insert into taxonomy
	(
		name,
		parent_id,
		created,
		created_by
	) values (?, ?, ?, ?)`

	stmt, err := repo.Prepare(insertSql)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare statement %s: %v", insertSql, err)
	}

	var parentID interface{}
	if t.ParentID != 0 {
		parentID = t.ParentID
	}

	result, err := stmt.Exec(
		t.Name,
		parentID,
		t.Created,
		t.CreatedBy,
	)

	if err != nil {
		return 0, fmt.Errorf("Failed to insert taxonomy: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id from result: %v", err)
	}

	return id, nil
}

func (repo *awsRepo) Update(t *cohesioned.Taxonomy) error {
	updateSql := `update taxonomy set
		name = ?,
		parent_id = ?,
		updated = ?,
		updated_by = ?
	where id = ?`

	stmt, err := repo.Prepare(updateSql)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement %s: %v", updateSql, err)
	}

	var parentID interface{}
	if t.ParentID != 0 {
		parentID = t.ParentID
	}

	result, err := stmt.Exec(
		t.Name,
		parentID,
		t.Updated,
		t.UpdatedBy,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update taxonomy: %v", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil || rowsEffected == 0 {
		return fmt.Errorf("Failed to update taxonomy: %v", err)
	}

	return nil
}

func (repo *awsRepo) Flatten(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error) {
	flattened := []*cohesioned.Taxonomy{}
	if t == nil {
		return flattened, nil
	}

	children, err := repo.ListChildren(t.ID)
	if err != nil {
		return flattened, fmt.Errorf("Failed to get children for %s %v\n", t.Name, err)
	}

	if len(children) == 0 {
		fmt.Printf("Flattened: %s\n", t.Name)
		flattened = append(flattened, t)
		return flattened, nil
	}

	for _, child := range children {
		child.Name = fmt.Sprintf("%s > %s", t.Name, child.Name)
		flattenedChildren, err := repo.Flatten(child)
		if err != nil {
			return flattened, fmt.Errorf("Failed to flatten children of %s %v", child.Name, err)
		}
		flattened = append(flattened, flattenedChildren...)
	}

	return flattened, nil
}

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.Taxonomy, error) {
	taxonomy := &cohesioned.Taxonomy{}
	var parentID sql.NullInt64
	var updated db.NullTime
	var updatedBy sql.NullInt64

	err := rs.Scan(
		&taxonomy.ID,
		&taxonomy.Name,
		&parentID,
		&taxonomy.Created,
		&taxonomy.CreatedBy,
		&updated,
		&updatedBy,
	)

	if err != nil {
		return taxonomy, fmt.Errorf("faled to map row to taxonomy: %v", err)
	}

	taxonomy.ParentID = parentID.Int64
	taxonomy.Updated = updated.Time
	taxonomy.UpdatedBy = updatedBy.Int64

	return taxonomy, nil
}
