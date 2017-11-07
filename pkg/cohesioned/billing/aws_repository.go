package billing

import (
	"database/sql"
	"fmt"

	"github.com/cohesion-education/api/pkg/cohesioned"
)

type awsRepo struct {
	*sql.DB
}

func NewAwsRepo(db *sql.DB) Repo {
	return &awsRepo{
		DB: db,
	}
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
		card_address_zip,
		card_address_zip_check
	) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

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
	//TODO - implement me!
	return fmt.Errorf("Not yet implemented")
}
