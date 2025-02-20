// Package logger fornece funcionalidades de logging estruturado.
// Este pacote implementa um wrapper sobre o zerolog para fornecer
// logs estruturados com níveis, campos contextuais e formatação consistente.
package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	// defaultLogger é a instância padrão do logger
	defaultLogger zerolog.Logger
)

// Config contém as configurações do logger
type Config struct {
	// Level define o nível mínimo de log
	Level string
	// Pretty indica se deve usar formatação legível
	Pretty bool
	// Output define o writer de saída
	Output io.Writer
}

// Field representa um campo de log estruturado
type Field struct {
	Key   string
	Value interface{}
}

// Setup inicializa o logger com as configurações fornecidas
func Setup(cfg *Config) {
	if cfg == nil {
		cfg = &Config{
			Level:  "info",
			Pretty: true,
			Output: os.Stdout,
		}
	}

	// Configura o nível de log
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configura o writer
	var w io.Writer = cfg.Output
	if cfg.Pretty {
		w = zerolog.ConsoleWriter{
			Out:        cfg.Output,
			TimeFormat: time.RFC3339,
		}
	}

	// Cria o logger
	defaultLogger = zerolog.New(w).With().Timestamp().Logger()
}

// WithContext retorna um logger com campos do contexto
func WithContext(ctx context.Context) *zerolog.Logger {
	return &defaultLogger
}

// WithFields adiciona campos ao logger
func WithFields(fields ...Field) *zerolog.Logger {
	ctx := defaultLogger.With()
	for _, f := range fields {
		ctx = ctx.Interface(f.Key, f.Value)
	}
	logger := ctx.Logger()
	return &logger
}

// Debug loga uma mensagem no nível debug
func Debug(msg string, fields ...Field) {
	event := defaultLogger.Debug()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

// Info loga uma mensagem no nível info
func Info(msg string, fields ...Field) {
	event := defaultLogger.Info()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

// Warn loga uma mensagem no nível warn
func Warn(msg string, fields ...Field) {
	event := defaultLogger.Warn()
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

// Error loga uma mensagem no nível error
func Error(msg string, err error, fields ...Field) {
	event := defaultLogger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

// Fatal loga uma mensagem no nível fatal e encerra o programa
func Fatal(msg string, err error, fields ...Field) {
	event := defaultLogger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for _, f := range fields {
		event = event.Interface(f.Key, f.Value)
	}
	event.Msg(msg)
}

// NewField cria um novo campo de log
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
