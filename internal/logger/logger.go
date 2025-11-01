package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"dolina-flower-order-backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger представляет кастомный логгер
type Logger struct {
	logger zerolog.Logger
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger возвращает единственный экземпляр логгера (синглтон)
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{}
	})
	return instance
}

// Initialize инициализирует логгер с конфигурацией
func (l *Logger) Initialize(cfg *config.Config) error {
	// Устанавливаем уровень логирования
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// Настраиваем вывод
	var output io.Writer
	switch cfg.Logger.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		file, err := os.OpenFile(cfg.Logger.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		output = file
	}

	// Настраиваем формат
	if cfg.Logger.Format == "console" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
	}

	// Создаем логгер с контекстом
	l.logger = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Str("service", "dolina-flower-order").
		Logger()

	// Устанавливаем глобальный логгер
	log.Logger = l.logger

	return nil
}

// WithContext добавляет контекст к логгеру
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger.With().Ctx(ctx).Logger(),
	}
}

// WithField добавляет поле к логгеру
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

// WithFields добавляет несколько полей к логгеру
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	logger := l.logger.With()
	for key, value := range fields {
		logger = logger.Interface(key, value)
	}
	return &Logger{
		logger: logger.Logger(),
	}
}

// WithError добавляет ошибку к логгеру
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Err(err).Logger(),
	}
}

// Trace логирует сообщение уровня trace
func (l *Logger) Trace(msg string) {
	l.logger.Trace().Msg(msg)
}

// Tracef логирует форматированное сообщение уровня trace
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

// Debug логирует сообщение уровня debug
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf логирует форматированное сообщение уровня debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info логирует сообщение уровня info
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof логирует форматированное сообщение уровня info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn логирует сообщение уровня warn
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf логирует форматированное сообщение уровня warn
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error логирует сообщение уровня error
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf логирует форматированное сообщение уровня error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatal логирует сообщение уровня fatal и завершает программу
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf логирует форматированное сообщение уровня fatal и завершает программу
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// LogRequest логирует HTTP запрос
func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	l.WithContext(ctx).
		WithFields(map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
		}).
		Info("HTTP request processed")
}

// LogError логирует ошибку с контекстом
func (l *Logger) LogError(ctx context.Context, err error, msg string) {
	l.WithContext(ctx).WithError(err).Error(msg)
}
