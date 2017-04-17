package cohesioned

import "time"

//Profile represents a User of this system
type Profile struct {
	Email         string    `json:"email"`
	FullName      string    `json:"name"`
	FirstName     string    `json:"given_name"`
	LastName      string    `json:"family_name"`
	PictureURL    string    `json:"picture"`
	Locale        string    `json:"locale"`
	Nickname      string    `json:"nickname"`
	EmailVerified bool      `json:"email_verified"`
	UserID        string    `json:"user_id"`
	ClientID      string    `json:"client_id"`
	DateCreated   time.Time `json:"created_at"`
	LastUpdated   time.Time `json:"updated_at"`
}
