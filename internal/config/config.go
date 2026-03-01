package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      App
	Database Database
	JWT      JWT
	MinIO    MinIO
}

type MinIO struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type App struct {
	APIHost string
	APIPort string
}

type Database struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

type JWT struct {
	JWTSecret            string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration
}

func Load() *Config {
	return &Config{
		App: App{
			APIHost: getEnv("API_HOST", "0.0.0.0"),
			APIPort: getEnv("API_PORT", "8081"),
		},

		Database: Database{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("POSTGRES_PORT", "5432"),
			DBUser:     getEnv("POSTGRES_USER", "transaction"),
			DBPassword: getEnv("POSTGRES_PASSWORD", "password"),
			DBName:     getEnv("POSTGRES_DB", "transaction_db"),
			DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWT{
			JWTSecret:            getEnv("JWT_SECRET", "secret"),
			JWTAccessExpiration:  getDuration("JWT_ACCESS_EXPIRATION", 15*time.Minute),
			JWTRefreshExpiration: getDuration("JWT_REFRESH_EXPIRATION", 24*time.Hour),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}
