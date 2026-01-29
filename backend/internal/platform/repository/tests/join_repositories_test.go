package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

func TestJoinRepositories_TemplateAssociations(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(
		t,
		&models.Template{},
		&models.ComplianceStandard{},
		&models.IACFormat{},
		&models.Technology{},
		&models.TemplateCompliance{},
		&models.TemplateIACFormat{},
		&models.TemplateTechnology{},
		&models.Category{},
		&models.User{},
	)
	base := repository.NewBaseRepositoryWithDB(db)

	cat := &models.Category{ID: uuid.New(), Name: "Cat", Slug: "cat", CreatedAt: time.Now()}
	user := &models.User{ID: uuid.New(), Name: "Author", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := db.Create(cat).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	template := &models.Template{
		ID:               uuid.New(),
		Title:            "Assoc Template",
		Description:      "desc",
		CategoryID:       cat.ID,
		CloudProvider:    "AWS",
		EstimatedCostMin: 1,
		EstimatedCostMax: 2,
		AuthorID:         user.ID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := db.Create(template).Error; err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	cs := &models.ComplianceStandard{ID: uuid.New(), Name: "SOC2", Slug: "soc2", CreatedAt: time.Now()}
	format := &models.IACFormat{ID: uuid.New(), Name: "Terraform", Slug: "terraform", CreatedAt: time.Now()}
	tech := &models.Technology{ID: uuid.New(), Name: "Kubernetes", Slug: "kubernetes", CreatedAt: time.Now()}
	if err := db.Create(cs).Error; err != nil {
		t.Fatalf("failed to create compliance: %v", err)
	}
	if err := db.Create(format).Error; err != nil {
		t.Fatalf("failed to create format: %v", err)
	}
	if err := db.Create(tech).Error; err != nil {
		t.Fatalf("failed to create technology: %v", err)
	}

	tcRepo := &repository.TemplateComplianceRepository{BaseRepository: base}
	tfRepo := &repository.TemplateIACFormatRepository{BaseRepository: base}
	ttRepo := &repository.TemplateTechnologyRepository{BaseRepository: base}

	if err := tcRepo.Create(ctx, &models.TemplateCompliance{TemplateID: template.ID, ComplianceID: cs.ID}); err != nil {
		t.Fatalf("TemplateComplianceRepository.Create error: %v", err)
	}
	if err := tfRepo.Create(ctx, &models.TemplateIACFormat{TemplateID: template.ID, IACFormatID: format.ID}); err != nil {
		t.Fatalf("TemplateIACFormatRepository.Create error: %v", err)
	}
	if err := ttRepo.Create(ctx, &models.TemplateTechnology{TemplateID: template.ID, TechnologyID: tech.ID}); err != nil {
		t.Fatalf("TemplateTechnologyRepository.Create error: %v", err)
	}

	if _, err := tcRepo.FindByTemplate(ctx, template.ID); err != nil {
		t.Fatalf("TemplateComplianceRepository.FindByTemplate error: %v", err)
	}
	if _, err := tfRepo.FindByTemplate(ctx, template.ID); err != nil {
		t.Fatalf("TemplateIACFormatRepository.FindByTemplate error: %v", err)
	}
	if _, err := ttRepo.FindByTemplate(ctx, template.ID); err != nil {
		t.Fatalf("TemplateTechnologyRepository.FindByTemplate error: %v", err)
	}
}

func TestJoinRepositories_ResourceRelationships(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(
		t,
		&models.Project{},
		&models.Resource{},
		&models.ResourceContainment{},
		&models.ResourceDependency{},
		&models.User{},
		&models.IACTarget{},
		&models.DependencyType{},
	)
	base := repository.NewBaseRepositoryWithDB(db)

	user := &models.User{ID: uuid.New(), Name: "Owner", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	target := &models.IACTarget{ID: 1, Name: "terraform"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := db.Create(target).Error; err != nil {
		t.Fatalf("failed to create iac target: %v", err)
	}

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        user.ID,
		InfraToolID:   target.ID,
		Name:          "RelProj",
		CloudProvider: "aws",
		Region:        "us-east-1",
		CreatedAt:     time.Now(),
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	r1 := &models.Resource{ID: uuid.New(), ProjectID: project.ID, ResourceTypeID: 1, Name: "r1", PosX: 0, PosY: 0, CreatedAt: time.Now()}
	r2 := &models.Resource{ID: uuid.New(), ProjectID: project.ID, ResourceTypeID: 1, Name: "r2", PosX: 1, PosY: 1, CreatedAt: time.Now()}
	if err := db.Create(r1).Error; err != nil {
		t.Fatalf("failed to create resource1: %v", err)
	}
	if err := db.Create(r2).Error; err != nil {
		t.Fatalf("failed to create resource2: %v", err)
	}

	dt := &models.DependencyType{Name: "uses"}
	if err := db.Create(dt).Error; err != nil {
		t.Fatalf("failed to create dependency type: %v", err)
	}

	rcRepo := &repository.ResourceContainmentRepository{BaseRepository: base}
	rdRepo := &repository.ResourceDependencyRepository{BaseRepository: base}

	if err := rcRepo.Create(ctx, &models.ResourceContainment{ParentResourceID: r1.ID, ChildResourceID: r2.ID}); err != nil {
		t.Fatalf("ResourceContainmentRepository.Create error: %v", err)
	}
	if err := rdRepo.Create(ctx, &models.ResourceDependency{FromResourceID: r1.ID, ToResourceID: r2.ID, DependencyTypeID: dt.ID}); err != nil {
		t.Fatalf("ResourceDependencyRepository.Create error: %v", err)
	}

	if _, err := rcRepo.FindChildren(ctx, r1.ID); err != nil {
		t.Fatalf("ResourceContainmentRepository.FindChildren error: %v", err)
	}
	if _, err := rcRepo.FindParents(ctx, r2.ID); err != nil {
		t.Fatalf("ResourceContainmentRepository.FindParents error: %v", err)
	}
	if _, err := rdRepo.FindByFromResource(ctx, r1.ID); err != nil {
		t.Fatalf("ResourceDependencyRepository.FindByFromResource error: %v", err)
	}
	if _, err := rdRepo.FindByToResource(ctx, r2.ID); err != nil {
		t.Fatalf("ResourceDependencyRepository.FindByToResource error: %v", err)
	}
}

