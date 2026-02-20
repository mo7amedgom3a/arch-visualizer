package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ProjectService handles project management and persistence.
// Mutation semantics:
//   - Project metadata (Update) is in-place – no snapshot is created.
//   - All architecture / state changes happen through CreateVersion, which always
//     clones the project into a new immutable snapshot.
type ProjectService interface {
	// ── Project CRUD (non-versioned) ─────────────────────────────────────────

	// Create creates a new root project row (no architecture yet).
	Create(ctx context.Context, req *CreateProjectRequest) (*models.Project, error)

	// GetByID retrieves a project snapshot by its exact ID.
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)

	// List retrieves projects with pagination and filtering.
	List(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error)

	// UpdateMetadata performs an in-place update of project metadata fields
	// (name, description, tags, cloud_provider, region, thumbnail, iac_tool_id).
	// No new snapshot is created.
	UpdateMetadata(ctx context.Context, project *models.Project) (*models.Project, error)

	// Duplicate creates an independent copy of the project (new root, version 1).
	Duplicate(ctx context.Context, projectID uuid.UUID, name string) (*models.Project, *models.ProjectVersion, error)

	// Delete hard-deletes a project snapshot and its resources.
	Delete(ctx context.Context, id uuid.UUID) error

	// ── Version CRUD ─────────────────────────────────────────────────────────

	// CreateVersion snapshots the supplied architecture as a new immutable version.
	// It clones the current project row and returns the new version metadata.
	CreateVersion(ctx context.Context, projectID uuid.UUID, req *CreateVersionRequest) (*ProjectVersionDetail, error)

	// GetVersions returns the full ordered version chain for any project ID in the lineage.
	GetVersions(ctx context.Context, projectID uuid.UUID) ([]*ProjectVersionSummary, error)

	// GetLatestVersion returns the latest version with full architecture state.
	GetLatestVersion(ctx context.Context, projectID uuid.UUID) (*ProjectVersionDetail, error)

	// GetVersionByID returns a specific version with full architecture state.
	GetVersionByID(ctx context.Context, projectID uuid.UUID, versionID uuid.UUID) (*ProjectVersionDetail, error)

	// DeleteVersion removes a single version entry (does not delete the snapshot project row).
	DeleteVersion(ctx context.Context, projectID uuid.UUID, versionID uuid.UUID) error

	// ── Architecture (read-only) ──────────────────────────────────────────────

	// GetArchitecture retrieves the full architecture for a project snapshot.
	GetArchitecture(ctx context.Context, projectID uuid.UUID) (*dto.ArchitectureResponse, error)

	// ── Utility actions (version-scoped) ─────────────────────────────────────

	// ValidateVersionArchitecture validates the architecture stored in a specific version.
	ValidateVersionArchitecture(ctx context.Context, versionID uuid.UUID) (*dto.ValidationResponse, error)

	// ── Orchestrator compatibility ────────────────────────────────────────────

	// PersistArchitecture persists an architecture as part of a project (used by pipeline orchestrator).
	PersistArchitecture(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}) error

	// PersistArchitectureWithPricing persists an architecture with pricing calculation.
	PersistArchitectureWithPricing(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}, pricingDuration time.Duration) (*ArchitecturePersistResult, error)

	// LoadArchitecture loads an architecture from the database for a project.
	LoadArchitecture(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error)

	// GetProjectPricing retrieves pricing for a project.
	GetProjectPricing(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error)
}

// ── Request / Response types ─────────────────────────────────────────────────

// CreateProjectRequest contains data needed to create a project.
type CreateProjectRequest struct {
	UserID        uuid.UUID
	Name          string
	Description   string
	Tags          []string
	IACTargetID   uint
	CloudProvider string
	Region        string
}

// CreateVersionRequest is the payload for POST /projects/{id}/versions.
type CreateVersionRequest struct {
	Nodes     []dto.ArchitectureNode     `json:"nodes"`
	Edges     []dto.ArchitectureEdge     `json:"edges"`
	Variables []dto.ArchitectureVariable `json:"variables"`
	Outputs   []dto.ArchitectureOutput   `json:"outputs"`
	Message   string                     `json:"message"`
}

// ProjectVersionSummary is a lightweight listing entry for GET /projects/{id}/versions.
type ProjectVersionSummary struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       uuid.UUID  `json:"project_id"`
	ParentVersionID *uuid.UUID `json:"parent_version_id"`
	VersionNumber   int        `json:"version_number"`
	Message         string     `json:"message,omitempty"`
	CreatedAt       string     `json:"created_at"`
	CreatedBy       uuid.UUID  `json:"created_by"`
}

// ProjectVersionDetail is the full version response including architecture state.
type ProjectVersionDetail struct {
	ProjectVersionSummary
	State *dto.ArchitectureResponse `json:"state"`
}

// ArchitecturePersistResult contains the result of persisting an architecture with pricing.
type ArchitecturePersistResult struct {
	ResourceIDMapping map[string]uuid.UUID      `json:"resource_id_mapping"`
	PricingEstimate   *ArchitectureCostEstimate `json:"pricing_estimate,omitempty"`
}

// ── Kept for backward compatibility with orchestrator types ──────────────────

// VersionedOperationResult is kept because the orchestrator test mock references it.
// New code should use ProjectVersionDetail / ProjectVersionSummary instead.
type VersionedOperationResult struct {
	NewProjectID  uuid.UUID `json:"project_id"`
	VersionID     uuid.UUID `json:"version_id"`
	VersionNumber int       `json:"version_number"`
}

// VersionedArchitectureResult is kept for scenario20 runner compat.
type VersionedArchitectureResult struct {
	VersionedOperationResult
	Architecture *dto.ArchitectureResponse `json:"architecture"`
}

// VersionedNodeResult is kept for pipeline_test.go mock compat.
type VersionedNodeResult struct {
	VersionedOperationResult
	Node *dto.ArchitectureNode `json:"node"`
}
