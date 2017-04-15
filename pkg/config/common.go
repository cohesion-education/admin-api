package config

type key int

const (
	AuthSessionCookieName string = "auth-session"
	CurrentUserKey        key    = iota
	CurrentUserSessionKey string = "profile"
)
