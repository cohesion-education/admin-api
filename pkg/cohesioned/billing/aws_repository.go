package billing

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

func (repo *awsRepo) List() ([]*cohesioned.PaymentDetails, error) {
	var list []*cohesioned.PaymentDetails

	selectQuery := `select
		id,
		created,
		created_by,
		updated,
		updated_by,
		token_created,
		token_id,
		token_client_ip,
		token_used,
		token_live_mode,
		token_type,
		card_id,
		card_brand,
		card_funding,
		card_tokenization_method,
		card_fingerprint,
		card_name,
		card_country,
		card_currency,
		card_exp_month,
		card_exp_year,
		card_cvc_check,
		card_last4,
		card_dynamic_last4,
		card_address_line1,
		card_address_line1_check,
		card_address_line2,
		card_address_line2_check,
		card_address_country,
		card_address_state,
		card_address_city,
		card_address_zip,
		card_address_zip_check
	from
		payment_detail`

	rows, err := repo.Query(selectQuery)
	if err != nil {
		return list, fmt.Errorf("Failed to execute list payment detail query: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		p, err := repo.mapRowToObject(rows)
		if err != nil {
			return list, fmt.Errorf("an unexpected error occurred while processing the list payment detail result set from the db: %v", err)
		}

		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return list, fmt.Errorf("list payment detail rows had an error: %v", err)
	}

	return list, nil
}

func (repo *awsRepo) FindByCreatedByID(id int64) (*cohesioned.PaymentDetails, error) {
	selectQuery := `select
		id,
		created,
		created_by,
		updated,
		updated_by,
		token_created,
		token_id,
		token_client_ip,
		token_used,
		token_live_mode,
		token_type,
		card_id,
		card_brand,
		card_funding,
		card_tokenization_method,
		card_fingerprint,
		card_name,
		card_country,
		card_currency,
		card_exp_month,
		card_exp_year,
		card_cvc_check,
		card_last4,
		card_dynamic_last4,
		card_address_line1,
		card_address_line1_check,
		card_address_line2,
		card_address_line2_check,
		card_address_country,
		card_address_state,
		card_address_city,
		card_address_zip,
		card_address_zip_check
	from
		payment_detail
	where
		created_by = ?`

	row := repo.QueryRow(selectQuery, id)
	paymentDetail, err := repo.mapRowToObject(row)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		default:
			return nil, fmt.Errorf("Unexpected error querying for payment detail by user id %d: %v", id, err)
		}
	}

	return paymentDetail, nil
}

func (repo *awsRepo) Save(p *cohesioned.PaymentDetails) (int64, error) {
	sql := `insert into payment_detail
	(
		created,
		created_by,
		token_created,
		token_id,
		token_client_ip,
		token_used,
		token_live_mode,
		token_type,
		card_id,
		card_brand,
		card_funding,
		card_tokenization_method,
		card_fingerprint,
		card_name,
		card_country,
		card_currency,
		card_exp_month,
		card_exp_year,
		card_cvc_check,
		card_last4,
		card_dynamic_last4,
		card_address_line1,
		card_address_line1_check,
		card_address_line2,
		card_address_line2_check,
		card_address_country,
		card_address_state,
		card_address_city,
		card_address_zip,
		card_address_zip_check
	) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return 0, fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		p.Created,
		p.CreatedBy,
		p.Token.Created,
		p.Token.ID,
		p.Token.ClientIP,
		p.Token.Used,
		p.Token.LiveMode,
		p.Token.Type,
		p.Token.Card.ID,
		p.Token.Card.Brand,
		p.Token.Card.Funding,
		p.Token.Card.TokenizationMethod,
		p.Token.Card.Fingerprint,
		p.Token.Card.Name,
		p.Token.Card.Country,
		p.Token.Card.Currency,
		p.Token.Card.ExpiryMonth,
		p.Token.Card.ExpiryYear,
		p.Token.Card.CvcCheck,
		p.Token.Card.LastFour,
		p.Token.Card.DynamicLastFour,
		p.Token.Card.AddressLine1,
		p.Token.Card.AddressLine1Check,
		p.Token.Card.AddressLine2,
		p.Token.Card.AddressLine2Check,
		p.Token.Card.Country,
		p.Token.Card.State,
		p.Token.Card.City,
		p.Token.Card.PostalCode,
		p.Token.Card.PostalCodeCheck,
	)

	if err != nil {
		return 0, fmt.Errorf("Failed to insert payment details: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id from result: %v", err)
	}

	return id, nil
}

func (repo *awsRepo) Update(p *cohesioned.PaymentDetails) error {
	sql := `update payment_detail set
		updated = ?,
		updated_by = ?,
		token_created = ?,
		token_id = ?,
		token_client_ip = ?,
		token_used = ?,
		token_live_mode = ?,
		token_type = ?,
		card_id = ?,
		card_brand = ?,
		card_funding = ?,
		card_tokenization_method = ?,
		card_fingerprint = ?,
		card_name = ?,
		card_country = ?,
		card_currency = ?,
		card_exp_month = ?,
		card_exp_year = ?,
		card_cvc_check = ?,
		card_last4 = ?,
		card_dynamic_last4 = ?,
		card_address_line1 = ?,
		card_address_line1_check = ?,
		card_address_line2 = ?,
		card_address_line2_check = ?,
		card_address_country = ?,
		card_address_state = ?,
		card_address_city = ?,
		card_address_zip = ?,
		card_address_zip_check = ?
	where
		id = ?`

	stmt, err := repo.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement %s: %v", sql, err)
	}

	result, err := stmt.Exec(
		p.Updated,
		p.UpdatedBy,
		p.Token.Created,
		p.Token.ID,
		p.Token.ClientIP,
		p.Token.Used,
		p.Token.LiveMode,
		p.Token.Type,
		p.Token.Card.ID,
		p.Token.Card.Brand,
		p.Token.Card.Funding,
		p.Token.Card.TokenizationMethod,
		p.Token.Card.Fingerprint,
		p.Token.Card.Name,
		p.Token.Card.Country,
		p.Token.Card.Currency,
		p.Token.Card.ExpiryMonth,
		p.Token.Card.ExpiryYear,
		p.Token.Card.CvcCheck,
		p.Token.Card.LastFour,
		p.Token.Card.DynamicLastFour,
		p.Token.Card.AddressLine1,
		p.Token.Card.AddressLine1Check,
		p.Token.Card.AddressLine2,
		p.Token.Card.AddressLine2Check,
		p.Token.Card.Country,
		p.Token.Card.State,
		p.Token.Card.City,
		p.Token.Card.PostalCode,
		p.Token.Card.PostalCodeCheck,
		p.ID,
	)

	if err != nil {
		fmt.Errorf("Failed to update payment details: %v", err)
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil || rowsEffected == 0 {
		return fmt.Errorf("Failed to update payment details: %v", err)
	}

	return nil
}

func (repo *awsRepo) mapRowToObject(rs db.RowScanner) (*cohesioned.PaymentDetails, error) {
	pd := &cohesioned.PaymentDetails{}
	var updated db.NullTime
	var updatedBy, cardExpiryMonth, cardExpiryYear sql.NullInt64
	var tokenUsed, tokenLiveMode sql.NullBool

	var tokenType, cardBrand, cardFunding, cardTokenizationMethod, cardFingerprint, name, cardCountry, cardCurrency, cardCvcCheck, cardLast4, cardDynamicLast4, cardAddressLine1, cardAddressLine1Check, cardAddressLine2, cardAddressLine2Check, cardAddressCountry, cardAddressState, cardAddressCity, cardAddressZip, cardAddressZipCheck sql.NullString

	err := rs.Scan(
		&pd.ID,
		&pd.Created,
		&pd.CreatedBy,
		&updated,
		&updatedBy,
		&pd.Token.Created,
		&pd.Token.ID,
		&pd.Token.ClientIP,
		&tokenUsed,
		&tokenLiveMode,
		&tokenType,
		&pd.Token.Card.ID,
		&cardBrand,
		&cardFunding,
		&cardTokenizationMethod,
		&cardFingerprint,
		&name,
		&cardCountry,
		&cardCurrency,
		&cardExpiryMonth,
		&cardExpiryYear,
		&cardCvcCheck,
		&cardLast4,
		&cardDynamicLast4,
		&cardAddressLine1,
		&cardAddressLine1Check,
		&cardAddressLine2,
		&cardAddressLine2Check,
		&cardAddressCountry,
		&cardAddressState,
		&cardAddressCity,
		&cardAddressZip,
		&cardAddressZipCheck,
	)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return pd, err
		default:
			return pd, fmt.Errorf("failed to map row to paymentdetail: %v", err)
		}
	}

	pd.Updated = updated.Time
	pd.UpdatedBy = updatedBy.Int64
	pd.Token.Type = tokenType.String
	pd.Token.LiveMode = tokenLiveMode.Bool
	pd.Token.Used = tokenLiveMode.Bool
	pd.Token.Card.Brand = cardBrand.String
	pd.Token.Card.Funding = cardFunding.String
	pd.Token.Card.TokenizationMethod = cardTokenizationMethod.String
	pd.Token.Card.Name = name.String
	pd.Token.Card.Country = cardCountry.String
	pd.Token.Card.Currency = cardCurrency.String
	pd.Token.Card.ExpiryMonth = int8(cardExpiryMonth.Int64)
	pd.Token.Card.ExpiryYear = int16(cardExpiryYear.Int64)
	pd.Token.Card.CvcCheck = cardCvcCheck.String
	pd.Token.Card.LastFour = cardLast4.String
	pd.Token.Card.DynamicLastFour = cardDynamicLast4.String
	pd.Token.Card.AddressLine1 = cardAddressLine1.String
	pd.Token.Card.AddressLine1Check = cardAddressLine1Check.String
	pd.Token.Card.AddressLine2 = cardAddressLine2.String
	pd.Token.Card.AddressLine2Check = cardAddressLine2Check.String
	pd.Token.Card.Country = cardAddressCountry.String
	pd.Token.Card.State = cardAddressState.String
	pd.Token.Card.City = cardAddressCity.String
	pd.Token.Card.PostalCode = cardAddressZip.String
	pd.Token.Card.PostalCodeCheck = cardAddressZipCheck.String

	return pd, nil
}
