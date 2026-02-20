package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectVersion represents a versioned snapshot link for a project.
// Each write creates a new Project row (the real snapshot) plus one of these
// chain entries so we can traverse history.
type ProjectVersion struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	ParentVersionID *uuid.UUID `gorm:"type:uuid;index" json:"parent_version_id"`
	VersionNumber   int        `gorm:"not null;default:1" json:"version_number"`
	Message         string     `gorm:"type:text" json:"message,omitempty"`
	CreatedAt       time.Time  `gorm:"default:now()" json:"created_at"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid" json:"created_by"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName specifies the table name for GORM
func (ProjectVersion) TableName() string {
	return "project_versions"
}
