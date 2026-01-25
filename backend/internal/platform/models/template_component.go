package models

import (
	"time"

	"github.com/google/uuid"
)

// TemplateComponent represents an individual component of a template
type TemplateComponent struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TemplateID   uuid.UUID `gorm:"type:uuid;not null;index" json:"template_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Service      string    `gorm:"type:varchar(255);not null" json:"service"`
	Configuration *string  `gorm:"type:text" json:"configuration,omitempty"`
	MonthlyCost  float64   `gorm:"type:decimal(10,2);default:0" json:"monthly_cost"`
	Purpose      *string   `gorm:"type:text" json:"purpose,omitempty"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Template Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateComponent) TableName() string {
	return "template_components"
}
