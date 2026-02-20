package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"gorm.io/datatypes"
)

// cloneProjectSnapshotOptions defines how the clone write path should behave.
type cloneProjectSnapshotOptions struct {
	// sourceProjectID is the project to clone from.
	sourceProjectID uuid.UUID
	// createdBy is the user responsible for the new version.
	createdBy uuid.UUID
	// applyArch, if non-nil, replaces the architecture with the provided data instead of cloning.
	applyArch *dto.UpdateArchitectureRequest
	// message is the human-readable version description stored in project_versions.message.
	message string
}

// cloneProjectSnapshot is the core immutable write path:
//  1. Load source project + architecture.
//  2. Create a brand-new project row (new UUID).
//  3. Persist a fresh copy of all resources (new UUIDs) into the new project.
//  4. Insert a project_versions row linking to the previous version.
//
// It returns a VersionedOperationResult with the new project ID and version info.
func (s *ProjectServiceImpl) cloneProjectSnapshot(ctx context.Context, opts cloneProjectSnapshotOptions) (*serverinterfaces.VersionedOperationResult, *models.Project, error) {
	// 1. Load source project
	srcProject, err := s.projectRepo.FindByID(ctx, opts.sourceProjectID)
	if err != nil {
		return nil, nil, fmt.Errorf("cloneProjectSnapshot: source project not found: %w", err)
	}

	// 2. Resolve root project ID (the logical "lineage" anchor)
	rootProjectID := srcProject.RootProjectID
	if rootProjectID == nil {
		// Source IS the root — its own ID becomes the root for child versions
		srcID := srcProject.ID
		rootProjectID = &srcID
	}

	// Find the latest version for the source project so we can chain parent_version_id
	parentVersion, _ := s.versionRepo.GetLatestVersionForProject(ctx, opts.sourceProjectID)

	// 3. Determine architecture to persist
	var archReq *dto.UpdateArchitectureRequest
	if opts.applyArch != nil {
		archReq = opts.applyArch
	} else {
		// Clone current architecture
		currentArch, err := s.GetArchitecture(ctx, opts.sourceProjectID)
		if err != nil {
			return nil, nil, fmt.Errorf("cloneProjectSnapshot: load source architecture: %w", err)
		}
		archReq = &dto.UpdateArchitectureRequest{
			Nodes:     currentArch.Nodes,
			Edges:     currentArch.Edges,
			Variables: currentArch.Variables,
			Outputs:   currentArch.Outputs,
		}
	}

	// 4. Create the new project row
	newProject := &models.Project{
		ID:            uuid.New(),
		RootProjectID: rootProjectID,
		UserID:        srcProject.UserID,
		InfraToolID:   srcProject.InfraToolID,
		Name:          srcProject.Name,
		Description:   srcProject.Description,
		Tags:          srcProject.Tags,
		CloudProvider: srcProject.CloudProvider,
		Region:        srcProject.Region,
		Thumbnail:     srcProject.Thumbnail,
	}
	if err := s.projectRepo.Create(ctx, newProject); err != nil {
		return nil, nil, fmt.Errorf("cloneProjectSnapshot: create new project: %w", err)
	}

	// 5. Build DiagramGraph from archReq and persist as new architecture
	diagramGraph := s.archReqToDiagramGraph(archReq)
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.CloudProvider(newProject.CloudProvider))
	if err != nil {
		_ = s.projectRepo.Delete(ctx, newProject.ID)
		return nil, nil, fmt.Errorf("cloneProjectSnapshot: map architecture: %w", err)
	}
	if err := s.PersistArchitecture(ctx, newProject.ID, arch, diagramGraph); err != nil {
		_ = s.projectRepo.Delete(ctx, newProject.ID)
		return nil, nil, fmt.Errorf("cloneProjectSnapshot: persist architecture: %w", err)
	}

	// 6. Determine new version number
	newVersionNumber := 1
	var parentVersionID *uuid.UUID
	if parentVersion != nil {
		newVersionNumber = parentVersion.VersionNumber + 1
		parentVersionID = &parentVersion.ID
	}

	// 7. Insert project_versions chain entry
	createdBy := opts.createdBy
	if createdBy == uuid.Nil {
		createdBy = srcProject.UserID
	}
	newVersion := &models.ProjectVersion{
		ID:              uuid.New(),
		ProjectID:       newProject.ID,
		ParentVersionID: parentVersionID,
		VersionNumber:   newVersionNumber,
		Message:         opts.message,
		CreatedBy:       createdBy,
	}
	if err := s.versionRepo.Create(ctx, newVersion); err != nil {
		// Non-fatal: the project snapshot is valid even if the version record fails
		fmt.Printf("⚠️  Failed to create project_versions entry: %v\n", err)
	}

	return &serverinterfaces.VersionedOperationResult{
		NewProjectID:  newProject.ID,
		VersionID:     newVersion.ID,
		VersionNumber: newVersionNumber,
	}, newProject, nil
}

// archReqToDiagramGraph converts an UpdateArchitectureRequest to a DiagramGraph for persistence.
func (s *ProjectServiceImpl) archReqToDiagramGraph(req *dto.UpdateArchitectureRequest) *graph.DiagramGraph {
	dg := &graph.DiagramGraph{
		Nodes:     make(map[string]*graph.Node),
		Edges:     make([]*graph.Edge, 0),
		Variables: make([]graph.Variable, 0),
		Outputs:   make([]graph.Output, 0),
	}

	for _, node := range req.Nodes {
		config := node.Data.Config
		if config == nil {
			config = make(map[string]interface{})
		}

		// Preserve UI state in config so that PersistArchitecture can extract it
		if node.UIState != nil {
			styleBytes, _ := json.Marshal(node.UIState.Style)
			measuredBytes, _ := json.Marshal(node.UIState.Measured)
			config["_ui"] = map[string]interface{}{
				"position":   map[string]interface{}{"x": node.UIState.X, "y": node.UIState.Y},
				"width":      node.UIState.Width,
				"height":     node.UIState.Height,
				"style":      json.RawMessage(styleBytes),
				"measured":   json.RawMessage(measuredBytes),
				"selected":   node.UIState.Selected,
				"dragging":   node.UIState.Dragging,
				"resizing":   node.UIState.Resizing,
				"focusable":  node.UIState.Focusable,
				"selectable": node.UIState.Selectable,
				"zIndex":     node.UIState.ZIndex,
			}
		}

		dg.Nodes[node.ID] = &graph.Node{
			ID:           node.ID,
			ResourceType: node.Data.ResourceType,
			Label:        node.Data.Label,
			ParentID:     node.ParentID,
			PositionX:    int(node.Position.X),
			PositionY:    int(node.Position.Y),
			Config:       config,
			IsVisualOnly: node.Data.IsVisualOnly,
		}
	}

	for _, edge := range req.Edges {
		if edge.Type == "depends_on" || edge.Type == "dependency" {
			cfg := make(map[string]interface{})
			if edge.Label != "" {
				cfg["label"] = edge.Label
			}
			dg.Edges = append(dg.Edges, &graph.Edge{
				ID:     edge.ID,
				Source: edge.Source,
				Target: edge.Target,
				Type:   "dependency",
				Config: cfg,
			})
		}
	}

	for _, v := range req.Variables {
		dg.Variables = append(dg.Variables, graph.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Default:     v.Value,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}
	for _, o := range req.Outputs {
		dg.Outputs = append(dg.Outputs, graph.Output{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	return dg
}

// GetArchitecture retrieves the full architecture for a project snapshot
func (s *ProjectServiceImpl) GetArchitecture(ctx context.Context, projectID uuid.UUID) (*dto.ArchitectureResponse, error) {
	arch, err := s.LoadArchitecture(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture: %w", err)
	}

	nodes := make([]dto.ArchitectureNode, len(arch.Resources))
	for i, res := range arch.Resources {
		posX, posY := 0.0, 0.0
		if pos, ok := res.Metadata["position"].(map[string]interface{}); ok {
			if x, ok := pos["x"].(float64); ok {
				posX = x
			}
			if y, ok := pos["y"].(float64); ok {
				posY = y
			}
		}

		isVisualOnly := false
		if v, ok := res.Metadata["isVisualOnly"].(bool); ok {
			isVisualOnly = v
		}

		// Extract UI state if present
		var uiState *dto.NodeUIState
		if ui, ok := res.Metadata["ui"].(*graph.UIState); ok && ui != nil {
			w := 0.0
			h := 0.0
			if ui.Width != nil {
				w = *ui.Width
			}
			if ui.Height != nil {
				h = *ui.Height
			}
			uiState = &dto.NodeUIState{
				X:          ui.Position.X,
				Y:          ui.Position.Y,
				Width:      w,
				Height:     h,
				Style:      ui.Style,
				Measured:   ui.Measured,
				Selected:   ui.Selected,
				Dragging:   ui.Dragging,
				Resizing:   ui.Resizing,
				Focusable:  ui.Focusable,
				Selectable: ui.Selectable,
				ZIndex:     ui.ZIndex,
			}
		}

		nodes[i] = dto.ArchitectureNode{
			ID:   res.ID,
			Type: res.Type.Name,
			Position: dto.NodePosition{
				X: posX,
				Y: posY,
			},
			Data: dto.ArchitectureNodeData{
				Label:        res.Name,
				ResourceType: res.Type.Name,
				Config:       res.Metadata,
				IsVisualOnly: isVisualOnly,
			},
			ParentID: res.ParentID,
			UIState:  uiState,
		}
	}

	edges := make([]dto.ArchitectureEdge, 0)
	for parentID, children := range arch.Containments {
		for _, childID := range children {
			edges = append(edges, dto.ArchitectureEdge{
				ID:     fmt.Sprintf("contain-%s-%s", parentID, childID),
				Source: parentID,
				Target: childID,
				Type:   "contains",
			})
		}
	}
	for fromID, toIDs := range arch.Dependencies {
		for _, toID := range toIDs {
			edges = append(edges, dto.ArchitectureEdge{
				ID:     fmt.Sprintf("depend-%s-%s", fromID, toID),
				Source: fromID,
				Target: toID,
				Type:   "depends_on",
			})
		}
	}

	variables := make([]dto.ArchitectureVariable, 0, len(arch.Variables))
	for _, v := range arch.Variables {
		variables = append(variables, dto.ArchitectureVariable{
			Name:        v.Name,
			Type:        v.Type,
			Value:       v.Default,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	outputs := make([]dto.ArchitectureOutput, 0, len(arch.Outputs))
	for _, o := range arch.Outputs {
		outputs = append(outputs, dto.ArchitectureOutput{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	warnings := make([]dto.ValidationIssue, 0, len(arch.Warnings))
	for _, w := range arch.Warnings {
		warnings = append(warnings, dto.ValidationIssue{
			Type:     "auto-correction",
			Message:  w.Message,
			NodeID:   w.ResourceID,
			Severity: "warning",
		})
	}

	versionID := ""
	versions, err := s.versionRepo.ListByProjectID(ctx, projectID)
	if err == nil && len(versions) > 0 {
		versionID = versions[0].ID.String()
	}

	return &dto.ArchitectureResponse{
		VersionID: versionID,
		Nodes:     nodes,
		Edges:     edges,
		Variables: variables,
		Outputs:   outputs,
		Warnings:  warnings,
	}, nil
}

// SaveArchitecture creates a new immutable project snapshot with the given architecture.
// Kept for backward compatibility with the scenario runners and pipeline orchestrator tests.
func (s *ProjectServiceImpl) SaveArchitecture(ctx context.Context, projectID uuid.UUID, req *dto.UpdateArchitectureRequest) (*serverinterfaces.VersionedArchitectureResult, error) {
	versionedResult, newProject, err := s.cloneProjectSnapshot(ctx, cloneProjectSnapshotOptions{
		sourceProjectID: projectID,
		applyArch:       req,
	})
	if err != nil {
		return nil, fmt.Errorf("SaveArchitecture: %w", err)
	}

	arch, err := s.GetArchitecture(ctx, newProject.ID)
	if err != nil {
		return nil, fmt.Errorf("SaveArchitecture: load saved architecture: %w", err)
	}

	return &serverinterfaces.VersionedArchitectureResult{
		VersionedOperationResult: *versionedResult,
		Architecture:             arch,
	}, nil
}

// ensure datatypes import is used (it's referenced in PersistArchitecture in project_service.go)
var _ = datatypes.JSON(nil)
