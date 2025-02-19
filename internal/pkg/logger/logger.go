package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger zerolog.Logger

type contextKey string

const loggerKey = contextKey("logger")

// CustomLogger é uma interface para logging
type CustomLogger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	With(args ...interface{}) CustomLogger
}

type slogWrapper struct {
	logger *slog.Logger
}

func (s *slogWrapper) Info(msg string, args ...interface{}) {
	s.logger.Info(msg, args...)
}

func (s *slogWrapper) Error(msg string, args ...interface{}) {
	s.logger.Error(msg, args...)
}

func (s *slogWrapper) With(args ...interface{}) CustomLogger {
	return &slogWrapper{logger: s.logger.With(args...)}
}

// NewLogger cria uma nova instância do logger
func NewLogger() CustomLogger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &slogWrapper{logger: logger}
}

// FromContext retorna o logger do contexto
func FromContext(ctx context.Context) CustomLogger {
	if logger, ok := ctx.Value(loggerKey).(CustomLogger); ok {
		return logger
	}
	return NewLogger()
}

// WithContext adiciona o logger ao contexto
func WithContext(ctx context.Context, logger CustomLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Setup() {
	level := getLogLevel()

	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		globalLogger = log.Output(os.Stdout).Level(level)
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		globalLogger = log.Output(output).Level(level)
	}

	globalLogger = globalLogger.With().
		Str("app", "k8s-resource-analyzer").
		Str("env", os.Getenv("GIN_MODE")).
		Logger()
}

func getLogLevel() zerolog.Level {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func Info() *zerolog.Event {
	return globalLogger.Info()
}

func Error() *zerolog.Event {
	return globalLogger.Error()
}

func Fatal() *zerolog.Event {
	return globalLogger.Fatal()
}
