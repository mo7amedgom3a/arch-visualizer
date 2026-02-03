package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"gorm.io/datatypes"
)

// GetArchitecture retrieves the full architecture for a project
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

		nodes[i] = dto.ArchitectureNode{
			ID:   res.ID,
			Type: res.Type.Name, // Or specialized type from metadata
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
		}
	}

	edges := make([]dto.ArchitectureEdge, 0)
	// Add containment edges
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
	// Add dependency edges
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

	// Helper to load variables from DB (TODO: implement variable repo)
	variables := make([]dto.ArchitectureVariable, 0)
	// mock variables for now if not implemented
	for _, v := range arch.Variables {
		variables = append(variables, dto.ArchitectureVariable{
			Name:        v.Name,
			Type:        v.Type,
			Value:       v.Default,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	outputs := make([]dto.ArchitectureOutput, 0)
	for _, o := range arch.Outputs {
		outputs = append(outputs, dto.ArchitectureOutput{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	return &dto.ArchitectureResponse{
		Nodes:     nodes,
		Edges:     edges,
		Variables: variables,
		Outputs:   outputs,
	}, nil
}

// SaveArchitecture saves the full architecture for a project
func (s *ProjectServiceImpl) SaveArchitecture(ctx context.Context, projectID uuid.UUID, req *dto.UpdateArchitectureRequest) (*dto.ArchitectureResponse, error) {
	// 1. Convert DTO to DiagramGraph (simplified) or directly to Architecture
	// Using DiagramGraph intermediate allows validation reuse
	diagramGraph := &graph.DiagramGraph{
		Nodes:     make(map[string]*graph.Node),
		Outputs:   make([]graph.Output, 0),
		Variables: make([]graph.Variable, 0),
	}

	// Map generic architecture DTO to DiagramGraph
	// Note: This is a simplified mapping. Real mapping might need more details.
	for _, node := range req.Nodes {
		var config map[string]interface{}
		if node.Data.Config != nil {
			config = node.Data.Config
		} else {
			config = make(map[string]interface{})
		}

		diagramGraph.Nodes[node.ID] = &graph.Node{
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

	// Edges -> Dependencies
	if diagramGraph.Edges == nil {
		diagramGraph.Edges = make([]*graph.Edge, 0)
	}

	for _, edge := range req.Edges {
		if edge.Type == "depends_on" || edge.Type == "dependency" {
			// Add dependency edge
			config := make(map[string]interface{})
			if edge.Label != "" {
				config["label"] = edge.Label
			}

			newEdge := &graph.Edge{
				ID:     edge.ID,
				Source: edge.Source,
				Target: edge.Target,
				Type:   "dependency",
				Config: config,
			}
			diagramGraph.Edges = append(diagramGraph.Edges, newEdge)
		}
		// Containment already handled by ParentID in nodes
	}

	// 1.1 Map Variables and Outputs from DTO
	for _, v := range req.Variables {
		diagramGraph.Variables = append(diagramGraph.Variables, graph.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Default:     v.Value,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	for _, o := range req.Outputs {
		diagramGraph.Outputs = append(diagramGraph.Outputs, graph.Output{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	// 2. Fetch project to know provider
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// 3. Map to Architecture
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.CloudProvider(project.CloudProvider))
	if err != nil {
		return nil, fmt.Errorf("failed to map architecture: %w", err)
	}

	// 4. Clean up existing architecture (naive approach: delete all, then create)
	// TODO: Implement transaction-safe cleanup and recreate

	// 5. Persist
	err = s.PersistArchitecture(ctx, projectID, arch, diagramGraph)
	if err != nil {
		return nil, fmt.Errorf("failed to persist architecture: %w", err)
	}

	// 6. Create Version Snapshot
	changes := "Updated architecture via API"
	snapshotJSON, _ := json.Marshal(req) // Store request payload as snapshot
	version := &models.ProjectVersion{
		ID:        uuid.New(),
		ProjectID: projectID,
		CreatedAt: time.Now(),
		CreatedBy: project.UserID, // Set to project owner for now
		Changes:   changes,
		Snapshot:  datatypes.JSON(snapshotJSON),
	}
	_ = s.versionRepo.Create(ctx, version) // Log version but don't fail if error

	// 7. Return updated
	return s.GetArchitecture(ctx, projectID)
}

// UpdateNode updates a single node in the architecture
func (s *ProjectServiceImpl) UpdateNode(ctx context.Context, projectID uuid.UUID, nodeID string, req *dto.UpdateNodeRequest) (*dto.ArchitectureNode, error) {
	// For now, efficient partial update requires DB support.
	// As fallback: Load -> Update in memory -> Save
	// This is inefficient but safe.

	fullArch, err := s.GetArchitecture(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var targetNode *dto.ArchitectureNode
	for i := range fullArch.Nodes {
		if fullArch.Nodes[i].ID == nodeID {
			targetNode = &fullArch.Nodes[i]
			break
		}
	}

	if targetNode == nil {
		return nil, fmt.Errorf("node not found")
	}

	// Apply updates
	if req.Position != nil {
		targetNode.Position = *req.Position
	}
	if req.Data != nil {
		if req.Data.Label != "" {
			targetNode.Data.Label = req.Data.Label
		}
		if req.Data.Config != nil {
			// Merge config
			for k, v := range req.Data.Config {
				if targetNode.Data.Config == nil {
					targetNode.Data.Config = make(map[string]interface{})
				}
				targetNode.Data.Config[k] = v
			}
		}
	}

	// Save full
	updateReq := &dto.UpdateArchitectureRequest{
		Nodes:     fullArch.Nodes,
		Edges:     fullArch.Edges,
		Variables: fullArch.Variables,
	}
	_, err = s.SaveArchitecture(ctx, projectID, updateReq)
	if err != nil {
		return nil, err
	}

	return targetNode, nil
}

// DeleteNode deletes a node from the architecture
func (s *ProjectServiceImpl) DeleteNode(ctx context.Context, projectID uuid.UUID, nodeID string) error {
	fullArch, err := s.GetArchitecture(ctx, projectID)
	if err != nil {
		return err
	}

	newNodes := make([]dto.ArchitectureNode, 0)
	for _, n := range fullArch.Nodes {
		if n.ID != nodeID {
			newNodes = append(newNodes, n)
		}
	}

	if len(newNodes) == len(fullArch.Nodes) {
		return fmt.Errorf("node not found")
	}

	// Filter edges
	newEdges := make([]dto.ArchitectureEdge, 0)
	for _, e := range fullArch.Edges {
		if e.Source != nodeID && e.Target != nodeID {
			newEdges = append(newEdges, e)
		}
	}

	updateReq := &dto.UpdateArchitectureRequest{
		Nodes:     newNodes,
		Edges:     newEdges,
		Variables: fullArch.Variables,
	}

	_, err = s.SaveArchitecture(ctx, projectID, updateReq)
	return err
}

// ValidateArchitecture validates the current architecture
func (s *ProjectServiceImpl) ValidateArchitecture(ctx context.Context, projectID uuid.UUID) (*dto.ValidationResponse, error) {
	// Load architecture
	arch, err := s.LoadArchitecture(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement actual validation using architecture.ValidateRules
	// For now returns mock valid response or rudimentary checks
	valid := true
	errors := make([]dto.ValidationIssue, 0)

	// Example check
	if len(arch.Resources) == 0 {
		valid = false
		errors = append(errors, dto.ValidationIssue{
			Type:     "structural",
			Message:  "Architecture is empty",
			Severity: "warning",
		})
	}

	return &dto.ValidationResponse{
		Valid:    valid,
		Errors:   errors,
		Warnings: []dto.ValidationIssue{},
	}, nil
}
