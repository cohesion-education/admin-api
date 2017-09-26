package cohesioned

import (
	"fmt"
	"strings"
	"time"
)

type Preferences struct {
	Newsletter  bool `json:"newsletter"`
	BetaProgram bool `json:"beta_program"`
}

type Student struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	School    string    `json:"school"`
	ParentID  int64     `json:"parent_id"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

func (s *Student) String() string {
	return fmt.Sprintf("ID: %d Name: %s Grade: %s School: %s Parent ID: %d", s.ID, s.Name, s.Grade, s.School, s.ParentID)
}

//Profile represents a User of this system
type Profile struct {
	ID            int64     `json:"id"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Enabled       bool      `json:"enabled"`
	EmailVerified bool      `json:"email_verified"`
	Email         string    `json:"email"`
	FullName      string    `json:"name"`
	FirstName     string    `json:"given_name"`
	LastName      string    `json:"family_name"`
	Nickname      string    `json:"nickname"`
	PictureURL    string    `json:"picture"`
	Locale        string    `json:"locale"`

	// UserID        string      `json:"user_id"`
	// ClientID      string      `json:"client_id"`
	Sub         string      `json:"sub"`
	Preferences Preferences `json:"preferences"`
	State       string      `json:"state"`
	County      string      `json:"county"`
	Students    []*Student  `json:"students"`
}

func (p *Profile) String() string {
	return fmt.Sprintf("ID: %d Full Name: %s", p.ID, p.FullName)
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
