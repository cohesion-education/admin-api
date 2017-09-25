package profile

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

func (repo *awsRepo) Save(p *cohesioned.Profile) (int64, error) {
	sql := `insert into user
	(
		created,
		email,
		full_name,
		first_name,
		last_name,
		nickname,
		profile_pic_url,
		locale,
		enabled,
		verified,
		sub,
		state,
		county
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		p.Created,
		p.Email,
		p.FullName,
		p.FirstName,
		p.LastName,
		p.Nickname,
		p.PictureURL,
		p.Locale,
		p.Enabled,
		p.EmailVerified,
		p.Sub,
		p.State,
		p.County,
	)

	if err != nil {
		return 0, fmt.Errorf("Failed to insert user: %v", err)
	}

	profileID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id from result: %v", err)
	}

	return profileID, nil
}

func (repo *awsRepo) Update(p *cohesioned.Profile) error {
	sql := `update user set
		updated = ?,
		email = ?,
		full_name = ?,
		first_name = ?,
		last_name = ?,
		nickname = ?,
		profile_pic_url = ?,
		locale = ?,
		enabled = ?,
		verified = ?,
		sub = ?,
		state = ?,
		county = ?
	where id = ?`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		p.Updated,
		p.Email,
		p.FullName,
		p.FirstName,
		p.LastName,
		p.Nickname,
		p.PictureURL,
		p.Locale,
		p.Enabled,
		p.EmailVerified,
		p.Sub,
		p.State,
		p.County,
		p.ID,
	)

	if err != nil {
		return fmt.Errorf("Failed to update user: %v", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil || rowsEffected == 0 {
		return fmt.Errorf("Failed to update user: %v", err)
	}

	return nil
}

func (repo *awsRepo) FindByEmail(email string) (*cohesioned.Profile, error) {
	profile := new(cohesioned.Profile)

	query := `select
		id,
		created,
		updated,
		email,
		full_name,
		first_name,
		last_name,
		nickname,
		profile_pic_url,
		locale,
		enabled,
		verified,
		sub,
		state,
		county
	from user
		where email = ?`

	row := repo.QueryRow(query, email)

	var updated db.NullTime
	var nickname sql.NullString
	var pictureURL sql.NullString
	var locale sql.NullString
	var sub sql.NullString
	var state sql.NullString
	var county sql.NullString
	var enabled sql.NullInt64
	var verified sql.NullInt64

	err := row.Scan(
		&profile.ID,
		&profile.Created,
		&updated,
		&profile.Email,
		&profile.FullName,
		&profile.FirstName,
		&profile.LastName,
		&nickname,
		&pictureURL,
		&locale,
		&enabled,
		&verified,
		&sub,
		&state,
		&county,
	)

	if err != nil {
		//TODO - check error types

		return profile, err
	}

	profile.Updated = updated.Time
	profile.Nickname = nickname.String
	profile.PictureURL = pictureURL.String
	profile.Locale = locale.String
	profile.Sub = sub.String
	profile.State = state.String
	profile.County = county.String

	profile.Enabled = enabled.Int64 == 1
	profile.EmailVerified = verified.Int64 == 1

	return profile, nil
}
