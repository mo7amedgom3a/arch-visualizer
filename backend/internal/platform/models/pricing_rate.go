package models

import (
	"time"

	"gorm.io/datatypes"
)

// PricingRate represents pricing rates stored in the database
type PricingRate struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Provider       string         `gorm:"type:varchar(20);not null;index" json:"provider"`
	ResourceType   string         `gorm:"type:varchar(100);not null;index" json:"resource_type"`
	ComponentName  string         `gorm:"type:varchar(100);not null" json:"component_name"`
	PricingModel   string         `gorm:"type:varchar(50);not null" json:"pricing_model"`
	Unit           string         `gorm:"type:varchar(50);not null" json:"unit"`
	Rate           float64        `gorm:"type:numeric(14,6);not null" json:"rate"`
	Currency       string         `gorm:"type:varchar(10);default:'USD'" json:"currency"`
	Region         *string        `gorm:"type:varchar(50);index" json:"region,omitempty"`
	InstanceType   *string        `gorm:"type:varchar(50);index" json:"instance_type,omitempty"`
	OperatingSystem *string       `gorm:"type:varchar(20);default:'linux';index" json:"operating_system,omitempty"`
	EffectiveFrom  time.Time      `gorm:"not null;default:now()" json:"effective_from"`
	EffectiveTo    *time.Time     `gorm:"index" json:"effective_to,omitempty"`
	Metadata       datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (PricingRate) TableName() string {
	return "pricing_rates"
}
