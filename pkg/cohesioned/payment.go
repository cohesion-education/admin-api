package cohesioned

import "time"

type StripePaymentToken struct {
	Created   time.Time  `json:"created"`
	Updated   time.Time  `json:"updated"`
	CreatedBy int64      `json:"created_by"`
	UpdatedBy int64      `json:"updated_by"`
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Card      StripeCard `json:"card,omitempty"`
}

type StripeCard struct {
	ID                 string `json:"id"`
	City               string `json:"address_city"`
	Brand              string `json:"brand"`
	AddressCountry     string `json:"address_country"`
	AddressLine1       string `json:"address_line1"`
	AddressLine1Check  string `json:"address_line1_check"`
	AddressLine2       string `json:"address_line2"`
	AddressLine2Check  string `json:"address_state"`
	PostalCode         string `json:"address_zip"`
	PostalCodeCheck    string `json:"address_zip_check"`
	Country            string `json:"country"`
	Currency           string `json:"currency"`
	CvcCheck           string `json:"cvc_check"`
	DynamicLastFour    string `json:"dynamic_last4"`
	ExpiryMonth        int8   `json:"exp_month"`
	ExpiryYear         int16  `json:"exp_year"`
	Funding            string `json:"funding"`
	LastFour           string `json:"last4"`
	Name               string `json:"name"`
	TokenizationMethod string `json:"tokenization_method"`
	Fingerprint        string `json:"fingerprint"`
}
