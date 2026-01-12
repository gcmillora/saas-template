package provider

import (
	"log"
	"os"
)
const (
	AppEnvLocal = "local"
	
)

type EnvProvider struct {
	appBaseUrl string
	appEnv string
	appSecret string
	databaseUrl string
	logLevel string
	serverPort string
	supabaseURL string
	supabaseKey string
	sessionSecret string
}


func (e *EnvProvider) AppEnv() string {
	return e.appEnv
}

func (e *EnvProvider) AppBaseUrl() string {
	return e.appBaseUrl
}

func (e *EnvProvider) AppSecret() string {
	return e.appSecret
}

func (e *EnvProvider) DatabaseUrl() string {
	return e.databaseUrl
}

func (e *EnvProvider) LogLevel() string {
	return e.logLevel
}

func (e *EnvProvider) ServerPort() string {
	return e.serverPort
}

func (e *EnvProvider) SupabaseURL() string {
	return e.supabaseURL
}

func (e *EnvProvider) SupabaseKey() string {
	return e.supabaseKey
}

func (e *EnvProvider) SessionSecret() string {
	return e.sessionSecret
}

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

	appServer := fallbackLookupEnv("APP_ENV", AppEnvLocal)
	serverPort := fallbackLookupEnv("SERVER_PORT", "8080")

	logLevel := fallbackLookupEnv("LOG_LEVEL", "error")
	appSecret:= requireLookupEnv("APP_SECRET")
	appBaseUrl := requireLookupEnv("APP_BASE_URL")
	databaseUrl := requireLookupEnv("DATABASE_URL")
	supabaseURL := requireLookupEnv("SUPABASE_URL")
	supabaseKey := requireLookupEnv("SUPABASE_KEY")
	sessionSecret := requireLookupEnv("SESSION_SECRET")

	envProvider := EnvProvider{
		appEnv: appServer,
		serverPort: serverPort,
		logLevel: logLevel,
		appSecret: appSecret,
		appBaseUrl: appBaseUrl,
		databaseUrl: databaseUrl,
		supabaseURL: supabaseURL,
		supabaseKey: supabaseKey,
		sessionSecret: sessionSecret,
	}

	return &envProvider
}