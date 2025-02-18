package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

func Setup() {
	level := getLogLevel()

	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		Logger = log.Output(os.Stdout).Level(level)
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		Logger = log.Output(output).Level(level)
	}

	Logger = Logger.With().
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
	return Logger.Info()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
