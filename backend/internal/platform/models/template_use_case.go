package models

import (
	"time"

	"github.com/google/uuid"
)

// TemplateUseCase represents a template use case item
type TemplateUseCase struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TemplateID   uuid.UUID `gorm:"type:uuid;not null;index" json:"template_id"`
	Icon         *string   `gorm:"type:varchar(100)" json:"icon,omitempty"`
	Title        string    `gorm:"type:varchar(255);not null" json:"title"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Template Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateUseCase) TableName() string {
	return "template_use_cases"
}
