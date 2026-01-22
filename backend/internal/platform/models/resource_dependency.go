package models

import (
	"github.com/google/uuid"
)

// ResourceDependency represents directed dependencies between resources
type ResourceDependency struct {
	FromResourceID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"from_resource_id"`
	ToResourceID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"to_resource_id"`
	DependencyTypeID uint      `gorm:"not null;index" json:"dependency_type_id"`

	// Relationships
	FromResource   Resource       `gorm:"foreignKey:FromResourceID;constraint:OnDelete:CASCADE" json:"from_resource,omitempty"`
	ToResource     Resource       `gorm:"foreignKey:ToResourceID;constraint:OnDelete:CASCADE" json:"to_resource,omitempty"`
	DependencyType DependencyType `gorm:"foreignKey:DependencyTypeID" json:"dependency_type,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceDependency) TableName() string {
	return "resource_dependencies"
}
