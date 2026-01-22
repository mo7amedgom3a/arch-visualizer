package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceTypePricing represents pricing per service type (resource type)
type ServiceTypePricing struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	ResourceTypeID uint           `gorm:"not null;index" json:"resource_type_id"`
	TotalCost      float64        `gorm:"type:numeric(12,4);not null" json:"total_cost"`
	Currency       string         `gorm:"type:text;not null;check:currency IN ('USD','EUR','GBP')" json:"currency"`
	Period         string         `gorm:"type:text;not null;check:period IN ('hourly','monthly','yearly')" json:"period"`
	DurationSeconds int64         `gorm:"type:bigint;not null" json:"duration_seconds"`
	Provider       string         `gorm:"type:text;not null;check:provider IN ('aws','azure','gcp')" json:"provider"`
	Region         *string        `gorm:"type:text" json:"region,omitempty"`
	CalculatedAt   time.Time      `gorm:"default:now()" json:"calculated_at"`

	// Relationships
	Project      Project      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project,omitempty"`
	ResourceType ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type,omitempty"`
}

// TableName specifies the table name for GORM
func (ServiceTypePricing) TableName() string {
	return "service_type_pricing"
}
