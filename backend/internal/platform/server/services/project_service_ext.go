package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ── Project CRUD ──────────────────────────────────────────────────────────────

// Duplicate creates an independent copy of the project (new root, version 1).
func (s *ProjectServiceImpl) Duplicate(ctx context.Context, projectID uuid.UUID, name string) (*models.Project, *models.ProjectVersion, error) {
	originalProject, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find original project: %w", err)
	}

	arch, err := s.LoadArchitecture(ctx, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load architecture: %w", err)
	}

	newProject := &models.Project{
		ID:            uuid.New(),
		RootProjectID: nil,
		UserID:        originalProject.UserID,
		InfraToolID:   originalProject.InfraToolID,
		Name:          name,
		Description:   originalProject.Description,
		Tags:          originalProject.Tags,
		CloudProvider: originalProject.CloudProvider,
		Region:        originalProject.Region,
		Thumbnail:     originalProject.Thumbnail,
	}
	if err := s.projectRepo.Create(ctx, newProject); err != nil {
		return nil, nil, fmt.Errorf("failed to create duplicated project: %w", err)
	}

	if err := s.PersistArchitecture(ctx, newProject.ID, arch, nil); err != nil {
		_ = s.projectRepo.Delete(ctx, newProject.ID)
		return nil, nil, fmt.Errorf("failed to persist duplicated architecture: %w", err)
	}

	version := &models.ProjectVersion{
		ID:            uuid.New(),
		ProjectID:     newProject.ID,
		VersionNumber: 1,
		Message:       "Duplicated from " + originalProject.Name,
		CreatedBy:     newProject.UserID,
	}
	_ = s.versionRepo.Create(ctx, version)

	return newProject, version, nil
}

// UpdateMetadata performs an in-place update of project metadata (no snapshot created).
func (s *ProjectServiceImpl) UpdateMetadata(ctx context.Context, project *models.Project) (*models.Project, error) {
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("UpdateMetadata: %w", err)
	}
	updated, err := s.projectRepo.FindByID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("UpdateMetadata: reload: %w", err)
	}
	return updated, nil
}

// ── Version CRUD ────────────────────────────────────────────────────────────

// CreateVersion snapshots the supplied architecture as a new immutable version.
func (s *ProjectServiceImpl) CreateVersion(ctx context.Context, projectID uuid.UUID, req *serverinterfaces.CreateVersionRequest) (*serverinterfaces.ProjectVersionDetail, error) {
	archReq := &dto.UpdateArchitectureRequest{
		Nodes:     req.Nodes,
		Edges:     req.Edges,
		Variables: req.Variables,
		Outputs:   req.Outputs,
	}

	result, newProject, err := s.cloneProjectSnapshot(ctx, cloneProjectSnapshotOptions{
		sourceProjectID: projectID,
		applyArch:       archReq,
		message:         req.Message,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateVersion: %w", err)
	}

	arch, err := s.GetArchitecture(ctx, newProject.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateVersion: load architecture: %w", err)
	}

	// Fetch the version record we just created
	ver, err := s.versionRepo.FindByID(ctx, result.VersionID)
	if err != nil {
		return nil, fmt.Errorf("CreateVersion: fetch version: %w", err)
	}

	return versionDetail(ver, arch), nil
}

// GetVersions returns the full ordered version chain for any project ID in the lineage.
func (s *ProjectServiceImpl) GetVersions(ctx context.Context, projectID uuid.UUID) ([]*serverinterfaces.ProjectVersionSummary, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	rootID := projectID
	if project.RootProjectID != nil {
		rootID = *project.RootProjectID
	}

	versions, err := s.versionRepo.ListByRootProjectID(ctx, rootID)
	if err != nil {
		return nil, err
	}

	out := make([]*serverinterfaces.ProjectVersionSummary, len(versions))
	for i, v := range versions {
		out[i] = versionSummary(v)
	}
	return out, nil
}

// GetLatestVersion returns the most recent version with full architecture state.
func (s *ProjectServiceImpl) GetLatestVersion(ctx context.Context, projectID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	rootID := projectID
	if project.RootProjectID != nil {
		rootID = *project.RootProjectID
	}

	versions, err := s.versionRepo.ListByRootProjectID(ctx, rootID)
	if err != nil || len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for project %s", projectID)
	}

	latest := versions[len(versions)-1]
	arch, err := s.GetArchitecture(ctx, latest.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("GetLatestVersion: load architecture: %w", err)
	}
	return versionDetail(latest, arch), nil
}

// GetVersionByID returns a specific version with full architecture state.
func (s *ProjectServiceImpl) GetVersionByID(ctx context.Context, projectID uuid.UUID, versionID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error) {
	ver, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	arch, err := s.GetArchitecture(ctx, ver.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("GetVersionByID: load architecture: %w", err)
	}
	return versionDetail(ver, arch), nil
}

// DeleteVersion removes a single version entry (does NOT delete the project snapshot).
func (s *ProjectServiceImpl) DeleteVersion(ctx context.Context, projectID uuid.UUID, versionID uuid.UUID) error {
	// Verify the version belongs to this project lineage
	ver, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	project, _ := s.projectRepo.FindByID(ctx, projectID)
	if project == nil {
		return fmt.Errorf("project not found")
	}

	// Verify version is in the same lineage
	rootID := projectID
	if project.RootProjectID != nil {
		rootID = *project.RootProjectID
	}

	verProject, err := s.projectRepo.FindByID(ctx, ver.ProjectID)
	if err != nil {
		return fmt.Errorf("version project snapshot not found: %w", err)
	}

	verRootID := ver.ProjectID
	if verProject.RootProjectID != nil {
		verRootID = *verProject.RootProjectID
	}

	if verRootID != rootID {
		return fmt.Errorf("version %s does not belong to project lineage %s", versionID, projectID)
	}

	return s.versionRepo.Delete(ctx, versionID)
}

// ValidateVersionArchitecture validates the architecture stored in a specific version.
func (s *ProjectServiceImpl) ValidateVersionArchitecture(ctx context.Context, versionID uuid.UUID) (*dto.ValidationResponse, error) {
	ver, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}
	arch, err := s.LoadArchitecture(ctx, ver.ProjectID)
	if err != nil {
		return nil, err
	}

	errs := make([]dto.ValidationIssue, 0)
	valid := true
	if len(arch.Resources) == 0 {
		valid = false
		errs = append(errs, dto.ValidationIssue{
			Type:     "structural",
			Message:  "Architecture is empty",
			Severity: "warning",
		})
	}

	return &dto.ValidationResponse{
		Valid:    valid,
		Errors:   errs,
		Warnings: []dto.ValidationIssue{},
	}, nil
}

// ── internal helpers ──────────────────────────────────────────────────────────

func versionSummary(v *models.ProjectVersion) *serverinterfaces.ProjectVersionSummary {
	return &serverinterfaces.ProjectVersionSummary{
		ID:              v.ID,
		ProjectID:       v.ProjectID,
		ParentVersionID: v.ParentVersionID,
		VersionNumber:   v.VersionNumber,
		Message:         v.Message,
		CreatedAt:       v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:       v.CreatedBy,
	}
}

func versionDetail(v *models.ProjectVersion, arch *dto.ArchitectureResponse) *serverinterfaces.ProjectVersionDetail {
	return &serverinterfaces.ProjectVersionDetail{
		ProjectVersionSummary: *versionSummary(v),
		State:                 arch,
	}
}
