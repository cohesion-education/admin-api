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

func (repo *awsRepo) FindGradeByName(name string) (*cohesioned.Taxonomy, error) {
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
		name = ?
	and
		parent_id is null`

	row := repo.QueryRow(query, name)
	taxonomy, err := repo.mapRowToObject(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("can't find grade with name %s", name)
		}

		return nil, err
	}

	return taxonomy, nil
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
			return nil, nil
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
		taxonomy
	where
		parent_id is null`

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
		// fmt.Printf("Flattened: %s\n", t.Name)
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

func (repo *awsRepo) ReverseFlatten(t *cohesioned.Taxonomy) (*cohesioned.Taxonomy, error) {
	if t == nil {
		return nil, nil
	}

	parent, err := repo.Get(t.ParentID)
	if err != nil {
		return t, fmt.Errorf("Failed to get parent for %s %v\n", t.Name, err)
	}

	if parent == nil {
		return t, nil
	}
	t.Parent = parent
	flattenedParent, err := repo.ReverseFlatten(parent)
	if err != nil {
		return t, fmt.Errorf("Failed to flatten parent of %s: %v", t.Name, err)
	}

	t.Name = fmt.Sprintf("%s > %s", flattenedParent.Name, t.Name)
	return t, nil
}

func (repo *awsRepo) ListRecursive() ([]*cohesioned.Taxonomy, error) {
	var list []*cohesioned.Taxonomy

	list, err := repo.List()
	if err != nil {
		return list, err
	}

	for _, t := range list {
		children, err := repo.listChildrenRecursive(t)
		if err != nil {
			return list, fmt.Errorf("Failed to get listChildrenRecursive for %s %v\n", t.ID, err)
		}
		t.Children = children
	}

	return list, nil
}

func (repo *awsRepo) ListChildrenRecursive(parentID int64) ([]*cohesioned.Taxonomy, error) {
	parent, err := repo.Get(parentID)
	if err != nil {
		return nil, err
	}

	children, err := repo.listChildrenRecursive(parent)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ListChildrenRecursive for %s %v\n", parentID, err)
	}

	return children, nil
}

func (repo *awsRepo) listChildrenRecursive(t *cohesioned.Taxonomy) ([]*cohesioned.Taxonomy, error) {
	var list []*cohesioned.Taxonomy

	if t == nil {
		return list, nil
	}

	list, err := repo.ListChildren(t.ID)
	if err != nil {
		return list, err
	}

	t.Children = list
	for _, child := range t.Children {
		children, err := repo.listChildrenRecursive(child)
		if err != nil {
			return list, fmt.Errorf("Failed to get children for %s %v\n", child.ID, err)
		}

		child.Children = children
	}

	return list, nil
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
		if err == sql.ErrNoRows {
			return nil, err
		}

		return taxonomy, fmt.Errorf("faled to map row to taxonomy: %v", err)
	}

	taxonomy.ParentID = parentID.Int64
	taxonomy.Updated = updated.Time
	taxonomy.UpdatedBy = updatedBy.Int64

	return taxonomy, nil
}
