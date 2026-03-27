package provider

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func NewSessionProvider(env *EnvProvider) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(env.SessionSecret()))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   env.AppEnv() != AppEnvLocal,
		SameSite: http.SameSiteLaxMode,
	}
	return store
}
