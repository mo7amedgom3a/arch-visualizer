package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ProjectVersion represents a snapshot of a project's architecture
type ProjectVersion struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	CreatedBy uuid.UUID      `gorm:"type:uuid" json:"created_by"` // Nullable if system gen or unknown
	Changes   string         `gorm:"type:text" json:"changes"`
	Snapshot  datatypes.JSON `gorm:"type:jsonb" json:"snapshot"` // Stores full architecture state
}

// TableName specifies the table name for GORM
func (ProjectVersion) TableName() string {
	return "project_versions"
}
