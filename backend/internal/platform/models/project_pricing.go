package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectPricing represents pricing estimates for entire projects
type ProjectPricing struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	TotalCost      float64        `gorm:"type:numeric(12,4);not null" json:"total_cost"`
	Currency       string         `gorm:"type:text;not null;check:currency IN ('USD','EUR','GBP')" json:"currency"`
	Period         string         `gorm:"type:text;not null;check:period IN ('hourly','monthly','yearly')" json:"period"`
	DurationSeconds int64         `gorm:"type:bigint;not null" json:"duration_seconds"`
	Provider       string         `gorm:"type:text;not null;check:provider IN ('aws','azure','gcp')" json:"provider"`
	Region         *string        `gorm:"type:text" json:"region,omitempty"`
	CalculatedAt   time.Time      `gorm:"default:now()" json:"calculated_at"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project,omitempty"`
}

// TableName specifies the table name for GORM
func (ProjectPricing) TableName() string {
	return "project_pricing"
}
