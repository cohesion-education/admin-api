package cohesioned

import (
	"strings"
)

type Preferences struct {
	Newsletter  bool `json:"newsletter"`
	BetaProgram bool `json:"beta_program"`
}

type Student struct {
	Name   string `json:"name"`
	Grade  string `json:"grade"`
	School string `json:"school"`
}

//Profile represents a User of this system
type Profile struct {
	Auditable
	Enabled       bool        `json:"enabled" datastore:"Enabled,omitempty"`
	Email         string      `json:"email" datastore:"Email,omitempty"`
	FullName      string      `json:"name" datastore:"FullName,omitempty"`
	FirstName     string      `json:"given_name" datastore:"FirstName,omitempty"`
	LastName      string      `json:"family_name" datastore:"LastName,omitempty"`
	PictureURL    string      `json:"picture" datastore:"PictureURL,omitempty"`
	Locale        string      `json:"locale" datastore:"Locale,omitempty"`
	Nickname      string      `json:"nickname" datastore:"Nickname,omitempty"`
	EmailVerified bool        `json:"email_verified" datastore:"EmailVerified,omitempty"`
	UserID        string      `json:"user_id" datastore:"UserID,omitempty"`
	ClientID      string      `json:"client_id" datastore:"ClientID,omitempty"`
	Sub           string      `json:"sub" datastore:"Sub,omitempty"`
	Preferences   Preferences `json:"preferences" datastore:"Preferences,omitempty"`
	State         string      `json:"state" datastore:"State,omitempty"`
	County        string      `json:"county" datastore:"County,omitempty"`
	Students      []Student   `json:"students" datastore:"Students,omitempty"`
}

//IsAdmin returns true if the user has a verified email address in the cohesioned.io domain
func (p *Profile) IsAdmin() bool {
	if len(p.Email) == 0 {
		return false
	}

	if !p.EmailVerified {
		return false
	}

	return strings.HasSuffix(p.Email, "@cohesioned.io")
}
