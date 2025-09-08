package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration settings.
type Config struct {
	DatabaseURL  string
	DatabaseName string
	ServerPort   string
	JWTSecret    string
}

// Load loads configuration from environment variables or .env file.
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading from environment variables.")
	}

	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "portaria_keys"),
		ServerPort:   getEnv("SERVER_PORT", ":8080"),
		JWTSecret:    getEnv("JWT_SECRET", "supersecretkey"), // TODO: Change for production
	}
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}