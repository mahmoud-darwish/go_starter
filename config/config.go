package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

var config *Config

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		// Ignore if .env file is missing
	}

	config = &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://mahmoudibrahim:kokowawa@localhost:5432/shouf?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "my_ultra_secure_jwt_secret_key"),
		Port:        getEnv("PORT", "8080"),
	}
	return nil
}

func GetConfig() *Config {
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
