package models

import (
	"time"

	"github.com/google/uuid"
)

// ServicePricing represents pricing per service (resource category)
type ServicePricing struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	CategoryID     uint           `gorm:"not null;index" json:"category_id"`
	TotalCost      float64        `gorm:"type:numeric(12,4);not null" json:"total_cost"`
	Currency       string         `gorm:"type:text;not null;check:currency IN ('USD','EUR','GBP')" json:"currency"`
	Period         string         `gorm:"type:text;not null;check:period IN ('hourly','monthly','yearly')" json:"period"`
	DurationSeconds int64         `gorm:"type:bigint;not null" json:"duration_seconds"`
	Provider       string         `gorm:"type:text;not null;check:provider IN ('aws','azure','gcp')" json:"provider"`
	Region         *string        `gorm:"type:text" json:"region,omitempty"`
	CalculatedAt   time.Time      `gorm:"default:now()" json:"calculated_at"`

	// Relationships
	Project  Project         `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project,omitempty"`
	Category ResourceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName specifies the table name for GORM
func (ServicePricing) TableName() string {
	return "service_pricing"
}
