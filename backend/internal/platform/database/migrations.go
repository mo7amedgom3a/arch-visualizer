package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/config"
	"github.com/pressly/goose/v3"
)

// RunMigrations runs all pending migrations from the migrations directory
// It loads configuration from .env file in the backend root
func RunMigrations(migrationsDir string) error {
	// Load configuration from .env file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get database connection string from config
	dsn := cfg.Database.GetDSN()

	// Open database connection for goose
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection for migrations: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Resolve migrations directory path
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations directory path: %w", err)
	}

	// Check if migrations directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", absPath)
	}

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Run migrations
	if err := goose.Up(db, absPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// RollbackMigrations rolls back the last migration
// It loads configuration from .env file in the backend root
func RollbackMigrations(migrationsDir string) error {
	// Load configuration from .env file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get database connection string from config
	dsn := cfg.Database.GetDSN()

	// Open database connection for goose
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection for migrations: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Resolve migrations directory path
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations directory path: %w", err)
	}

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Rollback last migration
	if err := goose.Down(db, absPath); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
// It loads configuration from .env file in the backend root
func GetMigrationStatus(migrationsDir string) error {
	// Load configuration from .env file
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get database connection string from config
	dsn := cfg.Database.GetDSN()

	// Open database connection for goose
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection for migrations: %w", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Resolve migrations directory path
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations directory path: %w", err)
	}

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Get migration status (prints to stdout)
	if err := goose.Status(db, absPath); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
