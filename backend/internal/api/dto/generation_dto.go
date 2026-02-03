package dto

import (
	"time"

	"github.com/google/uuid"
)

// GenerateCodeRequest represents the request body for generating code
type GenerateCodeRequest struct {
	Tool    string             `json:"tool" binding:"required"` // terraform, pulumi, cdk
	Options *GenerationOptions `json:"options"`
}

type GenerationOptions struct {
	Format           string `json:"format"` // hcl, json, yaml
	IncludeOutputs   bool   `json:"includeOutputs"`
	IncludeVariables bool   `json:"includeVariables"`
	Modularity       string `json:"modularity"` // low, medium, high
}

// GeneratedFileResponse represents a single file in the response
type GeneratedFileResponse struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Content  string `json:"content,omitempty"` // Included in immediate generation response, excluded in history
	Size     int    `json:"size"`
}

// GenerationResponse represents the full response object
type GenerationResponse struct {
	ID           uuid.UUID               `json:"id"` // alias for generationId
	GenerationID uuid.UUID               `json:"generationId"`
	ProjectID    uuid.UUID               `json:"projectId,omitempty"`
	Status       string                  `json:"status"`
	Tool         string                  `json:"tool"`
	Files        []GeneratedFileResponse `json:"files,omitempty"`
	DownloadURL  string                  `json:"downloadUrl"`
	ExpiresAt    *time.Time              `json:"expiresAt"`
	CreatedAt    time.Time               `json:"createdAt"`
	ErrorMessage string                  `json:"errorMessage,omitempty"`
}

// GenerationHistoryResponse returns a list of generations
type GenerationHistoryResponse struct {
	Generations []GenerationResponse `json:"generations"`
}
