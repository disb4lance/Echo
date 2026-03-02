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
	UseSSL     bool
	Location   string
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
			DBPort:     getEnv("DB_PORT", "5434"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "postgres"),
			DBName:     getEnv("DB_NAME", "mydatabase"),
			DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWT{
			JWTSecret:            getEnv("JWT_SECRET", "xITmnGcNLhpe/R@qfuA>R1lfV/KTkF1v9J8++:mfS*l"),
			JWTAccessExpiration:  getDuration("JWT_ACCESS_EXPIRATION", 15*time.Minute),
			JWTRefreshExpiration: getDuration("JWT_REFRESH_EXPIRATION", 24*time.Hour*7),
		},
		MinIO: MinIO{
			Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
			BucketName: getEnv("MINIO_BUCKET_NAME", "uploads"),
			UseSSL:     getEnvAsBool("MINIO_USE_SSL", false),
			Location:   getEnv("MINIO_LOCATION", "us-east-1"),
		},
	}
}

// Вспомогательная функция для bool значений
func getEnvAsBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return boolVal
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
