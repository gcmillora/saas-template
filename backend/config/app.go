package config

import (
	"log/slog"
	"path"
	"runtime"
	"saas-template/config/provider"

	"github.com/gorilla/sessions"
	"github.com/patrickmn/go-cache"
	"github.com/resend/resend-go/v2"
	storage_go "github.com/supabase-community/storage-go"
)

type App struct {
	env     *provider.EnvProvider
	db      *provider.DbProvider
	rootDir string
	logger  *slog.Logger
	cache   *cache.Cache
	session *sessions.CookieStore
	storage *storage_go.Client
	resend  *resend.Client
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

func (app *App) Storage() *storage_go.Client {
	if app.storage == nil {
		app.storage = provider.NewSupabaseStorageClient(app.env)
	}
	return app.storage
}

func (app *App) Resend() *resend.Client {
	if app.resend == nil {
		app.resend = provider.NewResendProvider(app.env)
	}
	return app.resend
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
	app.session = provider.NewSessionProvider(app.env)

	return &app
}
