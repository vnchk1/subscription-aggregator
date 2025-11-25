package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Logger   LoggerConfig
	Server   ServerConfig
	Database DatabaseConfig
}

type LoggerConfig struct {
	LogLevel string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	Name           string
	User           string
	Password       string
	SSLMode        string
	MaxConnections int
	MigrationPath  string
}

func Load() (*Config, error) {
	dbConfig := DatabaseConfig{
		Host:           getEnv("DB_HOST", "localhost"),
		Port:           getEnv("DB_PORT", "5432"),
		Name:           getEnv("DB_NAME", "sub_aggregator"),
		User:           getEnv("DB_USER", "postgres"),
		Password:       getEnv("DB_PASSWORD", "postgres"),
		SSLMode:        getEnv("DB_SSLMODE", "disable"),
		MaxConnections: getEnvAsInt("DATABASE_MAX_CONNECTIONS", 10),
		MigrationPath:  getEnv("MIGRATION_PATH", "./migrations"),
	}

	// Собираем URL из отдельных параметров
	dbConfig.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSLMode,
	)

	return &Config{
		Logger: LoggerConfig{
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
		},
		Database: dbConfig,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}
