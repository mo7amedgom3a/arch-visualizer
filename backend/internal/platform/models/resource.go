package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Resource represents a resource instance in a project
type Resource struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OriginalID     string         `gorm:"type:varchar(255);index" json:"original_id"` // Frontend node ID for reference resolution
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	ResourceTypeID uint           `gorm:"not null;index" json:"resource_type_id"`
	Name           string         `gorm:"type:text;not null" json:"name"`
	IsVisualOnly   bool           `gorm:"default:false" json:"is_visual_only"`
	Config         datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"config"`
	CreatedAt      time.Time      `gorm:"default:now()" json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Project      Project      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project,omitempty"`
	ResourceType ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type,omitempty"`

	// Containment relationships
	ParentResources []ResourceContainment `gorm:"foreignKey:ChildResourceID;constraint:OnDelete:CASCADE" json:"parent_containments,omitempty"`
	ChildResources  []ResourceContainment `gorm:"foreignKey:ParentResourceID;constraint:OnDelete:CASCADE" json:"child_containments,omitempty"`

	// Dependency relationships
	FromDependencies []ResourceDependency `gorm:"foreignKey:FromResourceID;constraint:OnDelete:CASCADE" json:"from_dependencies,omitempty"`
	ToDependencies   []ResourceDependency `gorm:"foreignKey:ToResourceID;constraint:OnDelete:CASCADE" json:"to_dependencies,omitempty"`

	// Pricing
	ResourcePricing []ResourcePricing `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"resource_pricing,omitempty"`

	// UI State
	UIState *ResourceUIState `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE" json:"ui_state,omitempty"`
}

// TableName specifies the table name for GORM
func (Resource) TableName() string {
	return "resources"
}
