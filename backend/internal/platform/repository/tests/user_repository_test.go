package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := newTestDB(t, &models.User{})
	base := repository.NewBaseRepositoryWithDB(db)
	repo := &userrepo.UserRepository{BaseRepository: base}

	ctx := context.Background()

	user := &models.User{
		ID:   uuid.New(),
		Name: "Alice",
	}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	foundByID, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if foundByID.Name != user.Name {
		t.Fatalf("expected name %q, got %q", user.Name, foundByID.Name)
	}
	fmt.Println("foundByID", foundByID.ID)
	fmt.Println("user", user.Name)

	list, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 user in list, got %d", len(list))
	}
	fmt.Println("list", list)
}
