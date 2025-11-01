package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoadConfig загружает конфигурацию из переменных окружения и файла
func LoadConfig() (*Config, error) {
	cfg := GetConfig()

	// Загружаем из переменных окружения
	if err := loadFromEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// Загружаем из файла конфигурации, если он существует
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		if err := loadFromFile(cfg, configFile); err != nil {
			return nil, fmt.Errorf("failed to load config from file %s: %w", configFile, err)
		}
	}

	// Валидируем конфигурацию
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// loadFromEnv загружает конфигурацию из переменных окружения
func loadFromEnv(cfg *Config) error {
	return setFieldsFromEnv(reflect.ValueOf(cfg).Elem(), reflect.TypeOf(cfg).Elem())
}

// setFieldsFromEnv рекурсивно устанавливает значения полей из переменных окружения
func setFieldsFromEnv(v reflect.Value, t reflect.Type) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := setFieldsFromEnv(field, fieldType.Type); err != nil {
				return err
			}
			continue
		}

		envTag := fieldType.Tag.Get("env")
		defaultTag := fieldType.Tag.Get("default")

		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" && defaultTag != "" {
			envValue = defaultTag
		}

		if envValue == "" {
			continue
		}

		if err := setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue устанавливает значение поля в зависимости от его типа
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			values := strings.Split(value, ",")
			for i, v := range values {
				values[i] = strings.TrimSpace(v)
			}
			field.Set(reflect.ValueOf(values))
		}
	case reflect.TypeOf(time.Duration(0)).Kind():
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(duration))
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// loadFromFile загружает конфигурацию из JSON файла
func loadFromFile(cfg *Config, filename string) error {
	// Проверяем, что путь не содержит опасных символов
	if strings.Contains(filename, "..") {
		return fmt.Errorf("invalid file path: %s", filename)
	}

	data, err := os.ReadFile(filename) // #nosec G304 - путь проверен выше
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
}

// validateConfig валидирует конфигурацию
func validateConfig(cfg *Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
	}

	validLogLevels := map[string]bool{
		"trace": true, "debug": true, "info": true,
		"warn": true, "error": true, "fatal": true, "panic": true,
	}

	if !validLogLevels[cfg.Logger.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Logger.Level)
	}

	return nil
}
