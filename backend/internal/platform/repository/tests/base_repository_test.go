package repository_test

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

func TestBaseRepository_CommitAndRollback(t *testing.T) {
	db := newTestDB(t, &models.User{})
	base := repository.NewBaseRepositoryWithDB(db)

	ctx := context.Background()

	// Commit path
	{
		tx, txCtx := base.BeginTransaction(ctx)
		user := &models.User{Name: "Commit User"}
		if err := base.GetDB(txCtx).Create(user).Error; err != nil {
			t.Fatalf("failed to create user in transaction: %v", err)
		}
		if err := base.CommitTransaction(tx); err != nil {
			t.Fatalf("CommitTransaction returned error: %v", err)
		}
	}

	// Rollback path
	{
		tx, txCtx := base.BeginTransaction(ctx)
		user := &models.User{Name: "Rollback User"}
		if err := base.GetDB(txCtx).Create(user).Error; err != nil {
			t.Fatalf("failed to create user in transaction: %v", err)
		}
		if err := base.RollbackTransaction(tx); err != nil {
			t.Fatalf("RollbackTransaction returned error: %v", err)
		}
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 committed user and 0 rolled-back users, got %d total", count)
	}
}

