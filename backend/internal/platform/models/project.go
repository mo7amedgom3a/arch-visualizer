package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents a cloud architecture design project
type Project struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	InfraToolID   uint           `gorm:"column:infra_tool;not null;index" json:"infra_tool"`
	Name          string         `gorm:"type:text;not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	CloudProvider string         `gorm:"type:text;not null;check:cloud_provider IN ('aws','azure','gcp')" json:"cloud_provider"`
	Region        string         `gorm:"type:text;not null" json:"region"`
	Thumbnail     string         `gorm:"type:text" json:"thumbnail"`
	Tags          []string       `gorm:"type:text[]" json:"tags"`
	ResourceCount int            `gorm:"-" json:"resourceCount"` // Calculated field
	EstimatedCost float64        `gorm:"-" json:"estimatedCost"` // Calculated field
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User               User                 `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	IACTarget          IACTarget            `gorm:"foreignKey:InfraToolID" json:"iac_target,omitempty"`
	Resources          []Resource           `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"resources,omitempty"`
	ProjectPricing     []ProjectPricing     `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project_pricing,omitempty"`
	ServicePricing     []ServicePricing     `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"service_pricing,omitempty"`
	ServiceTypePricing []ServiceTypePricing `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"service_type_pricing,omitempty"`
	ResourcePricing    []ResourcePricing    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"resource_pricing,omitempty"`
	Variables          []ProjectVariable    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"variables,omitempty"`
	Outputs            []ProjectOutput      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"outputs,omitempty"`
}

// TableName specifies the table name for GORM
func (Project) TableName() string {
	return "projects"
}
