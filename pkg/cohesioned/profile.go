package cohesioned

import "time"

//AppMetadata mapping of metadata provided by auth0
type AppMetadata struct {
	Roles []string `json:"roles"`
}

type Preferences struct {
	Newsletter  bool `json:"newsletter"`
	BetaProgram bool `json:"beta_program"`
}

//Profile represents a User of this system
type Profile struct {
	Auditable
	Enabled       bool        `json:"enabled"`
	Email         string      `json:"email"`
	FullName      string      `json:"name"`
	FirstName     string      `json:"given_name"`
	LastName      string      `json:"family_name"`
	PictureURL    string      `json:"picture"`
	Locale        string      `json:"locale"`
	Nickname      string      `json:"nickname"`
	EmailVerified bool        `json:"email_verified"`
	UserID        string      `json:"user_id"`
	ClientID      string      `json:"client_id"`
	DateCreated   time.Time   `json:"created_at"`
	LastUpdated   time.Time   `json:"updated_at"`
	Metadata      AppMetadata `json:"app_metadata"`
	Preferences   Preferences `json:"preferences"`
}

func (p *Profile) HasRole(roleName string) bool {
	if len(roleName) == 0 {
		return false
	}

	if len(p.Metadata.Roles) == 0 {
		return false
	}

	for _, role := range p.Metadata.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}
