package repository

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository using the default database
// connection configured via the platform database package.
func NewBaseRepository() (*BaseRepository, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &BaseRepository{db: db}, nil
}

// NewBaseRepositoryWithDB creates a new base repository with an injected *gorm.DB
// instance. This is primarily intended for tests or alternative storage backends
// (e.g., in-memory SQLite) where you want to control the underlying connection.
func NewBaseRepositoryWithDB(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the GORM database instance
func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if ctx != nil {
		if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
			return tx
		}
	}
	return r.db
}

// BeginTransaction starts a new database transaction
func (r *BaseRepository) BeginTransaction(ctx context.Context) (*gorm.DB, context.Context) {
	tx := r.db.Begin()
	return tx, context.WithValue(ctx, "tx", tx)
}

// CommitTransaction commits a transaction
func (r *BaseRepository) CommitTransaction(tx *gorm.DB) error {
	err := tx.Commit().Error
	if err != nil {
		return platformerrors.NewDatabaseTransactionFailed("BaseRepository.CommitTransaction", err)
	}
	return nil
}

// RollbackTransaction rolls back a transaction
func (r *BaseRepository) RollbackTransaction(tx *gorm.DB) error {
	err := tx.Rollback().Error
	if err != nil {
		return platformerrors.NewDatabaseTransactionFailed("BaseRepository.RollbackTransaction", err)
	}
	return nil
}
