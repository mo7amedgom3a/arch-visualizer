package models

import (
	"github.com/google/uuid"
)

// ResourceContainment represents parent-child relationships between resources (VPC → Subnet → EC2)
type ResourceContainment struct {
	ParentResourceID uuid.UUID `gorm:"type:uuid;primaryKey" json:"parent_resource_id"`
	ChildResourceID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"child_resource_id"`

	// Relationships
	ParentResource Resource `gorm:"foreignKey:ParentResourceID;constraint:OnDelete:CASCADE" json:"parent_resource,omitempty"`
	ChildResource  Resource `gorm:"foreignKey:ChildResourceID;constraint:OnDelete:CASCADE" json:"child_resource,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceContainment) TableName() string {
	return "resource_containment"
}
