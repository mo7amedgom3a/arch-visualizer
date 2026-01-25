package models

import (
	"time"

	"github.com/google/uuid"
)

// TemplateFeature represents a template feature item
type TemplateFeature struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TemplateID   uuid.UUID `gorm:"type:uuid;not null;index" json:"template_id"`
	Feature      string    `gorm:"type:text;not null" json:"feature"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Template Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateFeature) TableName() string {
	return "template_features"
}
