package config

import "os"

type Config struct {
	Port string
	DB   *DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func LoadConfig() (*Config, error) {
	// TODO: Implement configuration loading logic
	return &Config{
		Port: getEnv("PORT", "8080"),
		DB:   newDBConfig(),
	}, nil
}

func newDBConfig() *DBConfig {
	// TODO: Implement DB configuration loading logic
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "todo_app"),
	}
}

func getEnv(key string, defaultValue string) string {
	// TODO: Implement environment variable retrieval logic
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
