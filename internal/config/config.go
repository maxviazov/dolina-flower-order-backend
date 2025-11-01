package config

import (
	"fmt"
	"sync"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logger   LoggerConfig   `json:"logger"`
	Security SecurityConfig `json:"security"`
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Host         string        `json:"host" env:"SERVER_HOST" default:"localhost"`
	Port         int           `json:"port" env:"SERVER_PORT" default:"8081"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `json:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" default:"60s"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host         string `json:"host" env:"DB_HOST" default:"localhost"`
	Port         int    `json:"port" env:"DB_PORT" default:"5432"`
	Name         string `json:"name" env:"DB_NAME" default:"dolina_flowers"`
	User         string `json:"user" env:"DB_USER" default:"postgres"`
	Password     string `json:"password" env:"DB_PASSWORD"`
	SSLMode      string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns int    `json:"max_open_conns" env:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns int    `json:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" default:"5"`
}

// LoggerConfig конфигурация логгера
type LoggerConfig struct {
	Level      string `json:"level" env:"LOG_LEVEL" default:"info"`
	Format     string `json:"format" env:"LOG_FORMAT" default:"console"`
	Output     string `json:"output" env:"LOG_OUTPUT" default:"stdout"`
	TimeFormat string `json:"time_format" env:"LOG_TIME_FORMAT" default:"2006-01-02T15:04:05.000Z07:00"`
}

// SecurityConfig конфигурация безопасности
type SecurityConfig struct {
	JWTSecret     string        `json:"jwt_secret" env:"JWT_SECRET"`
	JWTExpiration time.Duration `json:"jwt_expiration" env:"JWT_EXPIRATION" default:"24h"`
	CORSOrigins   []string      `json:"cors_origins" env:"CORS_ORIGINS" default:"*"`
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig возвращает единственный экземпляр конфигурации (синглтон)
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

// IsProduction проверяет, запущено ли приложение в продакшене
func (c *Config) IsProduction() bool {
	return c.Logger.Level == "error" || c.Logger.Level == "warn"
}

// GetServerAddress возвращает полный адрес сервера
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
