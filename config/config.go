package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5/middleware"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	LogLevel    string
	Env         string
	CORS        middleware.CORSConfig
}

// LoadConfig loads configuration from environment variables
func LoadConfig(env string) *Config {
	var envFile string
	switch env {
	case "DEV":
		envFile = "config/.env.dev"
	case "PROD":
		envFile = "config/.env.prod"
	default:
		log.Printf("Warning: Invalid env '%s', defaulting to DEV", env)
		envFile = "config/.env.dev"
	}

	// Load .env file
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: Error loading %s file: %v", envFile, err)
	}

	config := &Config{
		ServerPort:  getEnv("SERVER_PORT", "1323"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Env:         getEnv("ENV", "development"),
		CORS: middleware.CORSConfig{
			AllowOrigins:     strings.Split(getEnv("ALLOW_ORIGINS", "*"), ","),
			AllowMethods:     strings.Split(getEnv("ALLOW_METHODS", "GET,POST,PUT,DELETE"), ","),
			AllowHeaders:     strings.Split(getEnv("ALLOW_HEADERS", ""), ","),
			AllowCredentials: getEnvAsBool("ALLOW_CREDENTIALS", false),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 0),
		},
	}

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as bool or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
