package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	SupabaseURL    string
	SupabaseJWTKey string
	AllowedOrigins []string
	Environment    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	requiredEnvVars := []string{
		"PORT",
		"SUPABASE_URL",
		"SUPABASE_JWT_KEY",
	}

	missingEnvVars := []string{}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			missingEnvVars = append(missingEnvVars, envVar)
		}
	}

	if len(missingEnvVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingEnvVars, ","))
	}

	// Set default values for optional env variables
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	// Parse allowed origins
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
		// Default to localhost if not specified
		allowedOrigins = []string{"http://localhost:3000"}
	}

	// Trim spaces
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		SupabaseURL:    os.Getenv("SUPABASE_URL"),
		SupabaseJWTKey: os.Getenv("SUPABASE_JWT_KEY"),
		AllowedOrigins: allowedOrigins,
		Environment:    environment,
	}, nil
}
