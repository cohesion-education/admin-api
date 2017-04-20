package cohesioned

import (
	"fmt"
	"net/http"
	"reflect"
)

type key int

const (
	AuthSessionCookieName string = "auth-session"
	CurrentUserKey        key    = iota
	CurrentUserSessionKey string = "profile"
)

func GetProfile(req *http.Request) (*Profile, error) {
	profile, ok := req.Context().Value(CurrentUserKey).(*Profile)
	if !ok {
		return nil, fmt.Errorf("profile not of the proper type: %s", reflect.TypeOf(profile).String())
	}

	return profile, nil
}
