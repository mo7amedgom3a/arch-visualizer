package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	templaterepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/template"
)

func TestMarketplaceRepositories_BasicOperations(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(
		t,
		&models.Category{},
		&models.Template{},
		&models.Review{},
		&models.IACFormat{},
		&models.Technology{},
		&models.ComplianceStandard{},
		&models.User{},
	)
	base := repository.NewBaseRepositoryWithDB(db)

	// CategoryRepository
	catRepo := &templaterepo.CategoryRepository{BaseRepository: base}
	category := &models.Category{
		ID:        uuid.New(),
		Name:      "Security",
		Slug:      "security",
		CreatedAt: time.Now(),
	}
	if err := catRepo.Create(ctx, category); err != nil {
		t.Fatalf("CategoryRepository.Create error: %v", err)
	}
	if _, err := catRepo.FindBySlug(ctx, "security"); err != nil {
		t.Fatalf("CategoryRepository.FindBySlug error: %v", err)
	}

	// IACFormatRepository
	iacRepo := &templaterepo.IACFormatRepository{BaseRepository: base}
	format := &models.IACFormat{
		ID:        uuid.New(),
		Name:      "Terraform",
		Slug:      "terraform",
		CreatedAt: time.Now(),
	}
	if err := iacRepo.Create(ctx, format); err != nil {
		t.Fatalf("IACFormatRepository.Create error: %v", err)
	}
	if _, err := iacRepo.FindBySlug(ctx, "terraform"); err != nil {
		t.Fatalf("IACFormatRepository.FindBySlug error: %v", err)
	}

	// TechnologyRepository
	techRepo := &templaterepo.TechnologyRepository{BaseRepository: base}
	tech := &models.Technology{
		ID:        uuid.New(),
		Name:      "Kubernetes",
		Slug:      "kubernetes",
		CreatedAt: time.Now(),
	}
	if err := techRepo.Create(ctx, tech); err != nil {
		t.Fatalf("TechnologyRepository.Create error: %v", err)
	}
	if _, err := techRepo.FindBySlug(ctx, "kubernetes"); err != nil {
		t.Fatalf("TechnologyRepository.FindBySlug error: %v", err)
	}

	// ComplianceStandardRepository
	csRepo := &templaterepo.ComplianceStandardRepository{BaseRepository: base}
	cs := &models.ComplianceStandard{
		ID:        uuid.New(),
		Name:      "SOC2",
		Slug:      "soc2",
		CreatedAt: time.Now(),
	}
	if err := csRepo.Create(ctx, cs); err != nil {
		t.Fatalf("ComplianceStandardRepository.Create error: %v", err)
	}
	if _, err := csRepo.FindBySlug(ctx, "soc2"); err != nil {
		t.Fatalf("ComplianceStandardRepository.FindBySlug error: %v", err)
	}

	// TemplateRepository
	tmplRepo := &templaterepo.TemplateRepository{BaseRepository: base}
	author := &models.User{
		ID:        uuid.New(),
		Name:      "Author",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(author).Error; err != nil {
		t.Fatalf("failed to create author: %v", err)
	}
	template := &models.Template{
		ID:               uuid.New(),
		Title:            "VPC Baseline",
		Description:      "Base networking",
		CategoryID:       category.ID,
		CloudProvider:    "AWS",
		Rating:           0,
		ReviewCount:      0,
		Downloads:        0,
		Price:            0,
		IsSubscription:   false,
		EstimatedCostMin: 1,
		EstimatedCostMax: 10,
		AuthorID:         author.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := tmplRepo.Create(ctx, template); err != nil {
		t.Fatalf("TemplateRepository.Create error: %v", err)
	}
	_, err := tmplRepo.FindByID(ctx, template.ID)
	if err != nil {
		t.Fatalf("TemplateRepository.FindByID error: %v", err)
	}

	// ReviewRepository
	reviewRepo := &templaterepo.ReviewRepository{BaseRepository: base}
	review := &models.Review{
		ID:         uuid.New(),
		TemplateID: template.ID,
		UserID:     author.ID,
		Rating:     5,
		Title:      "Great",
		Content:    "Works well",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := reviewRepo.Create(ctx, review); err != nil {
		t.Fatalf("ReviewRepository.Create error: %v", err)
	}
	if _, err := reviewRepo.FindByTemplate(ctx, template.ID, 10, 0); err != nil {
		t.Fatalf("ReviewRepository.FindByTemplate error: %v", err)
	}
}
