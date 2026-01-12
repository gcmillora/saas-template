package provider

import (
	"github.com/gorilla/sessions"
)

func NewSessionProvider(env *EnvProvider) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(env.SessionSecret()))
}
