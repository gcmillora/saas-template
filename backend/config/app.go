package config

import (
	"adobo/config/provider"
	"log/slog"
	"path"
	"runtime"

	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
)


type App struct {
	env *provider.EnvProvider
	db *provider.DbProvider
	supabase *provider.SupabaseProvider
	rootDir string
	logger *slog.Logger
	cache *cache.Cache
	session *sessions.CookieStore
}

func (app *App) Session() *sessions.CookieStore {
	if app.session == nil {
		app.session = provider.NewSessionProvider(app.env)
	}
	return app.session
}

func (app *App) EnvVars() *provider.EnvProvider {
	if app.env == nil {
		app.env = provider.NewEnvProvider(app.rootDir)
	}
	return app.env
}

func (app *App) DB() *provider.DbProvider {
	if app.db == nil {
		app.db = provider.NewDbProvider(app.env)
	}
	return app.db
}

func (app *App) Logger() *slog.Logger {
	if app.logger == nil {
		app.logger = provider.NewLoggerProvider(app.env)
	}
	return app.logger
}

func (app *App) Cache() *cache.Cache {
	if app.cache == nil {
		app.cache = provider.NewCacheProvider()
	}
	return app.cache
}

func (app *App) Supabase() *provider.SupabaseProvider {
	if app.supabase == nil {
		app.supabase = provider.NewSupabaseProvider(app.env)
	}
	return app.supabase
}

func (app *App) setRootDir() {
	_, b, _, _ := runtime.Caller(0)
	app.rootDir = path.Join(path.Dir(b), "..")
}

func NewApp() *App {
	app := App{}

	app.env = provider.NewEnvProvider(app.rootDir)
	app.db = provider.NewDbProvider(app.env)
	app.setRootDir()

	provider.NewValidationProvider()
	app.supabase = provider.NewSupabaseProvider(app.env)
	app.session = provider.NewSessionProvider(app.env)

	return &app
}