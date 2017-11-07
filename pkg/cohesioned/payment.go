package cohesioned

import "time"

type PaymentDetails struct {
	ID        int64              `json:"id"`
	Created   time.Time          `json:"created"`
	Updated   time.Time          `json:"updated"`
	CreatedBy int64              `json:"created_by"`
	UpdatedBy int64              `json:"updated_by"`
	Token     StripePaymentToken `json:"token"`
}

type StripePaymentToken struct {
	Created  int32      `json:"created"`
	ID       string     `json:"id"`
	Type     string     `json:"type"`
	ClientIP string     `json:"client_ip"`
	Used     bool       `json:"used"`
	LiveMode bool       `json:"livemode"`
	Card     StripeCard `json:"card,omitempty"`
}

type StripeCard struct {
	ID                 string `json:"id"`
	City               string `json:"address_city"`
	Brand              string `json:"brand"`
	AddressCountry     string `json:"address_country"`
	AddressLine1       string `json:"address_line1"`
	AddressLine1Check  string `json:"address_line1_check"`
	AddressLine2       string `json:"address_line2"`
	AddressLine2Check  string `json:"address_line2_check"`
	PostalCode         string `json:"address_zip"`
	PostalCodeCheck    string `json:"address_zip_check"`
	State              string `json:"address_state"`
	Country            string `json:"country"`
	Currency           string `json:"currency"`
	CvcCheck           string `json:"cvc_check"`
	ExpiryMonth        int8   `json:"exp_month"`
	ExpiryYear         int16  `json:"exp_year"`
	Funding            string `json:"funding"`
	LastFour           string `json:"last4"`
	DynamicLastFour    string `json:"dynamic_last4,omitempty"`
	Name               string `json:"name"`
	TokenizationMethod string `json:"tokenization_method,omitempty"`
	Fingerprint        string `json:"fingerprint"`
}
