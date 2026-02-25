package provider

import (
	"log"
	"os"
)

const (
	AppEnvLocal = "local"
)

type EnvProvider struct {
	appBaseUrl         string
	appEnv             string
	appSecret          string
	databaseUrl        string
	logLevel           string
	serverPort         string
	sessionSecret      string
	googleClientID     string
	googleClientSecret string
	githubClientID     string
	githubClientSecret string
}

func (e *EnvProvider) AppEnv() string             { return e.appEnv }
func (e *EnvProvider) AppBaseUrl() string          { return e.appBaseUrl }
func (e *EnvProvider) AppSecret() string           { return e.appSecret }
func (e *EnvProvider) DatabaseUrl() string         { return e.databaseUrl }
func (e *EnvProvider) LogLevel() string            { return e.logLevel }
func (e *EnvProvider) ServerPort() string          { return e.serverPort }
func (e *EnvProvider) SessionSecret() string       { return e.sessionSecret }
func (e *EnvProvider) GoogleClientID() string      { return e.googleClientID }
func (e *EnvProvider) GoogleClientSecret() string  { return e.googleClientSecret }
func (e *EnvProvider) GithubClientID() string      { return e.githubClientID }
func (e *EnvProvider) GithubClientSecret() string  { return e.githubClientSecret }

func NewEnvProvider(rootDir string) *EnvProvider {
	fallbackLookupEnv := func(key string, fallback string) string {
		value, exists := os.LookupEnv(key)
		if !exists {
			return fallback
		}
		return value
	}

	requireLookupEnv := func(key string) string {
		value, exists := os.LookupEnv(key)
		if !exists {
			log.Fatal("environment variable " + key + " is not set")
		}
		return value
	}

	return &EnvProvider{
		appEnv:             fallbackLookupEnv("APP_ENV", AppEnvLocal),
		serverPort:         fallbackLookupEnv("SERVER_PORT", "8080"),
		logLevel:           fallbackLookupEnv("LOG_LEVEL", "error"),
		appSecret:          requireLookupEnv("APP_SECRET"),
		appBaseUrl:         requireLookupEnv("APP_BASE_URL"),
		databaseUrl:        requireLookupEnv("DATABASE_URL"),
		sessionSecret:      requireLookupEnv("SESSION_SECRET"),
		googleClientID:     fallbackLookupEnv("GOOGLE_CLIENT_ID", ""),
		googleClientSecret: fallbackLookupEnv("GOOGLE_CLIENT_SECRET", ""),
		githubClientID:     fallbackLookupEnv("GITHUB_CLIENT_ID", ""),
		githubClientSecret: fallbackLookupEnv("GITHUB_CLIENT_SECRET", ""),
	}
}
