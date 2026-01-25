package models

import "github.com/google/uuid"

// TemplateTechnology represents the template_technologies join table
type TemplateTechnology struct {
	TemplateID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"template_id"`
	TechnologyID uuid.UUID `gorm:"type:uuid;primaryKey" json:"technology_id"`

	// Relationships
	Template   Template   `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	Technology Technology `gorm:"foreignKey:TechnologyID;constraint:OnDelete:CASCADE" json:"technology,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateTechnology) TableName() string {
	return "template_technologies"
}
