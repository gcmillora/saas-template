package provider

import (
	"log/slog"
	"os"
)
func NewLoggerProvider(env *EnvProvider) *slog.Logger {
	level := slog.LevelDebug

	if env.appEnv == "production" {
		level = slog.LevelInfo
	}

	loggerOpts := slog.HandlerOptions{
		Level: level,
	}

	withTextLogger := slog.NewTextHandler(os.Stdout, &loggerOpts)
	
	slog.SetDefault(slog.New(withTextLogger))

	return slog.New(withTextLogger)
}