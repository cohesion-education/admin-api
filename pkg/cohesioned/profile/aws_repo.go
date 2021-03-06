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

func (repo *awsRepo) List() ([]*cohesioned.Profile, error) {
	var list []*cohesioned.Profile

	selectQuery := `select
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
			beta_program,
			newsletter,
			sub,
			state,
			county,
			onboarded,
			billing_status,
			trial_start
		from
			user`

	rows, err := repo.Query(selectQuery)
	if err != nil {
		return list, fmt.Errorf("Failed to execute query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		p, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the result set from the db: %v", err)
		}

		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return list, fmt.Errorf("rows had an error: %v", err)
	}

	return list, nil
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
		beta_program,
		newsletter,
		sub,
		state,
		county,
		onboarded,
		billing_status,
		trial_start
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		p.Preferences.BetaProgram,
		p.Preferences.Newsletter,
		p.Sub,
		p.State,
		p.County,
		p.Onboarded,
		p.BillingStatus,
		p.TrialStart,
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
		beta_program = ?,
		newsletter = ?,
		sub = ?,
		state = ?,
		county = ?,
		onboarded = ?,
		billing_status = ?,
		trial_start = ?
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
		p.Preferences.BetaProgram,
		p.Preferences.Newsletter,
		p.Sub,
		p.State,
		p.County,
		p.Onboarded,
		p.BillingStatus,
		p.TrialStart,
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
		beta_program,
		newsletter,
		sub,
		state,
		county,
		onboarded,
		billing_status,
		trial_start
	from user
		where email = ?`

	row := repo.QueryRow(query, email)

	p, err := repo.mapRowToObject(row)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		default:
			return nil, fmt.Errorf("Unexpected error querying for user by email %s: %v", email, err)
		}
	}

	return p, nil
}

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.Profile, error) {
	profile := new(cohesioned.Profile)

	var updated db.NullTime
	var nickname sql.NullString
	var pictureURL sql.NullString
	var locale sql.NullString
	var sub sql.NullString
	var state sql.NullString
	var county sql.NullString
	var enabled sql.NullBool
	var verified sql.NullBool
	var betaProgram sql.NullBool
	var newsletter sql.NullBool
	var onboarded sql.NullBool
	var billingStatus sql.NullString
	var trialStart db.NullTime

	err := rs.Scan(
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
		&betaProgram,
		&newsletter,
		&sub,
		&state,
		&county,
		&onboarded,
		&billingStatus,
		&trialStart,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	profile.Updated = updated.Time
	profile.Nickname = nickname.String
	profile.PictureURL = pictureURL.String
	profile.Locale = locale.String
	profile.Sub = sub.String
	profile.State = state.String
	profile.County = county.String

	profile.Enabled = enabled.Bool
	profile.EmailVerified = verified.Bool
	profile.Preferences.BetaProgram = betaProgram.Bool
	profile.Preferences.Newsletter = newsletter.Bool

	profile.Onboarded = onboarded.Bool
	profile.BillingStatus = billingStatus.String
	profile.TrialStart = trialStart.Time

	return profile, nil
}
