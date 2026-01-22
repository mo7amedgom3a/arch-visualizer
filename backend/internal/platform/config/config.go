package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load loads configuration from .env file in the backend root directory
func Load() (*Config, error) {
	// Get the backend root directory (assuming this is called from backend/)
	// Try to find .env file in the backend root
	backendRoot := findBackendRoot()
	envPath := filepath.Join(backendRoot, ".env")

	// Load .env file if it exists (don't fail if it doesn't exist)
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "arch_visualizer"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// findBackendRoot finds the backend root directory by looking for go.mod
func findBackendRoot() string {
	// Start from current working directory and walk up
	dir, err := os.Getwd()
	if err != nil {
		// Fallback to relative path
		return "."
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, return current directory
			return "."
		}
		dir = parent
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
