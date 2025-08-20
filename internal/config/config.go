package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Port        string `yaml:"port"`
	Environment string `yaml:"environment"`
}

type DatabaseConfig struct {
	URL                string `yaml:"url"`
	MaxConnections     int    `yaml:"max_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load загружает конфигурацию из файла YAML или переменных окружения
func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	// Пытаемся загрузить из YAML файла
	config := &Config{}

	configPath := getEnv("CONFIG_PATH", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		if err := loadFromYAML(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to load config from YAML: %w", err)
		}
	}

	// Переопределяем значения из переменных окружения
	overrideFromEnv(config)

	// Устанавливаем значения по умолчанию если они не установлены
	setDefaults(config)

	return config, nil
}

func loadFromYAML(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(config)
}

func overrideFromEnv(config *Config) {
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Server.Environment = env
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Database.URL = dbURL
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	if logFormat := os.Getenv("LOG_FORMAT"); logFormat != "" {
		config.Logging.Format = logFormat
	}
}

func setDefaults(config *Config) {
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Environment == "" {
		config.Server.Environment = "development"
	}
	if config.Database.URL == "" {
		config.Database.URL = "postgres://user:password@localhost:5432/subscriptions?sslmode=disable"
	}
	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 25
	}
	if config.Database.MaxIdleConnections == 0 {
		config.Database.MaxIdleConnections = 5
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		if config.Server.Environment == "development" {
			config.Logging.Format = "text"
		} else {
			config.Logging.Format = "json"
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetConfigPath возвращает путь к конфигурационному файлу
func GetConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// Ищем config.yaml в текущей директории и в parent directories
	paths := []string{
		"config.yaml",
		"../config.yaml",
		"../../config.yaml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	return "config.yaml"
}
