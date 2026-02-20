package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	infrastructurerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/infrastructure"
	pricingrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/pricing"
	projectrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/project"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
	"gorm.io/gorm"
)

// Repository adapters that wrap concrete repositories to match interfaces

// ProjectRepositoryAdapter adapts repository.ProjectRepository to serverinterfaces.ProjectRepository
type ProjectRepositoryAdapter struct {
	Repo *projectrepo.ProjectRepository
}

func (a *ProjectRepositoryAdapter) Create(ctx context.Context, project *models.Project) error {
	return a.Repo.Create(ctx, project)
}

func (a *ProjectRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return a.Repo.FindByID(ctx, id)
}

func (a *ProjectRepositoryAdapter) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	return a.Repo.FindByUserID(ctx, userID)
}

func (a *ProjectRepositoryAdapter) Update(ctx context.Context, project *models.Project) error {
	return a.Repo.Update(ctx, project)
}

func (a *ProjectRepositoryAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	return a.Repo.Delete(ctx, id)
}

func (a *ProjectRepositoryAdapter) BeginTransaction(ctx context.Context) (interface{}, context.Context) {
	return a.Repo.BaseRepository.BeginTransaction(ctx)
}

func (a *ProjectRepositoryAdapter) CommitTransaction(tx interface{}) error {
	if gormTx, ok := tx.(*gorm.DB); ok {
		return a.Repo.BaseRepository.CommitTransaction(gormTx)
	}
	return nil
}

func (a *ProjectRepositoryAdapter) RollbackTransaction(tx interface{}) error {
	if gormTx, ok := tx.(*gorm.DB); ok {
		return a.Repo.BaseRepository.RollbackTransaction(gormTx)
	}
	return nil
}

func (a *ProjectRepositoryAdapter) FindAll(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error) {
	return a.Repo.FindAll(ctx, userID, page, limit, sort, order, search)
}

func (a *ProjectRepositoryAdapter) FindByRootProjectID(ctx context.Context, rootProjectID uuid.UUID) ([]*models.Project, error) {
	return a.Repo.FindByRootProjectID(ctx, rootProjectID)
}

// ProjectVersionRepositoryAdapter adapts repository.ProjectVersionRepository to serverinterfaces.ProjectVersionRepository
type ProjectVersionRepositoryAdapter struct {
	Repo *projectrepo.ProjectVersionRepository
}

func (a *ProjectVersionRepositoryAdapter) Create(ctx context.Context, version *models.ProjectVersion) error {
	return a.Repo.Create(ctx, version)
}

func (a *ProjectVersionRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*models.ProjectVersion, error) {
	return a.Repo.FindByID(ctx, id)
}

func (a *ProjectVersionRepositoryAdapter) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVersion, error) {
	return a.Repo.ListByProjectID(ctx, projectID)
}

func (a *ProjectVersionRepositoryAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	return a.Repo.Delete(ctx, id)
}

func (a *ProjectVersionRepositoryAdapter) ListByRootProjectID(ctx context.Context, rootProjectID uuid.UUID) ([]*models.ProjectVersion, error) {
	return a.Repo.ListByRootProjectID(ctx, rootProjectID)
}

func (a *ProjectVersionRepositoryAdapter) GetLatestVersionForProject(ctx context.Context, projectID uuid.UUID) (*models.ProjectVersion, error) {
	return a.Repo.GetLatestVersionForProject(ctx, projectID)
}

// ResourceRepositoryAdapter adapts repository.ResourceRepository to serverinterfaces.ResourceRepository
type ResourceRepositoryAdapter struct {
	Repo *resourcerepo.ResourceRepository
}

func (a *ResourceRepositoryAdapter) Create(ctx context.Context, resource *models.Resource) error {
	return a.Repo.Create(ctx, resource)
}

func (a *ResourceRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*models.Resource, error) {
	return a.Repo.FindByID(ctx, id)
}

func (a *ResourceRepositoryAdapter) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Resource, error) {
	return a.Repo.FindByProjectID(ctx, projectID)
}

func (a *ResourceRepositoryAdapter) CreateContainment(ctx context.Context, parentID, childID uuid.UUID) error {
	return a.Repo.CreateContainment(ctx, parentID, childID)
}

func (a *ResourceRepositoryAdapter) CreateDependency(ctx context.Context, dependency *models.ResourceDependency) error {
	return a.Repo.CreateDependency(ctx, dependency)
}

func (a *ResourceRepositoryAdapter) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return a.Repo.DeleteByProjectID(ctx, projectID)
}

// ResourceTypeRepositoryAdapter adapts repository.ResourceTypeRepository to serverinterfaces.ResourceTypeRepository
type ResourceTypeRepositoryAdapter struct {
	Repo *resourcerepo.ResourceTypeRepository
}

func (a *ResourceTypeRepositoryAdapter) FindByNameAndProvider(ctx context.Context, name, provider string) (*models.ResourceType, error) {
	return a.Repo.FindByNameAndProvider(ctx, name, provider)
}

func (a *ResourceTypeRepositoryAdapter) ListByProvider(ctx context.Context, provider string) ([]*models.ResourceType, error) {
	return a.Repo.ListByProvider(ctx, provider)
}

// ResourceConstraintRepositoryAdapter adapts repository.ResourceConstraintRepository to serverinterfaces.ResourceConstraintRepository
type ResourceConstraintRepositoryAdapter struct {
	Repo *resourcerepo.ResourceConstraintRepository
}

func (a *ResourceConstraintRepositoryAdapter) FindByResourceType(ctx context.Context, resourceTypeID uint) ([]*models.ResourceConstraint, error) {
	return a.Repo.FindByResourceType(ctx, resourceTypeID)
}

// DependencyTypeRepositoryAdapter adapts repository.DependencyTypeRepository to serverinterfaces.DependencyTypeRepository
type DependencyTypeRepositoryAdapter struct {
	Repo *resourcerepo.DependencyTypeRepository
}

func (a *DependencyTypeRepositoryAdapter) FindByName(ctx context.Context, name string) (*models.DependencyType, error) {
	return a.Repo.FindByName(ctx, name)
}

func (a *DependencyTypeRepositoryAdapter) Create(ctx context.Context, depType *models.DependencyType) error {
	// The concrete repository doesn't have Create, so we'll need to add it or use direct DB access
	// For now, return an error indicating this needs to be implemented
	return nil // TODO: Implement if needed
}

// UserRepositoryAdapter adapts repository.UserRepository to serverinterfaces.UserRepository
type UserRepositoryAdapter struct {
	Repo *userrepo.UserRepository
}

func (a *UserRepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return a.Repo.FindByID(ctx, id)
}

func (a *UserRepositoryAdapter) Create(ctx context.Context, user *models.User) error {
	return a.Repo.Create(ctx, user)
}

// IACTargetRepositoryAdapter adapts repository.IACTargetRepository to serverinterfaces.IACTargetRepository
type IACTargetRepositoryAdapter struct {
	Repo *infrastructurerepo.IACTargetRepository
}

func (a *IACTargetRepositoryAdapter) FindByName(ctx context.Context, name string) (*models.IACTarget, error) {
	return a.Repo.FindByName(ctx, name)
}

func (a *IACTargetRepositoryAdapter) Create(ctx context.Context, target *models.IACTarget) error {
	return a.Repo.Create(ctx, target)
}

// ResourceContainmentRepositoryAdapter adapts repository.ResourceContainmentRepository to serverinterfaces.ResourceContainmentRepository
type ResourceContainmentRepositoryAdapter struct {
	Repo *resourcerepo.ResourceContainmentRepository
}

func (a *ResourceContainmentRepositoryAdapter) Create(ctx context.Context, containment *models.ResourceContainment) error {
	return a.Repo.Create(ctx, containment)
}

func (a *ResourceContainmentRepositoryAdapter) FindChildren(ctx context.Context, parentID uuid.UUID) ([]*models.ResourceContainment, error) {
	return a.Repo.FindChildren(ctx, parentID)
}

func (a *ResourceContainmentRepositoryAdapter) FindParents(ctx context.Context, childID uuid.UUID) ([]*models.ResourceContainment, error) {
	return a.Repo.FindParents(ctx, childID)
}

// ResourceDependencyRepositoryAdapter adapts repository.ResourceDependencyRepository to serverinterfaces.ResourceDependencyRepository
type ResourceDependencyRepositoryAdapter struct {
	Repo *resourcerepo.ResourceDependencyRepository
}

func (a *ResourceDependencyRepositoryAdapter) Create(ctx context.Context, dependency *models.ResourceDependency) error {
	return a.Repo.Create(ctx, dependency)
}

func (a *ResourceDependencyRepositoryAdapter) FindByFromResource(ctx context.Context, fromID uuid.UUID) ([]*models.ResourceDependency, error) {
	return a.Repo.FindByFromResource(ctx, fromID)
}

func (a *ResourceDependencyRepositoryAdapter) FindByToResource(ctx context.Context, toID uuid.UUID) ([]*models.ResourceDependency, error) {
	return a.Repo.FindByToResource(ctx, toID)
}

// PricingRepositoryAdapter adapts repository.PricingRepository to serverinterfaces.PricingRepository
type PricingRepositoryAdapter struct {
	Repo *pricingrepo.PricingRepository
}

func (a *PricingRepositoryAdapter) CreateProjectPricing(ctx context.Context, pricing *models.ProjectPricing) error {
	return a.Repo.CreateProjectPricing(ctx, pricing)
}

func (a *PricingRepositoryAdapter) FindProjectPricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	return a.Repo.FindProjectPricingByProjectID(ctx, projectID)
}

func (a *PricingRepositoryAdapter) CreateResourcePricing(ctx context.Context, pricing *models.ResourcePricing) error {
	return a.Repo.CreateResourcePricing(ctx, pricing)
}

func (a *PricingRepositoryAdapter) FindResourcePricingByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error) {
	return a.Repo.FindResourcePricingByResourceID(ctx, resourceID)
}

func (a *PricingRepositoryAdapter) FindResourcePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ResourcePricing, error) {
	return a.Repo.FindResourcePricingByProjectID(ctx, projectID)
}

func (a *PricingRepositoryAdapter) CreatePricingComponent(ctx context.Context, component *models.PricingComponent) error {
	return a.Repo.CreatePricingComponent(ctx, component)
}

// ProjectVariableRepositoryAdapter adapts repository.ProjectVariableRepository to serverinterfaces.ProjectVariableRepository
type ProjectVariableRepositoryAdapter struct {
	Repo *projectrepo.ProjectVariableRepository
}

func (a *ProjectVariableRepositoryAdapter) Create(ctx context.Context, variable *models.ProjectVariable) error {
	return a.Repo.Create(ctx, variable)
}

func (a *ProjectVariableRepositoryAdapter) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVariable, error) {
	return a.Repo.FindByProjectID(ctx, projectID)
}

func (a *ProjectVariableRepositoryAdapter) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return a.Repo.DeleteByProjectID(ctx, projectID)
}

// ProjectOutputRepositoryAdapter adapts repository.ProjectOutputRepository to serverinterfaces.ProjectOutputRepository
type ProjectOutputRepositoryAdapter struct {
	Repo *projectrepo.ProjectOutputRepository
}

func (a *ProjectOutputRepositoryAdapter) Create(ctx context.Context, output *models.ProjectOutput) error {
	return a.Repo.Create(ctx, output)
}

func (a *ProjectOutputRepositoryAdapter) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectOutput, error) {
	return a.Repo.FindByProjectID(ctx, projectID)
}

func (a *ProjectOutputRepositoryAdapter) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return a.Repo.DeleteByProjectID(ctx, projectID)
}
