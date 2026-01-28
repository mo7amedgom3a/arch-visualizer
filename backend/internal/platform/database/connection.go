package database

import (
	"fmt"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection instance
var DB *gorm.DB

// Connect initializes and returns a GORM database connection
// It loads configuration from .env file in the backend root
func Connect() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	// Load configuration from .env file
	cfg, err := config.Load()
	if err != nil {
		return nil, platformerrors.NewDatabaseConfigError("failed to load configuration").WithMeta("cause", err.Error())
	}

	// Get database connection string from config
	dsn := cfg.Database.GetDSN()

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err).WithMeta("operation", "get_sql_db")
	}
	fmt.Println("Connected to database")

	if err := sqlDB.Ping(); err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err).WithMeta("operation", "ping")
	}

	DB = db
	return DB, nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
